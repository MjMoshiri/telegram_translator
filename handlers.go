package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func handleAddUser(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message.From.ID != OwnerID {
		return
	}
	parts := strings.Split(update.Message.Text, " ")
	if len(parts) != 2 {
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "Usage: /adduser <id>"})
		return
	}
	id, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "Invalid ID"})
		return
	}
	if err := AddUser(id); err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "Error adding user"})
		return
	}
	b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: fmt.Sprintf("User %d added successfully. Please instruct them to contact @%s to set their translation language.", id, BotUsername)})
}

func handleRemoveUser(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message.From.ID != OwnerID {
		return
	}
	parts := strings.Split(update.Message.Text, " ")
	if len(parts) != 2 {
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "Usage: /removeuser <id>"})
		return
	}
	id, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "Invalid ID"})
		return
	}
	if err := RemoveUser(id); err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "Error removing user"})
		return
	}
	b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "User removed"})
}

func handleLanguage(ctx context.Context, b *bot.Bot, update *models.Update) {
	userID := update.Message.From.ID
	if !IsUserAllowed(userID) {
		return
	}
	parts := strings.Split(update.Message.Text, " ")
	if len(parts) != 2 {
		langs := "Usage: /language <code>\nSupported languages:\n"
		for code, name := range validLangs {
			langs += fmt.Sprintf("%s: %s\n", code, name)
		}
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: langs})
		return
	}
	lang := parts[1]
	if _, ok := validLangs[lang]; !ok {
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "Invalid language code"})
		return
	}
	if err := SetUserLanguage(userID, lang); err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "Error setting language"})
		return
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf("Language set to %s. You can now use the bot in any chat by typing @%s <text to translate> to get an inline translation.", validLangs[lang], BotUsername),
	})
}

func handleInlineQuery(ctx context.Context, b *bot.Bot, update *models.Update) {
	userID := update.InlineQuery.From.ID
	if !IsUserAllowed(userID) {
		b.AnswerInlineQuery(ctx, &bot.AnswerInlineQueryParams{
			InlineQueryID: update.InlineQuery.ID,
			Results: []models.InlineQueryResult{
				&models.InlineQueryResultArticle{
					ID:          "1",
					Title:       "Access Denied",
					Description: "You are not authorized to use this bot.",
					InputMessageContent: &models.InputTextMessageContent{
						MessageText: "You are not authorized to use this bot.",
					},
				},
			},
		})
		return
	}

	query := update.InlineQuery.Query
	if query == "" {
		return
	}

	targetLang := GetUserLanguage(userID)
	translated, tokens, err := translateText(ctx, query, targetLang)
	if err != nil {
		log.Println("Translation error:", err)
		translated = "Translation error"
	} else {
		logUsage(userID, tokens)
	}

	result := &models.InlineQueryResultArticle{
		ID:                  "1",
		Title:               "Translate to " + validLangs[targetLang],
		Description:         translated,
		InputMessageContent: &models.InputTextMessageContent{MessageText: translated},
	}

	b.AnswerInlineQuery(ctx, &bot.AnswerInlineQueryParams{
		InlineQueryID: update.InlineQuery.ID,
		Results:       []models.InlineQueryResult{result},
	})
}

func handleHelp(ctx context.Context, b *bot.Bot, update *models.Update) {
	userID := update.Message.From.ID
	if !IsUserAllowed(userID) {
		b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "You are not authorized to use this bot."})
		return
	}

	helpText := "Congratulations on being authorized to use the bot!\n\n" +
		"The next step is to set your preferred language.\n" +
		"You can set or reset the language using the /language command.\n\n" +
		"Usage: /language <code>\n" +
		"Example: /language es\n\n" +
		"You can also use mixed languages in your input, and the bot will translate it based on context.\n" +
		"Example: 'Hola, Como was your dia?'\n\n" +
		"Supported languages:\n"

	for code, name := range validLangs {
		helpText += fmt.Sprintf("%s: %s\n", code, name)
	}

	b.SendMessage(ctx, &bot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: helpText})
}
