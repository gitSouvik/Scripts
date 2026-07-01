package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// ── Execution Logic ───────────────────────────────────────────────────────────

func compileFile(name string, mode string, ch chan outLine) bool {
	cppFile := name + ".cpp"
	ch <- outLine{"log", fmt.Sprintf("[%s] compiling %s...", mode, cppFile)}
	
	cmd := exec.Command("g++-15", "-O2", "-std=c++20", "-DWOOF_", "-o", name, cppFile)
	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf
	
	start := time.Now()
	err := cmd.Run()
	elapsed := time.Since(start)
	
	if err != nil {
		ch <- outLine{"err", fmt.Sprintf("compilation failed (%dms)", elapsed.Milliseconds())}
		for _, line := range strings.Split(strings.TrimSpace(errBuf.String()), "\n") {
			if line != "" {
				ch <- outLine{"err", "  " + line}
			}
		}
		return false
	}
	ch <- outLine{"ok", fmt.Sprintf("compiled successfully (%dms)", elapsed.Milliseconds())}
	return true
}

func launchRun(name string, ch chan outLine) tea.Cmd {
	return func() tea.Msg {
		defer close(ch)
		
		if !compileFile(name, "run", ch) {
			return nil
		}
		
		inFiles, _ := filepath.Glob(name + "-*.in")
		if len(inFiles) == 0 {
			ch <- outLine{"err", "No test cases found (*.in)"}
			return nil
		}
		
		passed, failed := 0, 0
		
		for _, inFile := range inFiles {
			base := filepath.Base(inFile)
			base = strings.TrimPrefix(base, name+"-")
			testNum := strings.TrimSuffix(base, ".in")
			outFile := name + "-" + testNum + ".out"
			
			ch <- outLine{"sep", "-------------------------"}
			
			inContent, _ := os.ReadFile(inFile)
			cmd := exec.Command("./" + name)
			cmd.Stdin = bytes.NewReader(inContent)
			var outBuf, errBuf bytes.Buffer
			cmd.Stdout = &outBuf
			cmd.Stderr = &errBuf
			
			start := time.Now()
			runErr := cmd.Run()
			elapsed := time.Since(start)
			
			if runErr != nil {
				ch <- outLine{"fail", fmt.Sprintf("Test %s Failed: Runtime error / crash (%dms)", testNum, elapsed.Milliseconds())}
				failed++
				continue
			}
			
			actualOut := strings.TrimSpace(outBuf.String())
			expectedBytes, expErr := os.ReadFile(outFile)
			expectedOut := ""
			if expErr == nil {
				expectedOut = strings.TrimSpace(string(expectedBytes))
			}
			
			if errBuf.Len() > 0 {
				ch <- outLine{"log", "Debug:"}
				for _, line := range strings.Split(strings.TrimSpace(errBuf.String()), "\n") {
					ch <- outLine{"debug", line}
				}
			}
			
			if expErr != nil {
				ch <- outLine{"fail", fmt.Sprintf("Test %s Failed: missing expected output file", testNum)}
				failed++
			} else if actualOut == expectedOut {
				ch <- outLine{"pass", fmt.Sprintf("Test %s Passed (%dms)", testNum, elapsed.Milliseconds())}
				passed++
			} else {
				ch <- outLine{"fail", fmt.Sprintf("Test %s Failed (%dms)", testNum, elapsed.Milliseconds())}
				failed++
				
				ch <- outLine{"log", "Diff (Actual vs Expected):"}
				actualLines := strings.Split(actualOut, "\n")
				expectedLines := strings.Split(expectedOut, "\n")
				
				maxLines := len(actualLines)
				if len(expectedLines) > maxLines {
					maxLines = len(expectedLines)
				}
				
				for i := 0; i < maxLines; i++ {
					aLine := ""
					if i < len(actualLines) {
						aLine = actualLines[i]
					}
					eLine := ""
					if i < len(expectedLines) {
						eLine = expectedLines[i]
					}
					
					if aLine == eLine {
						ch <- outLine{"output", aLine}
					} else {
						// Token diff
						aTokens := strings.Fields(aLine)
						eTokens := strings.Fields(eLine)
						
						maxTokens := len(aTokens)
						if len(eTokens) > maxTokens {
							maxTokens = len(eTokens)
						}
						
						var aDiff, eDiff []string
						for j := 0; j < maxTokens; j++ {
							aTok := ""
							if j < len(aTokens) {
								aTok = aTokens[j]
							}
							eTok := ""
							if j < len(eTokens) {
								eTok = eTokens[j]
							}
							
							if aTok == eTok {
								aDiff = append(aDiff, aTok)
								eDiff = append(eDiff, eTok)
							} else {
								if aTok != "" {
									aDiff = append(aDiff, "\033[1;31m"+aTok+"\033[0m") // red
								}
								if eTok != "" {
									eDiff = append(eDiff, "\033[1;32m"+eTok+"\033[0m") // green
								}
							}
						}
						
						if aLine != "" {
							ch <- outLine{"diff-act", "Act: " + strings.Join(aDiff, " ")}
						}
						if eLine != "" {
							ch <- outLine{"diff-exp", "Exp: " + strings.Join(eDiff, " ")}
						}
					}
				}
			}
		}
		
		ch <- outLine{"sep", "========================="}
		ch <- outLine{"log", fmt.Sprintf("Summary: %d passed, %d failed", passed, failed)}
		return nil
	}
}

func launchRrun(name, testNum string, ch chan outLine) tea.Cmd {
	return func() tea.Msg {
		defer close(ch)
		
		if !compileFile(name, "debug", ch) {
			return nil
		}
		
		if testNum == "" {
			// No interactive keyboard input in this TUI yet, because TUI owns stdin
			ch <- outLine{"err", "Keyboard interactive input is not supported inside TUI."}
			ch <- outLine{"err", "Please specify a test case number, or 0 for all."}
			return nil
		}
		
		if testNum == "0" {
			inFiles, _ := filepath.Glob(name + "-*.in")
			for _, inFile := range inFiles {
				runSingleTest(name, inFile, ch)
			}
		} else {
			inFile := fmt.Sprintf("%s-%s.in", name, testNum)
			if _, err := os.Stat(inFile); os.IsNotExist(err) {
				ch <- outLine{"err", fmt.Sprintf("Test file missing: %s", inFile)}
				return nil
			}
			runSingleTest(name, inFile, ch)
		}
		
		return nil
	}
}

func runSingleTest(name, inFile string, ch chan outLine) {
	ch <- outLine{"sep", "-------------------------"}
	ch <- outLine{"log", fmt.Sprintf("Running %s", inFile)}
	
	inContent, _ := os.ReadFile(inFile)
	ch <- outLine{"log", "Input:"}
	for _, line := range strings.Split(strings.TrimSpace(string(inContent)), "\n") {
		ch <- outLine{"input", line}
	}
	
	cmd := exec.Command("./" + name)
	cmd.Stdin = bytes.NewReader(inContent)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	
	start := time.Now()
	cmd.Run()
	elapsed := time.Since(start)
	
	ch <- outLine{"log", fmt.Sprintf("Output (%dms):", elapsed.Milliseconds())}
	for _, line := range strings.Split(strings.TrimSpace(outBuf.String()), "\n") {
		ch <- outLine{"output", line}
	}
	
	if errBuf.Len() > 0 {
		ch <- outLine{"log", "Debug:"}
		for _, line := range strings.Split(strings.TrimSpace(errBuf.String()), "\n") {
			ch <- outLine{"debug", line}
		}
	}
}

func launchGen(name, seed string, ch chan outLine) tea.Cmd {
	return func() tea.Msg {
		defer close(ch)
		
		if _, err := os.Stat("gen.cpp"); os.IsNotExist(err) {
			ch <- outLine{"err", "gen.cpp not found"}
			return nil
		}
		
		ch <- outLine{"log", "[gen] compiling gen.cpp..."}
		genCompile := exec.Command("g++-15", "-O2", "-std=c++20", "-o", "gen", "gen.cpp")
		var genErrBuf bytes.Buffer
		genCompile.Stderr = &genErrBuf
		if err := genCompile.Run(); err != nil {
			ch <- outLine{"err", "gen.cpp compilation failed"}
			for _, line := range strings.Split(strings.TrimSpace(genErrBuf.String()), "\n") {
				if line != "" {
					ch <- outLine{"err", "  " + line}
				}
			}
			return nil
		}
		ch <- outLine{"ok", "gen compiled"}
		
		if !compileFile(name, "gen", ch) {
			return nil
		}
		
		args := []string{}
		if seed != "" {
			args = append(args, seed)
		}
		genRun := exec.Command("./gen", args...)
		var inBuf bytes.Buffer
		genRun.Stdout = &inBuf
		if err := genRun.Run(); err != nil {
			ch <- outLine{"err", fmt.Sprintf("gen execution failed: %v", err)}
			return nil
		}
		
		inputData := inBuf.Bytes()
		ch <- outLine{"sep", "-------------------------"}
		ch <- outLine{"log", "Generated Input:"}
		for _, line := range strings.Split(strings.TrimSpace(string(inputData)), "\n") {
			ch <- outLine{"input", line}
		}
		
		solCmd := exec.Command("./" + name)
		solCmd.Stdin = bytes.NewReader(inputData)
		var outBuf, errBuf bytes.Buffer
		solCmd.Stdout = &outBuf
		solCmd.Stderr = &errBuf
		
		solCmd.Run()
		
		ch <- outLine{"log", "Output:"}
		for _, line := range strings.Split(strings.TrimSpace(outBuf.String()), "\n") {
			ch <- outLine{"output", line}
		}
		
		if errBuf.Len() > 0 {
			ch <- outLine{"log", "Debug:"}
			for _, line := range strings.Split(strings.TrimSpace(errBuf.String()), "\n") {
				ch <- outLine{"debug", line}
			}
		}
		
		return nil
	}
}

// ── Custom Test Run ───────────────────────────────────────────────────────────

func launchCustomRun(name string, input string, ch chan outLine) tea.Cmd {
	return func() tea.Msg {
		defer close(ch)

		if !compileFile(name, "custom", ch) {
			return nil
		}

		inputData := []byte(input)

		ch <- outLine{"sep", "-------------------------"}
		ch <- outLine{"log", "Input:"}
		for _, line := range strings.Split(strings.TrimSpace(input), "\n") {
			ch <- outLine{"input", line}
		}

		cmd := exec.Command("./" + name)
		cmd.Stdin = bytes.NewReader(inputData)
		var outBuf, errBuf bytes.Buffer
		cmd.Stdout = &outBuf
		cmd.Stderr = &errBuf

		start := time.Now()
		cmd.Run()
		elapsed := time.Since(start)

		ch <- outLine{"log", fmt.Sprintf("Output (%dms):", elapsed.Milliseconds())}
		for _, line := range strings.Split(strings.TrimSpace(outBuf.String()), "\n") {
			ch <- outLine{"output", line}
		}

		if errBuf.Len() > 0 {
			ch <- outLine{"log", "Debug:"}
			for _, line := range strings.Split(strings.TrimSpace(errBuf.String()), "\n") {
				ch <- outLine{"debug", line}
			}
		}

		// Drop compiled binary (no .in files were made)
		_ = filepath.Join(".", name) // keep import used
		return nil
	}
}
