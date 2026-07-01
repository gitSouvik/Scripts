package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

var currentApiKey string

type aiTipMsg struct {
	tip string
}

type aiAckMsg struct {
	text string
}

func loadApiKey(cwd string) {
	keyFile := filepath.Join(cwd, ".cpx_key")
	if data, err := os.ReadFile(keyFile); err == nil {
		currentApiKey = strings.TrimSpace(string(data))
	}
}

func saveApiKey(cwd, key string) error {
	currentApiKey = strings.TrimSpace(key)
	keyFile := filepath.Join(cwd, ".cpx_key")
	return os.WriteFile(keyFile, []byte(currentApiKey), 0644)
}

func launchAITip() tea.Cmd {
	return func() tea.Msg {
		if currentApiKey == "" {
			return aiTipMsg{"for funny facts click ctrl + p to add free gemini key"}
		}

		url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent?key=" + currentApiKey

		payload := map[string]interface{}{
			"contents": []map[string]interface{}{
				{
					"parts": []map[string]interface{}{
						{
							"text": "You are a witty developer. Provide a very short programming joke, a fun fact, or a quick dev talk. Do not repeat previous responses. Make it a single, short sentence without quotes or introductory text.",
						},
					},
				},
			},
			"generationConfig": map[string]interface{}{
				"temperature": 1.0,
			},
		}

		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			return aiTipMsg{"Could not load AI tip"}
		}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
		if err != nil {
			return aiTipMsg{"Could not load AI tip"}
		}

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return aiTipMsg{"Could not load AI tip"}
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return aiTipMsg{"Could not load AI tip"}
		}

		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			return aiTipMsg{"Could not parse AI response"}
		}

		if errObj, ok := result["error"].(map[string]interface{}); ok {
			if msgStr, ok := errObj["message"].(string); ok {
				if strings.Contains(msgStr, "quota") || strings.Contains(msgStr, "Quota") {
					return aiTipMsg{"API Error: Quota exceeded (check billing or limits)"}
				}
				if len(msgStr) > 60 {
					msgStr = msgStr[:57] + "..."
				}
				return aiTipMsg{"API Error: " + msgStr}
			}
		}

		if candidates, ok := result["candidates"].([]interface{}); ok && len(candidates) > 0 {
			if candidate, ok := candidates[0].(map[string]interface{}); ok {
				if content, ok := candidate["content"].(map[string]interface{}); ok {
					if parts, ok := content["parts"].([]interface{}); ok && len(parts) > 0 {
						if part, ok := parts[0].(map[string]interface{}); ok {
							if text, ok := part["text"].(string); ok {
								tip := strings.TrimSpace(text)
								return aiTipMsg{tip}
							}
						}
					}
				}
			}
		}

		return aiTipMsg{"Could not load AI tip (No response)"}
	}
}

func launchAIAck(snippetName string) tea.Cmd {
	return func() tea.Msg {
		if currentApiKey == "" {
			return aiAckMsg{"Saved snippet: " + snippetName}
		}

		url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent?key=" + currentApiKey

		payload := map[string]interface{}{
			"contents": []map[string]interface{}{
				{
					"parts": []map[string]interface{}{
						{
							"text": "You are a witty programming assistant. The user just saved a competitive programming code snippet named '" + snippetName + "'. Write a very short, one-sentence congratulatory or witty comment about this. Do not use quotes or introductory text, just return the comment.",
						},
					},
				},
			},
		}

		jsonPayload, _ := json.Marshal(payload)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
		if err != nil {
			return aiAckMsg{"Saved snippet: " + snippetName}
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return aiAckMsg{"Saved snippet: " + snippetName}
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			return aiAckMsg{"Saved snippet: " + snippetName}
		}

		if candidates, ok := result["candidates"].([]interface{}); ok && len(candidates) > 0 {
			if candidate, ok := candidates[0].(map[string]interface{}); ok {
				if content, ok := candidate["content"].(map[string]interface{}); ok {
					if parts, ok := content["parts"].([]interface{}); ok && len(parts) > 0 {
						if part, ok := parts[0].(map[string]interface{}); ok {
							if text, ok := part["text"].(string); ok {
								ack := strings.TrimSpace(text)
								return aiAckMsg{ack}
							}
						}
					}
				}
			}
		}

		return aiAckMsg{"Saved snippet: " + snippetName}
	}
}
