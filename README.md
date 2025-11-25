# Telegram Translator Bot

A Telegram bot that translates inline queries using Google Gemini API. It supports multiple languages and user management.

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
    OWNER_ID=your-telegram-user-id
    BOT_USERNAME=your-telegram-bot-username
    ```
4.  **Gemini API Key:**
    Obtain an API key from Google AI Studio and set it in `.env`.

## Running

```bash
go run .
```

## Usage

### Inline Translation
In any chat, type `@your_bot_username <text>` to translate text. The bot will offer the translation which you can click to send.

### Commands
- `/language <code>`: Set your preferred target language (e.g., `es`, `fr`, `de`).
- Supported languages: English (en), Spanish (es), French (fr), German (de), Italian (it), Portuguese (pt), Russian (ru), Japanese (ja), Korean (ko), Chinese (zh).

### Admin Commands
Only the owner can manage users.
- `/adduser <id>`: Allow a user to use the bot.
- `/removeuser <id>`: Revoke access for a user.

The bot is restricted to the `OWNER_ID` and explicitly allowed users.

## Demo

![Sample Run](sample_run.gif)
