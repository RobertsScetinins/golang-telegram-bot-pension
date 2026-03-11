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
const geminigenerateContent = "/models/gemini-3-flash-preview:generateContent"
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

func getApiKey() string {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		panic("Missing GEMINI_API_KEY env variable")
	}
	return apiKey
}

func getBaseUrl() string {
	apiKey := os.Getenv("GEMINI_BASE_URL")
	if apiKey == "" {
		panic("Missing GEMINI_BASE_URL env variable")
	}
	return apiKey
}

func GenResponse(prompt string) string {
	fmt.Println("[INFO] Starting GenResponse")
	client := &http.Client{
		Timeout: time.Second * defaultTimeout,
	}

	url, err := url.JoinPath(getBaseUrl(), geminigenerateContent)

	if err != nil {
		panic("Can't create URL for request to Gemini")
	}

	method := "POST"

	text := fmt.Sprintf(`
		You are a fact-checking assistant.

		Rules:
		- Keep the language of the claim
		- Respond concisely

		Output format:
		Claim: <claim>
		Assessment: Fact | False | Misleading | Opinion | Unverifiable
		Explanation: <max 40 words>

		Claim: %s
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
	}

	body, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		panic(err)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Error creating request:", err)
		panic("Error creating request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(geminiAuthHeader, getApiKey())

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error creating request:", err)
		panic("Error creating request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		fmt.Println("API error:", resp.StatusCode, string(respBody))
		panic("API request failed")
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Read error:", err)
		panic(err)
	}

	var result GeminiResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		fmt.Println("[ERROR] Failed to parse JSON response:", err)
		panic(err)
	}
	fmt.Println("[INFO] Response JSON parsed")

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		fmt.Println("[ERROR] Empty response from Gemini")
		panic("Empty response from Gemini")
	}

	fmt.Println("[INFO] Returning final text from Gemini")
	return result.Candidates[0].Content.Parts[0].Text

}
