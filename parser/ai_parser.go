package parser

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"ai-shell-windows/utils"
)

type AIParser struct {
	apiKey     string
	baseURL    string
	model      string
	httpClient *http.Client
}

type groqChatRequest struct {
	Model    string            `json:"model"`
	Messages []groqChatMessage `json:"messages"`
}

type groqChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type groqChatResponse struct {
	Choices []struct {
		Message groqChatMessage `json:"message"`
	} `json:"choices"`
}

type aiIntentPayload struct {
	Action        string `json:"action"`
	Target        string `json:"target"`
	Source        string `json:"source"`
	Destination   string `json:"destination"`
	RequiresInfo  bool   `json:"requires_info"`
	Clarification string `json:"clarification"`
	Explanation   string `json:"explanation"`
}

func NewAIParser(apiKey, baseURL, model string, timeout time.Duration) AIParser {
	return AIParser{
		apiKey:  apiKey,
		baseURL: strings.TrimRight(baseURL, "/"),
		model:   model,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (p AIParser) Parse(input string) (Intent, error) {
	intent := Intent{
		RawInput:   input,
		Normalized: utils.NormalizeText(input),
	}

	if p.apiKey == "" {
		return intent, ErrAIUnavailable
	}

	requestBody := groqChatRequest{
		Model: p.model,
		Messages: []groqChatMessage{
			{
				Role: "system",
				Content: "You translate Windows terminal requests into JSON. " +
					"Allowed actions: list_files, list_folders, print_working_dir, delete_file, create_folder, rename_file, unknown. " +
					"Return JSON only with keys action, target, source, destination, requires_info, clarification, explanation. " +
					"Use list_folders when the user only wants folders. " +
					"If required details are missing, set requires_info=true and provide a short clarification question. " +
					"Never include markdown fences.",
			},
			{
				Role:    "user",
				Content: input,
			},
		},
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return intent, err
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, p.baseURL+"/chat/completions", bytes.NewReader(bodyBytes))
	if err != nil {
		return intent, err
	}

	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return intent, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return intent, err
	}

	if resp.StatusCode >= 300 {
		return intent, fmt.Errorf("groq api error: %s", strings.TrimSpace(string(responseBody)))
	}

	var chatResponse groqChatResponse
	if err := json.Unmarshal(responseBody, &chatResponse); err != nil {
		return intent, err
	}
	if len(chatResponse.Choices) == 0 {
		return intent, ErrAIResponse
	}

	payload, err := decodeAIIntent(chatResponse.Choices[0].Message.Content)
	if err != nil {
		return intent, err
	}

	intent.Action = payload.Action
	intent.Target = payload.Target
	intent.Source = payload.Source
	intent.Destination = payload.Destination
	intent.RequiresInfo = payload.RequiresInfo
	intent.Clarification = payload.Clarification
	intent.Explanation = payload.Explanation

	switch intent.Action {
	case ActionListFiles, ActionListFolders, ActionPrintWorkingDir, ActionDeleteFile, ActionCreateFolder, ActionRenameFile:
	case ActionUnknown, "":
		return intent, ErrUnknownIntent
	default:
		return intent, ErrAIResponse
	}

	if intent.RequiresInfo {
		if intent.Clarification == "" {
			intent.Clarification = "Can you clarify the missing file or folder details?"
		}
		return intent, nil
	}

	return intent, nil
}

func decodeAIIntent(content string) (aiIntentPayload, error) {
	var payload aiIntentPayload
	rawJSON := strings.TrimSpace(content)
	rawJSON = strings.TrimPrefix(rawJSON, "```json")
	rawJSON = strings.TrimPrefix(rawJSON, "```")
	rawJSON = strings.TrimSuffix(rawJSON, "```")
	rawJSON = strings.TrimSpace(rawJSON)

	if !strings.HasPrefix(rawJSON, "{") {
		start := strings.Index(rawJSON, "{")
		end := strings.LastIndex(rawJSON, "}")
		if start == -1 || end == -1 || end < start {
			return payload, ErrAIResponse
		}
		rawJSON = rawJSON[start : end+1]
	}

	if err := json.Unmarshal([]byte(rawJSON), &payload); err != nil {
		return payload, ErrAIResponse
	}

	return payload, nil
}
