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

func requestHandlerVk(ctx *fasthttp.RequestCtx, ctl controller.VkController) {

	switch string(ctx.Path()) {
	case fmt.Sprintf("/vk/%s/webhook", ctl.Prefix):

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

	conf_instance := config.GetInstance()

	err := conf_instance.Init()

	if err != nil {
		log.Fatal(err)
	}

	conf := conf_instance.Get()

	logFile, err := os.OpenFile(conf["log.file"].(string), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)
	log.SetPrefix("vk ")

	db := dbal.Connect(conf["database.dsn"].(string))

	messages := make(chan model.Message)
	vk := model.Vk{Token: conf["vk.token"].(string)}
	go vk.SendMessage(messages)

	log.Printf("server vk run on %s", conf["vk.listen"].(string))

	ctl := controller.VkController{Db: db, Messages: messages, Prefix: conf["vk.prefix"].(string), ConfirmSecret: conf["vk.confirm_secret"].(string)}

	server_err := fasthttp.ListenAndServe(conf["vk.listen"].(string), func(ctx *fasthttp.RequestCtx) {
		requestHandlerVk(ctx, ctl)
	})

	if server_err != nil {
		log.Fatal(server_err)
	}
}
