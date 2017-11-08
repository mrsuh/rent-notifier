package model

import (
	"net/http"
	"time"
	"net/url"
	"strconv"
	"strings"
)

type Vk struct {
	Token string
}

func (vk *Vk) SendMessage(messages chan Message) {

	for message := range messages {

		form := url.Values{}
		form.Add("user_id", strconv.Itoa(message.ChatId))
		form.Add("access_token", vk.Token)
		form.Add("v", "5.62")
		form.Add("message", message.Text)

		http.Post("https://vk.com/messages.send", "application/json", strings.NewReader(form.Encode()))

		time.Sleep(50 * time.Millisecond) //20 rps
	}
}