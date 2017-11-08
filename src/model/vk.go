package model

import (
	"net/http"
	"time"
	"net/url"
	"strconv"
	"strings"
	"log"
)

type Vk struct {
	Token string
}

func (vk *Vk) SendMessage(messages chan Message) {

	for message := range messages {

		log.Printf("request chat_id: %s, message: %s", message.ChatId, message.Text)

		form := url.Values{}
		form.Add("user_id", strconv.Itoa(message.ChatId))
		form.Add("access_token", vk.Token)
		form.Add("v", "5.62")
		form.Add("message", message.Text)

		_, err := http.Post("https://vk.com/messages.send", "application/json", strings.NewReader(form.Encode()))

		if nil != err {
			log.Printf("request err: %s", err)
		}

		time.Sleep(50 * time.Millisecond) //20 rps
	}
}