# Telegram Bot with ChatGPT Integration

This project is a Telegram bot written in Go that uses OpenAI's ChatGPT to generate responses for `/topic` and `/phrase` commands. Configuration is managed via a YAML file.

## Features

- Responds to `/topic` and `/phrase` commands in Telegram chats.
- Uses OpenAI's GPT-3.5 Turbo model for generating responses.
- Configurable preamble for prompts.
- Easy configuration via `config.yaml`.

## Setup

### Prerequisites

- Go 1.23.2 or higher
- Telegram Bot Token
- OpenAI API Key

### Installation

1. Clone the repository:

    ```sh
    git clone https://github.com/rajapremsai/tgbot_go.git
    cd tgbot_go
    ```

2. Install dependencies:

    ```sh
    go mod tidy
    ```

3. Configure your tokens in `config.yaml`:

    ```yaml
    tgToken: "<YOUR_TELEGRAM_BOT_TOKEN>"
    gptToken: "<YOUR_OPENAI_API_KEY>"
    preamble: "E5: "
    ```

## Usage

Run the bot:

```sh
go run app.go
```

### Commands

- `/topic <your topic>`: Get a ChatGPT-generated response for a topic.
- `/phrase <your phrase>`: Get a ChatGPT-generated response for a phrase.

## File Structure

- `app.go` - Main application logic.
- `config.yaml` - Configuration file for tokens and preamble.
- `.gitignore` - Ignores config and info files.
- `go.mod`, `go.sum` - Go module files.

## What I Learned

- How to integrate Telegram bots with Go using the `go-telegram-bot-api` library.
- How to interact with OpenAI's ChatGPT API from Go.
- Managing configuration securely with `viper` and YAML files.
- Handling user input and commands in Telegram chats.
- Structuring a Go project for clarity and maintainability.
- Error handling and logging best practices in Go.

## License

MIT

## Author

[rajapremsai](https://github.com/rajapremsai)

