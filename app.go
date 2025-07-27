package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/sashabaranov/go-openai"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/spf13/viper"
)

type Config struct {
	TelegramToken string `mapstructure:"tgToken"`
	ChatGPTToken  string `mapstructure:"gptToken"`
	Preamble      string `mapstructure:"preamble"`
}

func LoadConfig(path string)(c Config,err error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	viper.AutomaticEnv()

	err=viper.ReadInConfig()
	if err!=nil{
		return c, fmt.Errorf("fatal error reading config file: %w", err)
	}

	err=viper.Unmarshal(&c)
	return c,err
}

func sendChatGPT(ctx context.Context, apiKey string, sendText string) (string, error) {
	client := openai.NewClient(apiKey)
	
	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: sendText,
				},
			},
			MaxTokens: 200,
		},
	)
	
	if err != nil {
		return "", fmt.Errorf("CreateChatCompletion failed: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no choices returned from OpenAI")
	}

	return resp.Choices[0].Message.Content, nil
}

func main() {
	var gptPrompt,userPrompt string

	config,err:=LoadConfig(".")
	if err!=nil{
		log.Fatalf("fatal error loading config.yaml: %v", err)
	}

	apiKey :=config.ChatGPTToken
	tgToken :=config.TelegramToken
	bot,err := tgbotapi.NewBotAPI(tgToken)
	if err!=nil{
		log.Fatalf("fatal error creating bot: %v", err)
	}

	bot.Debug=false
	log.Printf("Authorized on account : %s",bot.Self.UserName)

	u :=tgbotapi.NewUpdate(0)
	u.Timeout=10
	updates,err:=bot.GetUpdatesChan(u)
	if err!=nil{
		log.Fatalf("fatal error getting updates: %v", err)
	}

	for update := range updates {
		if update.Message==nil{
			continue
		}

		if !strings.HasPrefix(update.Message.Text,"/topic") && !strings.HasPrefix(update.Message.Text,"/phrase"){
			continue
		}

		if strings.HasPrefix(update.Message.Text,"/topic"){
			userPrompt =strings.TrimPrefix(update.Message.Text,"/topic ")
			gptPrompt=config.Preamble+"TOPIC: "
		}else if strings.HasPrefix(update.Message.Text,"/phrase"){	
			userPrompt =strings.TrimPrefix(update.Message.Text,"/phrase ")
			gptPrompt = config.Preamble+"PHRASE: "
		}

		if userPrompt != ""{
			gptPrompt += userPrompt
			response, err := sendChatGPT(context.Background(), apiKey, gptPrompt)
			if err != nil {
				log.Printf("Error from ChatGPT: %v", err)
				msgText := "Sorry, I encountered an error. Please try again later."
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
				bot.Send(msg)
				continue
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
			msg.ReplyToMessageID = update.Message.MessageID
			_, err = bot.Send(msg)
			if err != nil {
				log.Println(err)
			}
		}else{
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please enter your topic or phrase after the command.")
			bot.Send(msg)
		}
	}
}
