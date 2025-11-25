package main

import (
	"context"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func main() {
	LoadConfig()

	InitDB()

	b, err := bot.New(TelegramToken)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	client, err := InitGemini(ctx, GeminiKey)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	b.SetMyCommands(ctx, &bot.SetMyCommandsParams{
		Commands: []models.BotCommand{
			{Command: "language", Description: "Set target language"},
			{Command: "help", Description: "Show help text"},
		},
	})

	// Command Handlers
	b.RegisterHandler(bot.HandlerTypeMessageText, "/adduser", bot.MatchTypePrefix, handleAddUser)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/removeuser", bot.MatchTypePrefix, handleRemoveUser)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/language", bot.MatchTypePrefix, handleLanguage)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, handleHelp)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, handleHelp)

	// Inline Query Handler
	b.RegisterHandlerMatchFunc(func(update *models.Update) bool {
		return update.InlineQuery != nil
	}, handleInlineQuery)

	// Start flush routine
	go flushUsageRoutine()

	log.Println("Bot running...")
	b.Start(ctx)
}
