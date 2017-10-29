package main

import (
	"github.com/valyala/fasthttp"
	config "github.com/mrsuh/cli-config"
	"fmt"
	"log"
	"rent-notifier/src/db"
	"rent-notifier/src/model"
	"rent-notifier/src/controller"
)

func requestHandlerApi(ctx *fasthttp.RequestCtx, db *dbal.DBAL, messages chan model.Message) {
	switch string(ctx.Path()) {

	case "/api/v1/notify":

		if !ctx.IsPost() {
			ctx.Error("Method not allowed", fasthttp.StatusMethodNotAllowed)
			break
		}

		err := controller.Notify(ctx, db, messages)

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

	fmt.Println("server run on ", conf["api.listen"].(string))

	server_api_err := fasthttp.ListenAndServe(conf["api.listen"].(string), func(ctx *fasthttp.RequestCtx) {
		requestHandlerApi(ctx, db, messages)
	})

	if server_api_err != nil {
		log.Fatal(server_api_err)
	}
}
