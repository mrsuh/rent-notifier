package model

import (
	"fmt"
	"bytes"
	"net/http"
	"time"
	"log"
)

type Telegram struct {
	Token string
}

const TELEGRAM_URL = "https://api.telegram.org/bot%s/sendMessage"

func (telegram *Telegram) SendMessage(messages chan Message) {

	for message := range messages {

		body := fmt.Sprintf(`{"chat_id": %d, "text": "%s", "parse_mode": "HTML"}`, message.ChatId, message.Text)

		resp, err := http.Post(fmt.Sprintf(TELEGRAM_URL, telegram.Token), "application/json", bytes.NewBuffer([]byte(body)))

		defer resp.Body.Close()

		if nil != err {
			log.Printf("request err {chatId: %v, err: %s}", message.ChatId, err)
		}

		time.Sleep(35 * time.Millisecond) //30 rps
	}
}
