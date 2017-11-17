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

type TelegramBodyRequest struct {
	UpdateId int                    `json:"update_id"`
	Message  TelegramMessageRequest `json:"message"`
}

type TelegramMessageRequest struct {
	Chat TelegramChatRequest `json:"chat"`
	Text string              `json:"text"`
}

type TelegramChatRequest struct {
	Id int `json:"id"`
}

type TelegramMessageResponse struct {
	ChatId string `json:"chat_id"`
	Text   string `json:"text"`
}

type TelegramController struct {
	Messages chan model.Message
	DB       *dbal.DBAL
	Prefix   string
}

func (controller TelegramController) Parse(ctx *fasthttp.RequestCtx) error {

	ctx.SetContentType("application/json")

	body := string(ctx.PostBody())

	bodyRequest := TelegramBodyRequest{}

	err := json.Unmarshal([]byte(body), &bodyRequest)

	if nil != err {

		log.Printf("unmarshal error: %s", err)

		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBody([]byte(`{"status": "err"}`))

		return err
	}

	text := []byte(strings.TrimSpace(strings.ToLower(bodyRequest.Message.Text)))
	chatId := bodyRequest.Message.Chat.Id

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody([]byte(`{"status": "ok"}`))

	re_command_start := regexp.MustCompile(`^\/?start|привет`)
	if re_command_start.Match(text) {
		controller.onStart(chatId)

		return nil
	}

	re_command_help := regexp.MustCompile(`^\/?help`)
	if re_command_help.Match(text) {
		controller.onHelp(chatId)

		return nil
	}

	re_command_city := regexp.MustCompile(`^\/?city`)
	if re_command_city.Match(text) {
		controller.onCity(chatId)

		return nil
	}

	re_subscribe := regexp.MustCompile(`снять`)
	if re_subscribe.Match(text) {
		controller.onSubscribe(chatId, text)

		return nil
	}

	re_unsubscribe := regexp.MustCompile(`^\/?cancel`)
	if re_unsubscribe.Match(text) {
		controller.onUnSubscribe(chatId)

		return nil
	}

	var b bytes.Buffer
	b.WriteString("Не понимаю вас. Скорее всего вы ввели неизвествую команду.\n")
	b.WriteString("Для получения рассылки напишите: Снять <b>{тип жилья}</b> в <b>{городе}</b> около <b>{станции метро}</b>(если необходимо).\n")
	b.WriteString("\nНапример: <i>Снять комнату, однушку, двушку, студию в Питере около метро Академическая, Политехническая</i>.\n")
	b.WriteString("\nНапишите <b>help</b> для более подробной информации.\n")

	controller.Messages <- model.Message{ChatId: chatId, Text: b.String()}

	log.Printf("wrong message: %s", text)

	return nil
}

func (controller TelegramController) onSubscribe(chatId int, byte_text []byte) error {

	cities, err := controller.DB.FindCities()

	if err != nil {
		log.Printf("error find cities: %s", err)

		return err
	}

	city := dbal.City{}
	for _, _city := range cities {
		re := regexp.MustCompile(_city.Regexp)

		if re.Match(byte_text) {

			city = _city
			break
		}
	}

	if 0 == city.Id {
		var b bytes.Buffer
		b.WriteString("Вы не указали город.\n")
		b.WriteString("Для получения рассылки напишите: Снять <b>{тип жилья}</b> в <b>{городе}</b> около <b>{станции метро}</b>(если необходимо).\n")
		b.WriteString("\nНапример: <i>Снять комнату, однушку, двушку, студию в Питере около метро Академическая, Политехническая</i>.\n")
		b.WriteString("\nНапишите <b>help</b> для более подробной информации.\n")

		controller.Messages <- model.Message{ChatId: chatId, Text: b.String()}

		return nil
	}

	types := make([]int, 0)
	for _, _type := range controller.DB.FindTypes() {
		re := regexp.MustCompile(_type.Regexp)

		if re.Match(byte_text) {
			types = append(types, _type.Id)
		}
	}

	if 0 == len(types) {

		var b bytes.Buffer
		b.WriteString("Вы не указали тип жилья.\n")
		b.WriteString("Для получения рассылки напишите: Снять <b>{тип жилья}</b> в <b>{городе}</b> около <b>{станции метро}</b>(если необходимо).\n")
		b.WriteString("\nНапример: <i>Снять комнату, однушку, двушку, студию в Питере около метро Академическая, Политехническая</i>.\n")
		b.WriteString("\nНапишите <b>help</b> для более подробной информации.\n")

		controller.Messages <- model.Message{ChatId: chatId, Text: b.String()}

		return nil
	}

	subways := make([]int, 0)

	subwaysByCity, err := controller.DB.FindSubwaysByCity(city)

	if err != nil {
		log.Printf("error find subways by city: %s", err)

		return err
	}

	for _, _subway := range subwaysByCity {
		re := regexp.MustCompile(_subway.Regexp)

		if re.Match(byte_text) {
			subways = append(subways, _subway.Id)
		}
	}

	if 0 == len(subways) {
		log.Println("no subways")
	}

	exists_recipients, err := controller.DB.FindRecipientsByChatIdAndChatType(chatId, dbal.RECIPIENT_TELEGRAM)

	if err != nil {
		log.Printf("error find recipients by chat_id and chat_type: %s", err)

		return err
	}

	for _, exists_recipient := range exists_recipients {
		err := controller.DB.RemoveRecipient(exists_recipient)

		if err != nil {
			log.Printf("error remove recipient: %v %s", exists_recipient, err)

			return err
		}
	}

	recipient := dbal.Recipient{ChatId: chatId, ChatType: dbal.RECIPIENT_TELEGRAM, City: city.Id, Subways: subways, Types: types}

	err = controller.DB.AddRecipient(recipient)

	if err != nil {
		log.Printf("error add recipient: %v %s", recipient, err)

		return err
	}

	var b bytes.Buffer

	b.WriteString("Ваша подписка успешно оформлена!\n")
	b.WriteString(fmt.Sprintf("<b>Тип</b>: %s\n", model.FormatTypes(recipient.Types)))
	b.WriteString(fmt.Sprintf("<b>Город</b>: %s\n", city.Name))
	if city.HasSubway && len(recipient.Subways) > 0 {
		b.WriteString(fmt.Sprintf("<b>Метро</b>: %s\n", model.FormatSubways(controller.DB, recipient.Subways)))
	}
	b.WriteString("Вы получите новые объявления как только они появятся.\n")

	controller.Messages <- model.Message{ChatId: chatId, Text: b.String()}

	return nil
}

func (controller TelegramController) onUnSubscribe(chatId int) error {

	exists_recipients, err := controller.DB.FindRecipientsByChatIdAndChatType(chatId, dbal.RECIPIENT_TELEGRAM)

	if err != nil {
		log.Printf("error find recipients by chat_id and chat_type: %s", err)

		return err
	}

	for _, exists_recipient := range exists_recipients {
		err := controller.DB.RemoveRecipient(exists_recipient)

		if err != nil {
			log.Printf("error remove recipient: %v %s", exists_recipient, err)

			return err
		}
	}

	controller.Messages <- model.Message{ChatId: chatId, Text: "Вы успешно отменили подписку"}

	return nil
}

func (controller TelegramController) onStart(chat_id int) {
	var b bytes.Buffer

	b.WriteString("Добро пожаловать!\n")
	b.WriteString("<b>SocrentBot</b> предназначен для рассылки свежих объявлений жилья от собственников.\n")
	b.WriteString("\nДля получения рассылки напишите: Снять <b>{тип жилья}</b> в <b>{городе}</b> около <b>{станции метро}</b>(если необходимо).\n")
	b.WriteString("\nНапример: <i>Снять двушку в Питере около метро Академическая</i>.\n")
	b.WriteString("\nДополнительные команды:\n")
	b.WriteString("  <b>help</b> - более подробная информация\n")
	b.WriteString("  <b>city</b> - список доступных городов\n")
	b.WriteString("  <b>cancel</b> - отменить подписку\n")

	controller.Messages <- model.Message{ChatId: chat_id, Text: b.String()}
}

func (controller TelegramController) onHelp(chat_id int) {
	var b bytes.Buffer
	b.WriteString("<b>SocrentBot</b> предназначен для рассылки свежих объявлений жилья от собственников.\n")
	b.WriteString("Для получения рассылки напишите: Снять <b>{тип жилья}</b> в <b>{городе}</b> около <b>{станции метро}</b>(если необходимо).\n")
	b.WriteString("\nТипы жилья:\n")
	b.WriteString("  <b>комната</b>\n")
	b.WriteString("  <b>однушка</b> - 1 комнатная квартира\n")
	b.WriteString("  <b>двушка</b> - 2 комнатная квартира\n")
	b.WriteString("  <b>трешка</b> - 3 комнатная квартира\n")
	b.WriteString("  <b>студия</b> - квартира-студия\n")
	b.WriteString("  <b>квартира</b> - 1,2,3,4+ комнатная квартира, квартира-студия\n")
	b.WriteString("\nНапример:\n")
	b.WriteString("  <i>Снять комнату в Москве около метро Академическая</i>\n")
	b.WriteString("  <i>Снять однушку, студию в Питере около метро Рыбацкое</i>\n")
	b.WriteString("  <i>Снять квартиру в Волгограде</i>\n")
	b.WriteString("\nПри новой подписке старая подписка удаляется.\n")
	b.WriteString("\nДополнительные команды:\n")
	b.WriteString("  <b>city</b> - список доступных городов\n")
	b.WriteString("  <b>cancel</b> - отменить подписку\n")

	controller.Messages <- model.Message{ChatId: chat_id, Text: b.String()}
}

func (controller TelegramController) onCity(chat_id int) {
	citiesAll, err := controller.DB.FindCities()

	if err != nil {
		log.Printf("error find all cities: %s", err)

		return
	}

	cities := make([]string, 0)
	for _, city := range citiesAll {
		cities = append(cities, city.Name)
	}
	controller.Messages <- model.Message{ChatId: chat_id, Text: fmt.Sprintf("Список доступных городов:\n\n%s", strings.Join(cities, "\n"))}
}