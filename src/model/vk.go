package model

import (
	"net/http"
	"time"
	"net/url"
	"strconv"
	"strings"
	"log"
	"io/ioutil"
)

type Vk struct {
	Token string
}

func (vk *Vk) SendMessage(messages chan Message) {

	for message := range messages {

		form := url.Values{}

		if message.IsBulk {

			vkIds := make([]string, 0)
			for _, id := range message.ChatIds {
				vkIds = append(vkIds, strconv.Itoa(id))
			}

			form.Add("user_ids", strings.Join(vkIds, ","))
		} else {
			form.Add("user_id", strconv.Itoa(message.ChatId))
		}

		form.Add("access_token", vk.Token)
		form.Add("v", "5.64")
		form.Add("message", message.Text)

		resp, err := http.Post("https://api.vk.com/method/messages.send", "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))

		defer resp.Body.Close()

		if nil != err {
			log.Printf("request err: %s", err)
		}

		bodyBytes,_ := ioutil.ReadAll(resp.Body)
		log.Printf("response: %s", string(bodyBytes))

		time.Sleep(50 * time.Millisecond) //20 rps
	}
}