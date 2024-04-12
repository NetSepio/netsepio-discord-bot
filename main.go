package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var BotId string
var goBot *discordgo.Session

func Validator(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BotId {
		return
	}

	rxStrict := regexp.MustCompile(`(?i)\bhttps?://\S+\b`)
	url := rxStrict.FindString(m.Content)
	if url != "" {
		fmt.Println("Found URL:", url)
		_, _ = s.ChannelMessageSend(m.ChannelID, "Hang on! NetSepio is verifying the link")

		apiUrl := "https://gateway.netsepio.com/api/v1.0/stats?siteUrl=" + url

		request, err := http.NewRequest("GET", apiUrl, nil)
		if err != nil {
			fmt.Println("Error creating GET request:", err)
			_, _ = s.ChannelMessageSend(m.ChannelID, "Error contacting NetSepio API.")
			return
		}

		client := &http.Client{Timeout: time.Second * 10}

		response, err := client.Do(request)
		if err != nil {
			fmt.Println("Error sending GET request:", err)
			_, _ = s.ChannelMessageSend(m.ChannelID, "Error contacting NetSepio API.")
			return
		}
		defer response.Body.Close()

		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			_, _ = s.ChannelMessageSend(m.ChannelID, "Error contacting NetSepio API.")
			return
		}

		var data map[string]interface{}
		err = json.Unmarshal(responseBody, &data)
		if err != nil {
			fmt.Println("Error unmarshalling response:", err)
			_, _ = s.ChannelMessageSend(m.ChannelID, "Error processing NetSepio API response.")
			return
		}

		message := fmt.Sprintf("`The link %s is not tested.`", url)

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

		if siteSafety != "" {
			message = fmt.Sprintf("`The link %s is classified as`**`%s`**", url, siteSafety)
		}

		_, _ = s.ChannelMessageSend(m.ChannelID, message)
	}
}

func Start() {
	godotenv.Load(".env")
	goBot, err := discordgo.New("Bot " + os.Getenv("Token"))
	if err != nil {
		fmt.Println(err)
		return
	}
	u, err := goBot.User("@me")
	if err != nil {
		fmt.Println(err)
		return
	}
	BotId = u.ID
	goBot.AddHandler(Validator)
	err = goBot.Open()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Bot is live")
}

func main() {

	Start()
	<-make(chan struct{})
	return
}
