package model

import (
	"github.com/valyala/fasthttp"
	"encoding/json"
	"fmt"
	"rent-notifier/src/db"
)

func Parse(ctx *fasthttp.RequestCtx) {

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)

	note := dbal.Note{}

	err := json.Unmarshal([]byte(ctx.PostBody()), &note)

	if nil != err {

	}

	fmt.Println(note)

	//notifier.Note{}

	ctx.SetBody([]byte(`{"status": "ok"}`))
}
