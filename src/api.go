package main

import (
	"github.com/valyala/fasthttp"
	"github.com/mrsuh/cli-config"
	"fmt"
	"log"
	"rent-notifier/src/db"
	"rent-notifier/src/model"
	"rent-notifier/src/controller"
	"os"
)

func requestHandlerApi(ctx *fasthttp.RequestCtx, ctl controller.ApiController, connection *dbal.Connection) {

	session := connection.Session.Copy()
	defer session.Close()

	ctl.DB = &dbal.DBAL{DB: session.DB(connection.Database)}

	switch string(ctx.Path()) {

	case fmt.Sprintf("/%s/notify", ctl.Prefix):

		if !ctx.IsPost() {
			ctx.Error("Method not allowed", fasthttp.StatusMethodNotAllowed)
			break
		}

		ctl.Notify(ctx)

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

	connection := dbal.NewConnection(conf["database.dsn"].(string))

	logFile, err := os.OpenFile(conf["log.file"].(string), os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)
	log.SetPrefix("api ")

	telegramMessages := make(chan model.Message)
	telegram := model.Telegram{Token: conf["telegram.token"].(string)}
	go telegram.SendMessage(telegramMessages)

	vkMessages := make(chan model.Message)
	vk := model.Vk{Token: conf["vk.token"].(string)}
	go vk.SendMessage(vkMessages)

	log.Printf("server run on %s", conf["api.listen"].(string))

	ctl := controller.ApiController{TelegramMessages: telegramMessages, VkMessages: vkMessages, Prefix: conf["api.prefix"].(string)}

	serverErr := fasthttp.ListenAndServe(conf["api.listen"].(string), func(ctx *fasthttp.RequestCtx) {
		requestHandlerApi(ctx, ctl, connection)
	})

	if serverErr != nil {
		log.Fatal(serverErr)
	}
}
