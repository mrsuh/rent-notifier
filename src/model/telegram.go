package model

import (
	"net/http"
	"fmt"
	"encoding/json"
	"bytes"
)

type BodyRequest struct {
	UpdateId int         `json:"update_id"`
	Message  MessageRequest `json:"message"`
}

type MessageRequest struct {
	Chat Chat   `json:"chat"`
	Text string `json:"text"`
}

type Chat struct {
	Id int `json:"id"`
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

func Echo(bodyRequest BodyRequest, token string) {
	//https://api.telegram.org/bot%s/%s
	response := MessageResponse{ChatId: string(bodyRequest.Message.Chat.Id), Text: bodyRequest.Message.Text}

	body, err := json.Marshal(response)

	if nil != err {
		panic(err)
	}

	http.Post(fmt.Sprintf("https://api.telegram.org/bot%s/%s", token, "sendMessage"), "application/json", bytes.NewBuffer(body))
}
