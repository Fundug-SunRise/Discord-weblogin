package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type discordBot struct {
	Token   string
	Session discordgo.Session
}

func createBot(token string) (discordBot, error) {

	dg, err := discordgo.New("Bot " + token)

	err = dg.Open()
	if err != nil {
		fmt.Println("Ошибка:", err)
		return discordBot{}, err
	}

	bot := discordBot{
		Token:   token,
		Session: *dg,
	}

	return bot, nil
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
