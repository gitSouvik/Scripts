package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
)

var fetchServer *http.Server

func launchFetch(ch chan fetchEvtMsg, cwd string) tea.Cmd {
	return func() tea.Msg {
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "POST only", http.StatusMethodNotAllowed)
				return
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "read error", http.StatusBadRequest)
				return
			}

			var data struct {
				Name  string `json:"name"`
				URL   string `json:"url"`
				Tests []struct {
					Input  string `json:"input"`
					Output string `json:"output"`
				} `json:"tests"`
			}
			if err := json.Unmarshal(body, &data); err != nil {
				http.Error(w, "json error", http.StatusBadRequest)
				return
			}

			problemName := data.Name
			if len(problemName) > 0 {
				problemName = string(problemName[0])
			}

			// Get unique name
			counter := 1
			unique := problemName
			for {
				if _, err := os.Stat(filepath.Join(cwd, unique+".cpp")); os.IsNotExist(err) {
					break
				}
				unique = fmt.Sprintf("%s%d", problemName, counter)
				counter++
			}
			problemName = unique

			for i, test := range data.Tests {
				os.WriteFile(filepath.Join(cwd, fmt.Sprintf("%s-%d.in", problemName, i+1)), []byte(test.Input), 0644)
				os.WriteFile(filepath.Join(cwd, fmt.Sprintf("%s-%d.out", problemName, i+1)), []byte(test.Output), 0644)
			}

			cppFile := filepath.Join(cwd, problemName+".cpp")
			f, err := os.Create(cppFile)
			if err == nil {
				defer f.Close()
				fmt.Fprintf(f, "/*\n * Author: snowdust\n * Problem: %s\n * Source: %s\n */\n\n", problemName, data.URL)
				
				templatePath := "template.cpp"
				if _, err := os.Stat(templatePath); os.IsNotExist(err) {
					exePath, _ := os.Executable()
					templatePath = filepath.Join(filepath.Dir(exePath), "template.cpp")
				}
				
				if tmpl, err := os.ReadFile(templatePath); err == nil {
					f.Write(tmpl)
				}
			}

			w.WriteHeader(http.StatusOK)
			ch <- fetchEvtMsg(fmt.Sprintf("Received: %s (%d tests)", problemName, len(data.Tests)))
		})

		fetchServer = &http.Server{Addr: ":54321", Handler: mux}
		go fetchServer.ListenAndServe()
		
		return nil
	}
}

func stopFetch() {
	if fetchServer != nil {
		fetchServer.Shutdown(context.Background())
		fetchServer = nil
	}
}
