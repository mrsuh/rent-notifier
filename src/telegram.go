package main

import (
	"github.com/valyala/fasthttp"
	"github.com/mrsuh/cli-config"
	"fmt"
	"log"
	"rent-notifier/src/db"
	"rent-notifier/src/controller"
	"rent-notifier/src/model"
)

func requestHandlerTelegram(ctx *fasthttp.RequestCtx, db *dbal.DBAL, token string, messages chan model.Message) {
	switch string(ctx.Path()) {
	case fmt.Sprintf("/%s/webhook", token):

		if !ctx.IsPost() {
			ctx.Error("Method not allowed", fasthttp.StatusMethodNotAllowed)
			break
		}

		err := controller.Parse(ctx, db, messages)

		if err != nil {
			log.Println("error", err)
		}

		break
	default:
		ctx.Error("Not found", fasthttp.StatusNotFound)
	}
}

func main() {

	conf_instance := config.GetInstance()

	err := conf_instance.Init()

	if err != nil {
		log.Fatal(err)
	}

	conf := conf_instance.Get()

	db := dbal.Connect(conf["database.dsn"].(string))

	messages := make(chan model.Message)
	telegram := model.Telegram{Token: conf["telegram.token"].(string)}
	go telegram.SendMessage(messages)

	fmt.Println("server telegram run on ", conf["telegram.listen"].(string))

	server_telegram_err := fasthttp.ListenAndServe(conf["telegram.listen"].(string), func(ctx *fasthttp.RequestCtx) {
		requestHandlerTelegram(ctx, db, conf["telegram.token"].(string), messages)
	})

	if server_telegram_err != nil {
		log.Fatal(server_telegram_err)
	}
}
