package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

type TranslationResult struct {
	TranslatedText string `json:"translated_text"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, falling back to environment variables:", err)
	}
	telegramToken := os.Getenv("TELEGRAM_TOKEN")
	geminiKey := os.Getenv("GEMINI_API_KEY")
	userIDStr := os.Getenv("TELEGRAM_USER_ID")

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		log.Fatal("Invalid TELEGRAM_USER_ID:", err)
	}

	b, err := bot.New(telegramToken)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(geminiKey))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-2.5-flash")
	model.ResponseMIMEType = "application/json"

	model.ResponseSchema = &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"translated_text": {
				Type: genai.TypeString,
			},
		},
		Required: []string{"translated_text"},
	}

	b.RegisterHandlerMatchFunc(func(update *models.Update) bool {
		return update.Message != nil && update.Message.From.ID != userID
	}, func(ctx context.Context, b *bot.Bot, update *models.Update) {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "This bot is restricted to authorized users. Please contact the bot owner, @mjmoshiri, for further assistance.",
		})
	})

	b.RegisterHandlerMatchFunc(func(update *models.Update) bool {
		return update.InlineQuery != nil && update.InlineQuery.From.ID != userID
	}, func(ctx context.Context, b *bot.Bot, update *models.Update) {
		b.AnswerInlineQuery(ctx, &bot.AnswerInlineQueryParams{
			InlineQueryID: update.InlineQuery.ID,
			Results: []models.InlineQueryResult{
				&models.InlineQueryResultArticle{
					ID:          "1",
					Title:       "Access Denied",
					Description: "This bot is restricted to authorized users. Please contact the bot owner, @mjmoshiri, for further assistance.",
					InputMessageContent: &models.InputTextMessageContent{
						MessageText: "This bot is restricted to authorized users. Please contact the bot owner, @mjmoshiri, for further assistance.",
					},
				},
			},
		})
	})

	b.RegisterHandlerMatchFunc(func(update *models.Update) bool {
		return update.InlineQuery != nil && update.InlineQuery.From.ID == userID
	}, func(ctx context.Context, b *bot.Bot, update *models.Update) {
		query := update.InlineQuery.Query
		if query == "" {
			return
		}

		translated, err := translateText(ctx, model, query)
		if err != nil {
			log.Fatal("Translation error:", err)
			translated = "Translation error"
		}

		result := &models.InlineQueryResultArticle{
			ID:                  "1",
			Title:               "Translate",
			Description:         translated,
			InputMessageContent: &models.InputTextMessageContent{MessageText: translated},
		}

		b.AnswerInlineQuery(ctx, &bot.AnswerInlineQueryParams{
			InlineQueryID: update.InlineQuery.ID,
			Results:       []models.InlineQueryResult{result},
		})
	})

	log.Println("Bot running...")
	b.Start(ctx)
}

func translateText(ctx context.Context, model *genai.GenerativeModel, text string) (string, error) {
	prompt := fmt.Sprintf("Translate the following text to Spanish: %s", text)
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}

	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		jsonString := fmt.Sprint(resp.Candidates[0].Content.Parts[0])
		var result TranslationResult
		if err := json.Unmarshal([]byte(jsonString), &result); err != nil {
			return "", fmt.Errorf("error unmarshaling JSON: %v", err)
		}
		return result.TranslatedText, nil
	}

	return "", nil
}
