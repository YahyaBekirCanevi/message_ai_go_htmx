package models

// GeminiRequest is a struct for the Gemini API request body
// NewGeminiRequest creates a GeminiRequest with the given prompt text

type GeminiRequest struct {
	Contents []GeminiContent `json:"contents"`
}

type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
}

type GeminiPart struct {
	Text string `json:"text"`
}

func NewGeminiRequest(prompt string) GeminiRequest {
	return GeminiRequest{
		Contents: []GeminiContent{
			{
				Parts: []GeminiPart{
					{Text: prompt},
				},
			},
		},
	}
}

func NewGeminiRequestParts(messages []Message, userPrompt string) GeminiRequest {
	// Limit the number of previous messages sent for context
	const maxContextMessages = 10
	startIdx := 0
	if len(messages) > maxContextMessages {
		startIdx = len(messages) - maxContextMessages
	}
	trimmedMessages := messages[startIdx:]

	// Build context from previous messages
	var parts []GeminiPart
	
	for _, msg := range trimmedMessages {
		role := "User: "
		if msg.Sender == "ai" {
			role = "AI: "
		}
		parts = append(parts, GeminiPart{Text: role + msg.Content})
	}
	parts = append(parts, GeminiPart{Text: "User: " + userPrompt})
	return GeminiRequest{
		Contents: []GeminiContent{
			{
				Parts: parts,
			},
		},
	}
}
// GeminiResponse is a struct for the Gemini API response body

type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}
