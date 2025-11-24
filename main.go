package main

import (
	"context"
	"log"
	"os"
	"strconv"

	translate "cloud.google.com/go/translate/apiv3"
	translatepb "cloud.google.com/go/translate/apiv3/translatepb"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, falling back to environment variables:", err)
	}
	telegramToken := os.Getenv("TELEGRAM_TOKEN")
	googleCreds := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	projectID := os.Getenv("GOOGLE_PROJECT_ID")
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

	// Google Translate client
	client, err := translate.NewTranslationClient(ctx, option.WithCredentialsFile(googleCreds))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

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

		translated, err := translateText(ctx, client, projectID, query)
		if err != nil {
			log.Println("Translation error:", err)
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

func translateText(ctx context.Context, client *translate.TranslationClient, projectID, text string) (string, error) {
	req := &translatepb.TranslateTextRequest{
		Parent:             "projects/" + projectID + "/locations/global",
		Contents:           []string{text},
		MimeType:           "text/plain",
		TargetLanguageCode: "es", // translate → Spanish
	}

	resp, err := client.TranslateText(ctx, req)
	if err != nil {
		return "", err
	}

	if len(resp.Translations) == 0 {
		return "", nil
	}

	return resp.Translations[0].TranslatedText, nil
}
