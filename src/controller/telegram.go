package controller

import (
	"github.com/valyala/fasthttp"
	"encoding/json"
	"fmt"
	"rent-notifier/src/db"
	"rent-notifier/src/model"
)

func Parse(ctx *fasthttp.RequestCtx, db *dbal.DBAL, token string) bool {

	ctx.SetContentType("application/json")

	body := string(ctx.PostBody())

	fmt.Println(body)

	bodyRequest := model.BodyRequest{}

	err := json.Unmarshal([]byte(body), &bodyRequest)

	if nil != err {
		fmt.Println(err)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBody([]byte(`{"status": "err"}`))

		return false
	}

	fmt.Println(bodyRequest)

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody([]byte(`{"status": "ok"}`))

	subways := make([]int, 0)
	subways = append(subways, 1)

	recipient := dbal.Recipient{TelegramChatId: string(bodyRequest.Message.Chat.Id), City: 1, Subways: subways, Type: 1}

	model.Echo(bodyRequest, token)

	db.AddRecipient(recipient)

	return true
}
