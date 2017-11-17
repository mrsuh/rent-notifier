package main

import (
	"github.com/valyala/fasthttp"
	"github.com/mrsuh/cli-config"
	"fmt"
	"log"
	"rent-notifier/src/db"
	"rent-notifier/src/controller"
	"rent-notifier/src/model"
	"os"
)

func requestHandlerTelegram(ctx *fasthttp.RequestCtx, ctl controller.TelegramController, connection *dbal.Connection) {

	session := connection.Session.Copy()
	defer session.Close()

	ctl.DB = &dbal.DBAL{DB: session.DB(connection.Database)}

	switch string(ctx.Path()) {
	case fmt.Sprintf("/telegram/%s/webhook", ctl.Prefix):

		if !ctx.IsPost() {
			ctx.Error("Method not allowed", fasthttp.StatusMethodNotAllowed)
			break
		}

		ctl.Parse(ctx)

		break
	default:
		ctx.Error("Not found", fasthttp.StatusNotFound)
	}
}

func main() {

	confInstance := config.GetInstance()

	err := confInstance.Init()

	if err != nil {
		log.Fatal(err)
	}

	conf := confInstance.Get()

	logFile, err := os.OpenFile(conf["log.file"].(string), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)
	log.SetPrefix("telegram ")

	connection := dbal.NewConnection(conf["database.dsn"].(string))

	messages := make(chan model.Message)
	telegram := model.Telegram{Token: conf["telegram.token"].(string)}
	go telegram.SendMessage(messages)

	log.Printf("server telegram run on %s", conf["telegram.listen"].(string))

	ctl := controller.TelegramController{Messages: messages, Prefix: conf["telegram.prefix"].(string)}

	server_err := fasthttp.ListenAndServe(conf["telegram.listen"].(string), func(ctx *fasthttp.RequestCtx) {
		requestHandlerTelegram(ctx, ctl, connection)
	})

	if server_err != nil {
		log.Fatal(server_err)
	}
}
