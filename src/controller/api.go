package controller

import (
	"github.com/valyala/fasthttp"
	"encoding/json"
	"rent-notifier/src/db"
	"rent-notifier/src/model"
)

func Notify(ctx *fasthttp.RequestCtx, db *dbal.DBAL , messages chan model.Message) error {

	ctx.SetContentType("application/json")

	body := string(ctx.PostBody())

	note := dbal.Note{}

	err := json.Unmarshal([]byte(body), &note)

	if nil != err {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBody([]byte(`{"status": "err"}`))

		return err
	}

	for _, recipient := range db.FindRecipientsByNote(note) {
		text := model.FormatMessage(db, note)
		messages <- model.Message{ChatId: recipient.TelegramChatId, Text: text}
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody([]byte(`{"status": "ok"}`))

	return nil
}
