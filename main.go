package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/joho/godotenv"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

var Bot *tgbotapi.BotAPI

func Validator(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	rxStrict := regexp.MustCompile(`(?i)\bhttps?://\S+\b`)
	url := rxStrict.FindString(update.Message.Text)
	if url != "" {
		fmt.Println("Found URL:", url)

		apiUrl := "https://gateway.netsepio.com/api/v1.0/stats?siteUrl=" + url

		client := &http.Client{Timeout: time.Second * 10}

		response, err := client.Get(apiUrl)
		if err != nil {
			fmt.Println("Error sending GET request:", err)
			return
		}
		defer response.Body.Close()

		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return
		}

		var data map[string]interface{}
		err = json.Unmarshal(responseBody, &data)
		if err != nil {
			fmt.Println("Error unmarshalling response:", err)
			return
		}

		message := fmt.Sprintf("The link %s is not tested.", url)

		var siteSafety string
		if safety, ok := data["payload"]; ok {
			if safetyList, ok := safety.([]interface{}); ok {
				if len(safetyList) > 0 {
					if safetyMap, ok := safetyList[0].(map[string]interface{}); ok {
						if safetyValue, ok := safetyMap["siteSafety"].(string); ok {
							siteSafety = safetyValue
						}
					}
				}
			}
		}

		if siteSafety != "" && siteSafety != "Safe" {
			message = fmt.Sprintf("The link %s is classified as %s", url, siteSafety)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
			_, err := Bot.Send(msg)
			if err != nil {
				fmt.Println("Error sending message:", err)
			}
		}
	}
}

func Start() {
	godotenv.Load(".env")
	var err error
	Bot, err = tgbotapi.NewBotAPI(os.Getenv("TelegramToken"))
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("Bot started")
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := Bot.GetUpdatesChan(u)
	if err != nil {
		fmt.Println(err)
		return
	}

	for update := range updates {
		Validator(update)
	}
}

func main() {
	Start()
}
