package controller

import (
	"github.com/valyala/fasthttp"
	"encoding/json"
	"rent-notifier/src/db"
	"rent-notifier/src/model"
	"log"
	"bytes"
	"fmt"
)

type ApiController struct {
	TelegramMessages chan model.Message
	VkMessages chan model.Message
	Db       *dbal.DBAL
	Prefix   string
}

func (controller ApiController) Notify(ctx *fasthttp.RequestCtx) error {

	ctx.SetContentType("application/json")

	body := string(ctx.PostBody())

	note := dbal.Note{}

	err := json.Unmarshal([]byte(body), &note)

	if nil != err {
		log.Printf("unmarshal error: %s", err)

		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBody([]byte(`{"status": "err"}`))

		return err
	}

	for _, recipient := range controller.Db.FindRecipientsByNote(note) {

		switch(recipient.ChatType) {
		case dbal.RECIPIENT_TELEGRAM:
			text := controller.formatMessageTelegram(note)
			controller.TelegramMessages <- model.Message{ChatId: recipient.ChatId, Text: text}
			break;
		case dbal.RECIPIENT_VK:
			text := controller.formatMessageVk(note)
			controller.VkMessages <- model.Message{ChatId: recipient.ChatId, Text: text}
		default:
			log.Printf("invalid recipient chat type: %s", recipient.ChatType)
		}
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody([]byte(`{"status": "ok"}`))

	return nil
}

func (controller ApiController) formatMessageTelegram (note dbal.Note) string {

	var b bytes.Buffer

	b.WriteString("\n******socrent.ru******\n")
	b.WriteString(fmt.Sprintf("<b>%s</b>\n", model.FormatHeader(controller.Db, note)))
	b.WriteString(fmt.Sprintf("<a href='%s'>Перейти к объявлению</a>", note.Link))
	b.WriteString("\n******socrent.ru******\n")

	return b.String()
}

func (controller ApiController) formatMessageVk (note dbal.Note) string {

	var b bytes.Buffer

	b.WriteString("\n******socrent.ru******\n")
	b.WriteString(model.FormatHeader(controller.Db, note))
	b.WriteString(note.Link)
	b.WriteString("******socrent.ru******\n")

	return b.String()
}