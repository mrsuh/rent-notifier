package controller

import (
	"github.com/valyala/fasthttp"
	"encoding/json"
	"rent-notifier/src/db"
	"rent-notifier/src/model"
	"log"
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
		text := model.FormatMessage(controller.Db, note)

		message := model.Message{ChatId: recipient.ChatId, Text: text}

		switch(recipient.ChatType) {
		case dbal.RECIPIENT_TELEGRAM:
			controller.TelegramMessages <- message
			break;
		case dbal.RECIPIENT_VK:
			controller.VkMessages <- message
		default:
			log.Printf("invalid recipient chat type: %s", recipient.ChatType)
		}
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody([]byte(`{"status": "ok"}`))

	return nil
}
