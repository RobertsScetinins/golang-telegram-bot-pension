package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

const geminiAuthHeader = "x-goog-api-key"
const geminigenerateContent = "/models/gemini-2.5-flash-lite:generateContent"
const generatorSeconds = 5
const defaultTimeout = 60

type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

func GetStreamResponse(ctx context.Context, prompt string) <-chan string {
	ch := make(chan string)

	go func() {
		defer close(ch)

		for i := 0; i <= generatorSeconds; i++ {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second):
				ch <- fmt.Sprintf("token %d", i)
			}
		}
	}()

	return ch

}

func getApiKey() (string, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("missing GEMINI_API_KEY env variable")
	}
	return apiKey, nil
}

func getBaseUrl() (string, error) {
	apiKey := os.Getenv("GEMINI_BASE_URL")
	if apiKey == "" {
		return "", fmt.Errorf("missing GEMINI_BASE_URL env variable")
	}
	return apiKey, nil
}

func GenResponse(ctx context.Context, prompt string) (string, error) {
	fmt.Println("[INFO] Starting GenResponse")
	client := &http.Client{
		Timeout: time.Second * defaultTimeout,
	}
	baseURL, err := getBaseUrl()
	if err != nil {
		return "", err
	}

	url, err := url.JoinPath(baseURL, geminigenerateContent)
	if err != nil {
		return "", fmt.Errorf("failed to join URL: %w", err)
	}

	method := "POST"

	text := fmt.Sprintf(`
		Вы – помощник по проверке фактов.

		Правила:
		- Сохраняйте язык утверждения
		- Отвечайте кратко
		- Всегда отвечайте на русском языке

		Формат вывода:
		Утверждение: <claim>
		Оценка: Факт | Ложь | Вводящее в заблуждение | Мнение | Неверифицируемо
		Объяснение: <максимум 40 слов>

		Утверждение: %s
		`, prompt)

	payload := map[string]any{
		"contents": []any{
			map[string]any{
				"parts": []any{
					map[string]any{
						"text": text,
					},
				},
			},
		},
		"tools": []any{
			map[string]any{
				"google_search": map[string]any{},
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	apiKey, err := getApiKey()
	if err != nil {
		return "", err
	}
	req.Header.Set(geminiAuthHeader, apiKey)

	req = req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error %d: %s", resp.StatusCode, string(respBody))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var result GeminiResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("failed to unmarshal Gemini response: %w", err)
	}
	fmt.Println("[INFO] Response JSON parsed")

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("empty response from Gemini")
	}

	fmt.Println("[INFO] Returning final text from Gemini")
	return result.Candidates[0].Content.Parts[0].Text, nil

}
