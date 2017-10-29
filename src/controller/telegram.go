package controller

import (
	"github.com/valyala/fasthttp"
	"encoding/json"
	"fmt"
	"rent-notifier/src/db"
	"strings"
	"regexp"
	"rent-notifier/src/model"
	"log"
)

type BodyRequest struct {
	UpdateId int            `json:"update_id"`
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

func Parse(ctx *fasthttp.RequestCtx, db *dbal.DBAL, messages chan model.Message) error {

	ctx.SetContentType("application/json")

	body := string(ctx.PostBody())

	bodyRequest := BodyRequest{}

	err := json.Unmarshal([]byte(body), &bodyRequest)

	if nil != err {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBody([]byte(`{"status": "err"}`))

		return err
	}

	text := []byte(strings.ToLower(bodyRequest.Message.Text))
	chatId := bodyRequest.Message.Chat.Id

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody([]byte(`{"status": "ok"}`))

	re_command_start := regexp.MustCompile(`\/start`)
	if re_command_start.Match(text) {
		onStart(chatId, messages)

		return nil
	}

	re_command_help := regexp.MustCompile(`\/help`)
	if re_command_help.Match(text) {
		onHelp(chatId, messages)

		return nil
	}

	re_subscribe := regexp.MustCompile(`хочу|снять`)
	if re_subscribe.Match(text) {
		onSubscribe(db, text, chatId, messages)

		return nil
	}

	re_unsubscribe := regexp.MustCompile(`отписаться|\/unsubscribe`)
	if re_unsubscribe.Match(text) {
		onUnSubscribe(chatId, messages)

		return nil
	}

	log.Print("wrong message", text)

	return nil
}

func onSubscribe(db *dbal.DBAL, byte_text []byte, chat_id int, messages chan model.Message) {

	city := dbal.City{}
	for _, _city := range db.FindCities() {
		re := regexp.MustCompile(_city.Regexp)

		if re.Match(byte_text) {

			city = _city
			break
		}
	}

	if 0 == city.Id {
		//todo
	}

	types := make([]int, 0)
	for _, _type := range db.FindTypes() {
		re := regexp.MustCompile(_type.Regexp)

		if re.Match(byte_text) {
			types = append(types, _type.Id)
		}
	}

	if 0 == len(types) {
		fmt.Println("no types")
		//todo
	}

	subways := make([]int, 0)
	for _, _subway := range db.FindSubwaysByCity(city) {
		re := regexp.MustCompile(_subway.Regexp)

		if re.Match(byte_text) {
			subways = append(subways, _subway.Id)
		}
	}

	if 0 == len(subways) {
		fmt.Println("no subways")
		//todo
	}

	//find recipient
	//if exists -> send Message "already subscribed"
	//create recipient if not exists

	recipient := dbal.Recipient{TelegramChatId: chat_id, City: city.Id, Subways: subways, Types: types}

	db.AddRecipient(recipient)

	//send message  "Вы успешно подписаны на ..."

	messages <- model.Message{ChatId:chat_id, Text: "Вы успешно подписаны на ..."}
}

func onUnSubscribe(chat_id int, messages chan model.Message) {

	//recipients := db.FindRecipientsByChatId(bodyRequest.Message.Chat.Id)
	//
	//for _, recipient := range recipients {
	//	db.RemoveRecipient(recipient)
	//}

	//send message "вы успешно отписаны"

	messages <- model.Message{ChatId:chat_id, Text: "Вы успешно отписаны"}
}

func onStart(chat_id int, messages chan model.Message) {
	messages <- model.Message{ChatId:chat_id, Text: "Добро пожаловать..."}
}

func onHelp(chat_id int, messages chan model.Message){
	messages <- model.Message{ChatId:chat_id, Text: "Список городов, как подписаться, как отписаться"}
}