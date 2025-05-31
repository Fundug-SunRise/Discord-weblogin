package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"

	"os"

	"github.com/joho/godotenv"
)

type discordBot struct {
	Token   string
	Session discordgo.Session
}

func initbot() (discordBot, error) {

	token, err := getToken()

	if err != nil {
		fmt.Println("Ошибка:", err)
		return discordBot{}, err
	}

	dg, err := discordgo.New("Bot " + token)

	err = dg.Open()

	bot := discordBot{
		Token:   token,
		Session: *dg,
	}

	return bot, nil
}

func getToken() (string, error) {
	err := godotenv.Load()

	if err != nil {
		return "", err
	}
	return os.Getenv("DISCORD_BOT_TOKEN"), nil
}

func (d discordBot) messageCode(userID string, code string) error {

	channel, err := d.Session.UserChannelCreate(userID)
	if err != nil {
		return fmt.Errorf("не удалось создать DM канал: %v", err)
	}

	_, err = d.Session.ChannelMessageSend(channel.ID, "Ваш код: "+code)
	if err != nil {
		return fmt.Errorf("не удалось отправить сообщение: %v", err)
	}

	return nil
}
