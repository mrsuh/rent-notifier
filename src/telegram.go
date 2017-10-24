package main

import (
	"github.com/valyala/fasthttp"
	config "github.com/mrsuh/cli-config"
	"fmt"
	"log"
	"rent-notifier/src/db"
	"rent-notifier/src/controller"
)

func requestHandlerTelegram(ctx *fasthttp.RequestCtx, db *dbal.DBAL, token string) {
	switch string(ctx.Path()) {
	case fmt.Sprintf("/%s/webhook", token):

		if !ctx.IsPost() {
			ctx.Error("Method not allowed", fasthttp.StatusMethodNotAllowed)
			break
		}

		controller.Parse(ctx, db, token)

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

	fmt.Println("server run on ", conf["telegram.listen"].(string))

	server_api_err := fasthttp.ListenAndServe(conf["telegram.listen"].(string), func(ctx *fasthttp.RequestCtx){
		requestHandlerTelegram(ctx, db, conf["telegram.token"].(string))
	})

	if server_api_err != nil {
		log.Fatal(server_api_err)
	}
}
