package model

import (
	"net/http"
	"fmt"
	"encoding/json"
	"bytes"
)

type MessageRequest struct {
	Chat Chat   `json:"chat"`
	Text string `json:"text"`
}

type Chat struct {
	Id string `json:"id"`
}

type MessageResponse struct {
	ChatId string `json:"chat_id"`
	Text   string `json:"text"`
}

func AddRecipient() {

}

func RemoveRecipient() {

}

func Notify() {

}

func Echo(message MessageRequest, token string) {
	//https://api.telegram.org/bot%s/%s
	response := MessageResponse{ChatId: message.Chat.Id, Text: message.Text}

	body, err := json.Marshal(response)

	if nil != err {
		panic(err)
	}

	http.Post(fmt.Sprintf("https://api.telegram.org/bot%s/%s", token, "sendMessage"), "application/json", bytes.NewBuffer(body))
}
