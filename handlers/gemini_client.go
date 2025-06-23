package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/YahyaBekirCanevi/message_ai_go_htmx/config"
	"github.com/YahyaBekirCanevi/message_ai_go_htmx/models"
)

func GetGeminiAIResponse(prompt string) (string, error) {
	apiKey, apiURL := config.LoadConfig()
	if apiKey == "" {
		return "", fmt.Errorf("gemini_api_key not set in application.yml")
	}
	if apiURL == "" {
		return "", fmt.Errorf("gemini_api_url not set in application.yml")
	}
	url := fmt.Sprintf("%s?key=%s", apiURL, apiKey)

	requestBody := models.NewGeminiRequest(prompt)
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var geminiResp models.GeminiResponse
	err = json.Unmarshal(body, &geminiResp)
	if err != nil {
		return "", err
	}

	if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
		return geminiResp.Candidates[0].Content.Parts[0].Text, nil
	}
	return "", fmt.Errorf("no response from Gemini API")
}
