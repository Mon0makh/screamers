package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot *tgbotapi.BotAPI
var chatsList []int64
var runnersSendList [50]string
var sendedNow = 0
var config map[string]string

func findIdInList(chatId int64) bool {
	var findFlag = false
	for _, e := range chatsList {
		if e == chatId {
			findFlag = true
		}
	}
	return findFlag
}

func findRunnerInList(numb string) bool {
	var findFlag = false

	for _, e := range runnersSendList {
		if e == numb {
			findFlag = true
		}
	}
	return findFlag
}

func getRunnerNumber(numb string) {
	requestURL := fmt.Sprintf("%s/numb/?numb=%s", config["http_server"], numb)
	res, err := http.Get(requestURL)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		return
	}
	log.Printf("client: status code: %d\n", res.StatusCode)

}
func initConfig() {

	file, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	contents, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(contents, &config)
	if err != nil {
		log.Fatal(err)
	}

}

func main() {
	initConfig()
	// BOT
	bot, err := tgbotapi.NewBotAPI(config["bot_token"])
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		chatId := update.Message.Chat.ID
		msg := tgbotapi.NewMessage(chatId, "")

		if findIdInList(chatId) {
			if _, err := strconv.Atoi(update.Message.Text); err == nil && len(update.Message.Text) < 5 {
				if findRunnerInList(update.Message.Text) {
					msg.Text = "Номер уже был введён ранее!"
				} else {
					runnersSendList[sendedNow] = update.Message.Text
					sendedNow++
					if sendedNow == len(runnersSendList) {
						sendedNow = 0
					}
					go getRunnerNumber(update.Message.Text)
					msg.Text = fmt.Sprintf("Номер %s отправлен", update.Message.Text)
				}
			} else {
				msg.Text = "Не корректный номер!!!"
			}
		} else {
			if update.Message.Text == config["pin_code"] {
				chatsList = append(chatsList, update.Message.Chat.ID)
				msg.Text = "Вы успешно авторизованы! Можете продолжить работу!"
			} else {
				msg.Text = "Не верный Пин-Код!"
			}
		}

		switch update.Message.Command() {
		case "start":
			if findIdInList(chatId) {
				msg.Text = "Вы уже были добавлены в список операторов!"
			} else {
				msg.Text = "Вы не авторизованны! Пожалуйста введите Пин-код: "
			}

		}

		log.Printf("[%s %d] %s", update.Message.From.FirstName, chatId, update.Message.Text)

		bot.Send(msg)

	}
}
