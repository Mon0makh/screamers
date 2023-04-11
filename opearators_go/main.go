package main

import (
	"fmt"
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
	requestURL := fmt.Sprintf("http://127.0.0.1:8090/numb/?numb=%s", numb)
	res, err := http.Get(requestURL)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		os.Exit(1)
	}

	log.Printf("client: got response!\n")
	log.Printf("client: status code: %d\n", res.StatusCode)

}

func main() {

	// BOT
	bot, err := tgbotapi.NewBotAPI("6086471341:AAHmDf8Ypi2PxHo_nM-tZ3iJRQfDWmdo0uw")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		if findIdInList(update.Message.Chat.ID) {
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
			if update.Message.Text == "11111" {
				chatsList = append(chatsList, update.Message.Chat.ID)
				msg.Text = "Вы успешно авторизованы! Можете продолжить работу!"
			} else {
				msg.Text = "Не верный Пин-Код!"
			}
		}

		switch update.Message.Command() {
		case "start":
			if findIdInList(update.Message.Chat.ID) {
				msg.Text = "Вы уже были добавлены в список операторов!"
			} else {
				msg.Text = "Вы не авторизованны! Пожалуйста введите Пин-код: "
			}

		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		bot.Send(msg)

	}
}
