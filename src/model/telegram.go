package model

import (
	"fmt"
	"bytes"
	"net/http"
)

type Message struct {
	ChatId int
	Text   string
}

type Telegram struct {
	Token string
}

func (telegram *Telegram) SendMessage(messages chan Message) {

	for message := range messages {

		body := fmt.Sprintf(`{"chat_id": %d, "text": "%s"}`, message.ChatId, message.Text)

		http.Post(fmt.Sprintf("https://api.telegram.org/bot%s/%s", telegram.Token, "sendMessage"), "application/json", bytes.NewBuffer([]byte(body)))

		fmt.Println(body)
	}
}