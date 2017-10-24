package main

import (
	"github.com/valyala/fasthttp"
	config "github.com/mrsuh/cli-config"
	"fmt"
	"log"
	"rent-notifier/src/controller"
)

func requestHandlerApi(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {

	case "/api/v1/notify":

		if !ctx.IsPost() {
			ctx.Error("Method not allowed", fasthttp.StatusMethodNotAllowed)
			break
		}

		controller.Notify(ctx)

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

	fmt.Println("server run on ", conf["api.listen"].(string))

	server_api_err := fasthttp.ListenAndServe(conf["api.listen"].(string), requestHandlerApi)

	if server_api_err != nil {
		log.Fatal(server_api_err)
	}
}
