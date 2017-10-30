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
	"bytes"
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

	re_command_city := regexp.MustCompile(`\/city`)
	if re_command_city.Match(text) {
		onCity(db, chatId, messages)

		return nil
	}

	re_subscribe := regexp.MustCompile(`хочу|снять`)
	if re_subscribe.Match(text) {
		onSubscribe(db, text, chatId, messages)

		return nil
	}

	re_unsubscribe := regexp.MustCompile(`отписаться|\/unsubscribe`)
	if re_unsubscribe.Match(text) {
		onUnSubscribe(db, chatId, messages)

		return nil
	}

	messages <- model.Message{ChatId: chatId, Text: "Не понимаю вас. Попробуйте обратиться за помощью: <b>/help</b>"}

	log.Print("wrong message", text)

	return nil
}

func onSubscribe(db *dbal.DBAL, byte_text []byte, chatId int, messages chan model.Message) {

	city := dbal.City{}
	for _, _city := range db.FindCities() {
		re := regexp.MustCompile(_city.Regexp)

		if re.Match(byte_text) {

			city = _city
			break
		}
	}

	if 0 == city.Id {
		messages <- model.Message{ChatId: chatId, Text: "Вы не указали город"}

		return
	}

	types := make([]int, 0)
	for _, _type := range db.FindTypes() {
		re := regexp.MustCompile(_type.Regexp)

		if re.Match(byte_text) {
			types = append(types, _type.Id)
		}
	}

	if 0 == len(types) {

		messages <- model.Message{ChatId: chatId, Text: "Вы не указали тип жилья"}

		return
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
	}

	for _, exists_recipient := range db.FindRecipientsByChatId(chatId) {
		db.RemoveRecipient(exists_recipient)
	}

	recipient := dbal.Recipient{TelegramChatId: chatId, City: city.Id, Subways: subways, Types: types}

	db.AddRecipient(recipient)

	var b bytes.Buffer

	b.WriteString("Ваша подписка успешно оформлена!\n")
	b.WriteString(fmt.Sprintf("Тип: %s\n", model.FormatTypes(recipient.Types)))
	b.WriteString(fmt.Sprintf("Город: %s\n", city.Name))
	if city.HasSubway && len(recipient.Subways) > 0 {
		b.WriteString(fmt.Sprintf("Метро: %s\n", model.FormatSubways(db, recipient.Subways)))
	}

	messages <- model.Message{ChatId: chatId, Text: b.String()}
}

func onUnSubscribe(db *dbal.DBAL, chat_id int, messages chan model.Message) {

	for _, exists_recipient := range db.FindRecipientsByChatId(chat_id) {
		db.RemoveRecipient(exists_recipient)
	}

	messages <- model.Message{ChatId:chat_id, Text: "Вы успешно отписаны."}
}

func onStart(chat_id int, messages chan model.Message) {
	var b bytes.Buffer

	b.WriteString("Добро пожаловать!\n")
	b.WriteString("<b>SocrentBot</b> предназначен для рассылки свежих объявлений жилья от собственников.\n")
	b.WriteString("Для получения рассылки напишите тип жилья, ваш город и список станций метро(если необходимо)\n")
	b.WriteString("Например: <i>Снять комнату в Москве около метро Академическая</i>\n")
	b.WriteString("Чтобы получить более подробную информацю о подписках напишие: <b>/help</b>\n")
	b.WriteString("Чтобы получить список доступных городов напишите: <b>/сity</b>\n")
	b.WriteString("Чтобы отписаться напишие: <b>отписаться</b> или <b>/unsubscribe</b>\n")

	messages <- model.Message{ChatId:chat_id, Text: b.String()}
}

func onHelp(chat_id int, messages chan model.Message){
	var b bytes.Buffer
	b.WriteString("Для получения рассылки напишите тип жилья, ваш город и список станций метро(если необходимо)\n")
	b.WriteString("Например: <i>Снять комнату, однушку, двушку, трешку, студию в Москве около метро Академическая, Выхино, Дубровка</i>\n")
	b.WriteString("При новой подписке старая подписка удаляется\n")
	b.WriteString("Чтобы получить список доступных городов напишите: <b>/сity</b>\n")
	b.WriteString("Чтобы отписаться напишие: <b>отписаться</b> или <b>/unsubscribe</b>\n")

	messages <- model.Message{ChatId:chat_id, Text: b.String()}
}

func onCity(db *dbal.DBAL, chat_id int, messages chan model.Message) {
	cities := make([]string, 0)
	for _, city := range db.FindCities() {
		cities = append(cities, city.Name)
	}
	messages <- model.Message{ChatId:chat_id, Text: fmt.Sprintf("Список городов:\n %s", strings.Join(cities, "\n"))}
}