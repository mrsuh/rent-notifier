package controller

import (
	"github.com/valyala/fasthttp"
	"encoding/json"
	"fmt"
	"rent-notifier/src/db"
)

func Notify(ctx *fasthttp.RequestCtx) bool {

	// get chat_ids by city, subway and type
	// send notification to chat_ids by goroutine

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)

	body := string(ctx.PostBody())

	note := dbal.Note{}

	err := json.Unmarshal([]byte(body), &note)

	if nil != err {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBody([]byte(`{"status": "err"}`))

		return false
	}

	fmt.Println(note)

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody([]byte(`{"status": "ok"}`))

	return true
}
