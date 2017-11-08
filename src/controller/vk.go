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

type VkBodyRequest struct {
	Type    string          `json:"type"`
	Object  VkObjectRequest `json:"object"`
	GroupId int             `json:"group_id"`
	Secret  string          `json:"secret"`
}

type VkObjectRequest struct {
	Date      int    `json:"date"`
	Out       int    `json:"out"`
	UserId    int    `json:"user_id"`
	ReadState int    `json:"read_state"`
	Title     string `json:"title"`
	Body      string `json:"body"`
}

type VkController struct {
	Messages chan model.Message
	Db       *dbal.DBAL
	Prefix   string
	ConfirmSecret string
}

func (controller VkController) Parse(ctx *fasthttp.RequestCtx) error {

	ctx.SetContentType("application/json")

	body := string(ctx.PostBody())

	bodyRequest := VkBodyRequest{}

	err := json.Unmarshal([]byte(body), &bodyRequest)

	if nil != err {
		log.Printf("unmarshal error: %s", err)

		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBody([]byte("ok"))

		return err
	}

	if "confirmation" == bodyRequest.Type {

		log.Print("confirmation")

		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.SetBody([]byte(controller.ConfirmSecret))

		return nil
	}

	text := []byte(strings.TrimSpace(strings.ToLower(bodyRequest.Object.Body)))
	chatId := bodyRequest.Object.UserId

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody([]byte("ok"))

	re_command_start := regexp.MustCompile(`\/start`)
	if re_command_start.Match(text) {
		controller.onStart(chatId)

		return nil
	}

	re_command_help := regexp.MustCompile(`\/help`)
	if re_command_help.Match(text) {
		controller.onHelp(chatId)

		return nil
	}

	re_command_city := regexp.MustCompile(`\/city`)
	if re_command_city.Match(text) {
		controller.onCity(chatId)

		return nil
	}

	re_subscribe := regexp.MustCompile(`\/снять`)
	if re_subscribe.Match(text) {
		controller.onSubscribe(chatId, text)

		return nil
	}

	re_unsubscribe := regexp.MustCompile(`отписаться|\/unsubscribe`)
	if re_unsubscribe.Match(text) {
		controller.onUnSubscribe(chatId)

		return nil
	}

	controller.Messages <- model.Message{ChatId: chatId, Text: "Не понимаю вас. Попробуйте обратиться за помощью: напишие <b>/help</b>"}

	log.Printf("wrong message: %s", text)

	return nil
}

func (controller VkController) onSubscribe(chatId int, byte_text []byte) {

	city := dbal.City{}
	for _, _city := range controller.Db.FindCities() {
		re := regexp.MustCompile(_city.Regexp)

		if re.Match(byte_text) {

			city = _city
			break
		}
	}

	if 0 == city.Id {
		controller.Messages <- model.Message{ChatId: chatId, Text: "Вы не указали город"}

		return
	}

	types := make([]int, 0)
	for _, _type := range controller.Db.FindTypes() {
		re := regexp.MustCompile(_type.Regexp)

		if re.Match(byte_text) {
			types = append(types, _type.Id)
		}
	}

	if 0 == len(types) {

		controller.Messages <- model.Message{ChatId: chatId, Text: "Вы не указали тип жилья"}

		return
	}

	subways := make([]int, 0)
	for _, _subway := range controller.Db.FindSubwaysByCity(city) {
		re := regexp.MustCompile(_subway.Regexp)

		if re.Match(byte_text) {
			subways = append(subways, _subway.Id)
		}
	}

	if 0 == len(subways) {
		log.Println("no subways")
	}

	for _, exists_recipient := range controller.Db.FindRecipientsByChatIdAndChatType(chatId, dbal.RECIPIENT_VK) {
		controller.Db.RemoveRecipient(exists_recipient)
	}

	recipient := dbal.Recipient{ChatId: chatId, ChatType: dbal.RECIPIENT_VK, City: city.Id, Subways: subways, Types: types}

	controller.Db.AddRecipient(recipient)

	var b bytes.Buffer

	b.WriteString("Ваша подписка успешно оформлена!\n")
	b.WriteString(fmt.Sprintf("<b>Тип</b>: %s\n", model.FormatTypes(recipient.Types)))
	b.WriteString(fmt.Sprintf("<b>Город</b>: %s\n", city.Name))
	if city.HasSubway && len(recipient.Subways) > 0 {
		b.WriteString(fmt.Sprintf("<b>Метро</b>: %s\n", model.FormatSubways(controller.Db, recipient.Subways)))
	}

	controller.Messages <- model.Message{ChatId: chatId, Text: b.String()}
}

func (controller VkController) onUnSubscribe(chat_id int) {

	for _, exists_recipient := range controller.Db.FindRecipientsByChatIdAndChatType(chat_id, dbal.RECIPIENT_VK) {
		controller.Db.RemoveRecipient(exists_recipient)
	}

	controller.Messages <- model.Message{ChatId: chat_id, Text: "Вы успешно отписаны."}
}

func (controller VkController) onStart(chat_id int) {
	var b bytes.Buffer

	b.WriteString("Добро пожаловать!\n")
	b.WriteString("<b>SocrentBot</b> предназначен для рассылки свежих объявлений жилья от собственников.\n")
	b.WriteString("Для получения рассылки напишите тип жилья, ваш город и список станций метро(если необходимо)\n")
	b.WriteString("Например: <i>Снять двушку в Москве около метро Академическая</i>\n")
	b.WriteString("Более подробная информацю о подписках: <b>/help</b>\n")
	b.WriteString("Список доступных городов: <b>/сity</b>\n")
	b.WriteString("Чтобы отписаться напишие: <b>отписаться</b> или <b>/unsubscribe</b>\n")

	controller.Messages <- model.Message{ChatId: chat_id, Text: b.String()}
}

func (controller VkController) onHelp(chat_id int) {
	var b bytes.Buffer
	b.WriteString("Для получения рассылки напишите тип жилья, ваш город и список станций метро(если необходимо)\n")
	b.WriteString("Например: <i>Снять комнату, однушку, двушку, трешку, студию в Москве около метро Академическая, Выхино, Дубровка</i>\n")
	b.WriteString("При новой подписке старая подписка удаляется\n")
	b.WriteString("Список доступных городов напишите: <b>/сity</b>\n")
	b.WriteString("Чтобы отписаться напишие: <b>отписаться</b> или <b>/unsubscribe</b>\n")

	controller.Messages <- model.Message{ChatId: chat_id, Text: b.String()}
}

func (controller VkController) onCity(chat_id int) {
	cities := make([]string, 0)
	for _, city := range controller.Db.FindCities() {
		cities = append(cities, city.Name)
	}

	controller.Messages <- model.Message{ChatId: chat_id, Text: fmt.Sprintf("Список городов:\n\n%s", strings.Join(cities, "\n"))}
}
