package zokyo

import (
	"MusicDev33/mdapi3/internal/config"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ChatEngine string

const (
	EngineChatGPT  ChatEngine = "chatgpt"
	EngineClaude   ChatEngine = "claude"
	EngineDeepSeek ChatEngine = "deepseek"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AgentConfig struct {
	TopP        float64 `json:"top_p"`
	Temperature float64 `json:"temperature"`
}

// Token counting - simplified version, would need proper implementation for production
func CountTokens(messages []Message, threshold int) bool {
	total := 0
	for _, msg := range messages {
		total += len(msg.Content) / 4
	}
	return total <= threshold
}

func GetSystemPrompt() string {
	now := time.Now()
	return fmt.Sprintf("The date is %s in San Francisco.", now.Format("January 2, 2006, 3:04 PM (Monday)"))
}

func GenerateChat(engine ChatEngine, messages []Message, agentConfig AgentConfig) (string, error) {
	cfg := config.Get()

	switch engine {
	case EngineChatGPT:
		return generateChatOpenAI(messages, agentConfig, cfg.KeyOpenAI)
	case EngineClaude:
		return generateChatClaude(messages, agentConfig, cfg.KeyAnthropic)
	case EngineDeepSeek:
		return generateChatDeepSeek(messages, agentConfig, cfg.KeyDeepSeek)
	default:
		return "", errors.New("unsupported chat engine")
	}
}

func generateChatOpenAI(messages []Message, agentConfig AgentConfig, apiKey string) (string, error) {
	url := "https://api.openai.com/v1/chat/completions"

	reqBody := map[string]any{
		"messages":    messages,
		"top_p":       agentConfig.TopP,
		"temperature": agentConfig.Temperature,
		"model":       "gpt-5-2025-08-07",
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("OpenAI API error: %s", string(body))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", errors.New("no choices in response")
	}

	choice, ok := choices[0].(map[string]interface{})
	if !ok {
		return "", errors.New("invalid choice format")
	}

	message, ok := choice["message"].(map[string]interface{})
	if !ok {
		return "", errors.New("invalid message format")
	}

	content, ok := message["content"].(string)
	if !ok {
		return "", errors.New("invalid content format")
	}

	return content, nil
}

func generateChatClaude(messages []Message, agentConfig AgentConfig, apiKey string) (string, error) {
	url := "https://api.anthropic.com/v1/messages"

	reqBody := map[string]interface{}{
		"messages":    messages,
		"temperature": agentConfig.Temperature,
		"model":       "claude-sonnet-4-5",
		"max_tokens":  4096,
		"system":      GetSystemPrompt(),
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("anthropic API error: %s", string(body))
	}

	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	content, ok := result["content"].([]any)
	if !ok || len(content) == 0 {
		return "", errors.New("no content in response")
	}

	contentBlock, ok := content[0].(map[string]any)
	if !ok {
		return "", errors.New("invalid content format")
	}

	text, ok := contentBlock["text"].(string)
	if !ok {
		return "", errors.New("invalid text format")
	}

	return text, nil
}

func generateChatDeepSeek(messages []Message, agentConfig AgentConfig, apiKey string) (string, error) {
	url := "https://api.deepseek.com/v1/chat/completions"

	reqBody := map[string]any{
		"messages":    messages,
		"temperature": agentConfig.Temperature,
		"model":       "deepseek-chat",
		"max_tokens":  1024,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("DeepSeek API error: %s", string(body))
	}

	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	choices, ok := result["choices"].([]any)
	if !ok || len(choices) == 0 {
		return "", errors.New("no choices in response")
	}

	choice, ok := choices[0].(map[string]any)
	if !ok {
		return "", errors.New("invalid choice format")
	}

	message, ok := choice["message"].(map[string]any)
	if !ok {
		return "", errors.New("invalid message format")
	}

	content, ok := message["content"].(string)
	if !ok {
		return "", errors.New("invalid content format")
	}

	return content, nil
}
