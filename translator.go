package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type TranslationResult struct {
	TranslatedText string `json:"translated_text"`
}

var (
	genaiModel *genai.GenerativeModel
	validLangs = map[string]string{
		"en": "English",
		"es": "Spanish",
		"fr": "French",
		"de": "German",
		"it": "Italian",
		"pt": "Portuguese",
		"ru": "Russian",
		"ja": "Japanese",
		"ko": "Korean",
		"zh": "Chinese",
	}
)

func InitGemini(ctx context.Context, apiKey string) (*genai.Client, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	genaiModel = client.GenerativeModel("gemini-2.5-flash")
	genaiModel.ResponseMIMEType = "application/json"
	genaiModel.ResponseSchema = &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"translated_text": {
				Type: genai.TypeString,
			},
		},
		Required: []string{"translated_text"},
	}
	return client, nil
}

func translateText(ctx context.Context, text, targetLang string) (string, int, error) {
	prompt := fmt.Sprintf("Act as a professional translator. Translate the following text into %s. If the input contains multiple languages (e.g. 'Hola, Como was your dia?' (Target: ES) -> 'Hola, Como fue tu dia?'), translate the entire content based on the best interpretation of the context. Text: %s", validLangs[targetLang], text)
	resp, err := genaiModel.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", 0, err
	}

	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		jsonString := fmt.Sprint(resp.Candidates[0].Content.Parts[0])
		var result TranslationResult
		if err := json.Unmarshal([]byte(jsonString), &result); err != nil {
			return "", 0, fmt.Errorf("error unmarshaling JSON: %v", err)
		}

		tokens := 0
		if resp.UsageMetadata != nil {
			tokens = int(resp.UsageMetadata.TotalTokenCount)
		}
		return result.TranslatedText, tokens, nil
	}

	return "", 0, nil
}
