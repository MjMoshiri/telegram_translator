# Telegram Translator Bot

A simple Telegram bot that translates inline queries into Spanish using Google Gemini API.
 
 ## Setup
 
 1.  **Clone the repository.**
 2.  **Install dependencies:**
     ```bash
     go mod download
     ```
 3.  **Configure Environment Variables:**
     Create a `.env` file in the root directory with the following:
     ```env
     GEMINI_API_KEY=your-gemini-api-key
     TELEGRAM_TOKEN=your-telegram-bot-token
     TELEGRAM_USER_ID=your-allowed-user-id
     ```
 4.  **Gemini API Key:**
     Obtain an API key from Google AI Studio and set it in `.env`.

## Running

```bash
go run main.go
```

## Usage

In Telegram, type `@your_bot_username <text>` to translate text to Spanish. The bot is restricted to respond only to the user ID specified in `TELEGRAM_USER_ID`.
