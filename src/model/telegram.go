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

func (telegram *Telegram) SendMessage(messages chan Message) {

	for message := range messages {

		body := fmt.Sprintf(`{"chat_id": %d, "text": "%s", "parse_mode": "HTML"}`, message.ChatId, message.Text)

		_, err := http.Post(fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", telegram.Token), "application/json", bytes.NewBuffer([]byte(body)))

		if nil != err {
			log.Printf("request err: %s", err)
		}

		time.Sleep(35 * time.Millisecond) //30 rps
	}
}