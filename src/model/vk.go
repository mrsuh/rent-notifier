package model

import (
	"net/http"
	"time"
	"net/url"
	"strconv"
	"strings"
	"log"
	"io/ioutil"
	"encoding/json"
	"rent-notifier/src/controller"
	"rent-notifier/src/db"
)

type Vk struct {
	Token      string
	Connection *dbal.Connection
}

const VK_URL = "https://api.vk.com/method/messages.send"

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

		resp, err := http.Post(VK_URL, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))

		defer resp.Body.Close()

		if nil != err {
			log.Printf("request err {form: %v, err: %s}", form, err)

			continue
		}

		bodyBytes, _ := ioutil.ReadAll(resp.Body)

		bodyResponse := controller.VkBodyResponse{}
		err = json.Unmarshal(bodyBytes, &bodyResponse)

		if err == nil && bodyResponse.Error.Code == 901 {

			vk.RemoveInvalidRecipients(message)
		}

		log.Printf("response: %s", string(bodyBytes))

		time.Sleep(50 * time.Millisecond) //20 rps
	}
}

func (vk *Vk) RemoveInvalidRecipients(message Message) {
	recipientIds := make([]int, 0)
	if message.IsBulk {
		recipientIds = message.ChatIds
	} else {
		recipientIds = append(recipientIds, message.ChatId)
	}

	session := vk.Connection.Session.Copy()
	defer session.Close()

	db := &dbal.DBAL{DB: session.DB(vk.Connection.Database)}

	for chatId := range recipientIds {
		form := url.Values{}

		form.Add("user_id", strconv.Itoa(chatId))

		form.Add("access_token", vk.Token)
		form.Add("v", "5.64")
		form.Add("message", message.Text)

		resp, err := http.Post(VK_URL, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
		defer resp.Body.Close()

		if nil != err {
			log.Printf("request err {form: %v, err: %s}", form, err)

			continue
		}

		bodyBytes, _ := ioutil.ReadAll(resp.Body)

		bodyResponse := controller.VkBodyResponse{}
		err = json.Unmarshal(bodyBytes, &bodyResponse)

		if err == nil && bodyResponse.Error.Code == 901 {
			err = db.RemoveRecipient(dbal.Recipient{ChatId: chatId, ChatType: dbal.RECIPIENT_VK})

			if err != nil {
				log.Printf("remove recipient err: %s", err)
			}
		}

		log.Printf("response: %s", string(bodyBytes))
	}
}
