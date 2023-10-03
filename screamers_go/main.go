package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var bot *tgbotapi.BotAPI
var chatsList [15]int64
var screamersCount = 0
var sendedNow = 1
var senedPhrase = 0

var coordinatorId int64
var haveCoord = false

var phrases [9]string

var collection *mongo.Collection
var ctx = context.TODO()
var logger *log.Logger

var doneKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Готово ", "done"),
	),
)

type Runner struct {
	ID     primitive.ObjectID `bson:"_id"`
	Name   string
	number int64
}

func sendMessage(numb_id int64) {

	if screamersCount > 0 {
		if sendedNow >= screamersCount {
			sendedNow = 0
		}
		sendedNow++

		filter := bson.D{{Key: "number", Value: numb_id}}

		var result Runner
		err := collection.FindOne(context.TODO(), filter).Decode(&result)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				logger.Printf("No document %d", numb_id)
				return
			}
			logger.Printf("MONGO DB ERROR!!!: %s", err)
			log.Printf("SEND ERROR %d", err)
			panic(err)
		}
		var text_msg = ""
		if senedPhrase < 4 {
			text_msg = fmt.Sprintf("<b><u>%d</u></b>:\n %s", numb_id, fmt.Sprintf(phrases[senedPhrase], result.Name, result.Name))
			// } else if senedPhrase == 4 {
			// 	text_msg = fmt.Sprintf("<b><u>%d</u></b>:\n %s", numb_id, fmt.Sprintf(phrases[senedPhrase], result.Name, result.Name, result.Name, result.Name))
		} else if senedPhrase < 8 {
			text_msg = fmt.Sprintf("<b><u>%d</u></b>:\n %s", numb_id, fmt.Sprintf(phrases[senedPhrase], result.Name, result.Name, result.Name, result.Name, result.Name, result.Name, result.Name, result.Name))
		} else {
			text_msg = fmt.Sprintf("<b><u>%d</u></b>:\n %s", numb_id, phrases[senedPhrase])
			senedPhrase = -1
		}
		senedPhrase++

		msg := tgbotapi.NewMessage(chatsList[sendedNow], text_msg)
		msg.ReplyMarkup = doneKeyboard
		msg.ParseMode = "HTML"

		if _, err := bot.Send(msg); err != nil {
			logger.Printf("SEND ERROR %d", sendedNow)
			log.Printf("SEND ERROR %d", sendedNow)

		}

		if haveCoord {
			textMsgToCoordinator := fmt.Sprintf("Группа: %d\nБегун: %d %s", sendedNow, numb_id, result.Name)
			msg := tgbotapi.NewMessage(coordinatorId, textMsgToCoordinator)

			if _, err := bot.Send(msg); err != nil {
				log.Printf("SEND ERROR %d", sendedNow)
			}

		}

	}
}

func sendNewIdInfoMessage(newId int) {
	if screamersCount > 0 {

		text_msg := fmt.Sprintf("Ваш ID был изменён, ваш новый ID: %d", newId)
		msg := tgbotapi.NewMessage(chatsList[newId], text_msg)

		if _, err := bot.Send(msg); err != nil {
			logger.Printf("SEND ERROR %d", chatsList[newId])
			log.Printf("SEND ERROR %d", chatsList[newId])

		}

		if haveCoord {
			textMsgToCoordinator := fmt.Sprintf("Номер изменен УП изменен с %d на %d !", screamersCount, newId)
			msg := tgbotapi.NewMessage(coordinatorId, textMsgToCoordinator)

			if _, err := bot.Send(msg); err != nil {
				log.Printf("SEND ERROR %d", chatsList[newId])
			}

		}

	}
}

func numb(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "OK")
	numb_id, err := strconv.ParseInt(req.URL.Query()["numb"][0], 10, 0)
	if err != nil {
		return
	}
	go sendMessage(numb_id)

}

func headers(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func findInScreamersList(id int64) bool {
	findFlag := false
	for i := 1; i <= screamersCount; i++ {
		if chatsList[i] == id {
			findFlag = true
		}
	}
	return findFlag
}

func addToScreamersList(id int64) bool {
	findFlag := false
	for i := 1; i <= screamersCount; i++ {
		if chatsList[i] == id {
			findFlag = true
		}
	}
	if findFlag == false {
		screamersCount++
		chatsList[screamersCount] = id
	}
	logger.Printf("User added to screamers List. User id: %d  Screamers Count: %d", id, screamersCount)
	log.Printf("User added to screamers List. User id: %d  Screamers Count: %d", id, screamersCount)
	return findFlag
}

func delFromScreamersList(id int64) bool {
	findFlag := false
	idIndex := 0
	for i := 1; i <= screamersCount; i++ {
		if chatsList[i] == id {
			idIndex = i
			findFlag = true
		}
	}
	if findFlag == true {
		if idIndex < screamersCount {
			chatsList[idIndex] = chatsList[screamersCount]
			sendNewIdInfoMessage(idIndex)
		}
		screamersCount--
	}
	logger.Printf("User deleted from screamers List. User id: %d  User list id: %d Screamers Count: %d", id, idIndex, screamersCount)
	log.Printf("User deleted from screamers List. User id: %d  User list id: %d Screamers Count: %d", id, idIndex, screamersCount)
	return findFlag
}

func initPhrases() {
	phrases[0] = "Бегать быстро <b><u> %s </u></b> умеет, никого не пожалеет!\n\nЖылдам жүгіріңіз <b><u> %s </u></b> қалай біледі,ешкімге өкініш!"
	phrases[1] = "Беги всегда, беги везде — <b><u> %s </u></b>, беги к своей мечте!\n\nӘрқашан жүгіріңіз, барлық жерде жүгіріңіз <b><u> %s </u></b>, арманыңа жүгір!"
	phrases[2] = "За победу надо драться — <b><u> %s </u></b>, тебе придётся постараться!\n\nЖеңіс үшін күресу керек<b><u> %s </u></b>, тырысу керек!"
	phrases[3] = "Побеждать — твоя судьба! <b><u> %s </u></b>, победи всех как всегда!\n\nЖеңіс - сіздің тағдырыңыз!<b><u> %s </u></b>, әдеттегідей бәрін жең!"
	// phrases[4] = "Эй, <b><u> %s </u></b>, ты не на прогулке!\nДавай <b><u> %s </u></b>, шевели булками!\n\nЭй <b><u> %s </u></b>, сіз серуендеуге шыққан жоқсыз!\nКеліңіздер <b><u> %s </u></b>, орамдарды жылжытыңыз!"
	phrases[4] = "Оббеги хоть всю планету, \nБыстрее <b><u> %s </u></b> в мире нету.\n <b><u> %s </u></b> - вперед!\n <b><u> %s </u></b> - давай!\n <b><u> %s </u></b> - беги, не засыпай!\n\nБүкіл планетаны айналып жүгіріңіз\nӘлемде жылдам <b><u> %s </u></b> жоқ.\n <b><u> %s </u></b> - алға!\n <b><u> %s </u></b> - келіңіз!\n <b><u> %s </u></b> - жүгір, ұйықтап қалма!"
	phrases[5] = "Кто забыл фразу «ой, всё, не могу…»?! <b><u> %s </u></b>\nКто шепчет себе «я добегу!»?! <b><u> %s </u></b>\nКто победит на раз и два! <b><u> %s </u></b>\nВыше нос <b><u> %s </u></b>! Ура! Ура!\n\n«Ой, болды, мен алмаймын...» дегенді ұмытқан кім?! <b><u> %s </u></b>\n«Мен жүгіремін!» деп сыбырлайтын кім бар?! <b><u> %s </u></b>\nБір және екеуі кім жеңеді! <b><u> %s </u></b>\nМұрын жоғары <b><u> %s </u></b>! Ура! Ура!"
	phrases[6] = "Кто всегда вперед стремится?! <b><u> %s </u></b>\nКто летит быстрее птицы?! <b><u> %s </u></b>\nСпортивный дух у у кого в крови?! <b><u> %s </u></b>\nДавай, <b><u> %s </u></b>! Всех, всех порви!\n\nКім әрқашан алға ұмтылады? <b><u> %s </u></b>\nҚұстан да жылдам ұшатын кім?! <b><u> %s </u></b>\nБіреудің қанында спорттық рух бар ма?! <b><u> %s </u></b>\nКеліңіздер, <b><u> %s </u></b>! Барлығын құртыңдар!"
	phrases[7] = "Раз, и два — бежать пора!\nТри, четыре — лучшим в мире!\nПять и шесть — в ком сила есть!\nСемь и восемь — темп не бросим!\nДевять, десять – победу разделим все вместе!\n\nБір, екі, жүгіретін уақыт келді!\nҮш, төрт - әлемдегі ең жақсы!\nБес пен алты – кімде билік бар!\nЖеті мен сегіз - қарқыннан бас тартпайық!\nТоғыз, он – жеңісті бірге бөлісейік!"
}

func initConfig() map[string]string {

	file, err := os.Open("config.json")
	if err != nil {
		logger.Fatal(err)
		log.Fatal(err)
	}
	defer file.Close()
	contents, err := ioutil.ReadAll(file)
	if err != nil {
		logger.Fatal(err)
		log.Fatal(err)
	}
	var config map[string]string
	err = json.Unmarshal(contents, &config)
	if err != nil {
		logger.Fatal(err)
		log.Fatal(err)
	}
	return config
}

func main() {
	var err error

	// open file and create if non-existent
	file, err := os.OpenFile("custom.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	logger = log.New(file, "Custom Log", log.LstdFlags)

	config := initConfig()

	initPhrases()
	// BOT
	bot, err = tgbotapi.NewBotAPI(config["bot_token"])
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	logger.Printf("Authorized on account %s", bot.Self.UserName)
	log.Printf("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	// Mongo
	clientOptions := options.Client().ApplyURI(config["db_link"])
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logger.Fatal(err)
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		logger.Fatal(err)
		log.Fatal(err)
	}

	collection = client.Database(config["database_name"]).Collection(config["database_collection"])

	// HTTP Server
	http.HandleFunc("/numb/", numb)
	http.HandleFunc("/headers", headers)

	go http.ListenAndServe(config["http_server_port"], nil)

	for update := range updates {

		if update.Message != nil {
			logger.Printf("Update Message!")
			log.Printf("Update Message!")
			chatId := update.Message.Chat.ID
			msg := tgbotapi.NewMessage(chatId, "")

			if !findInScreamersList(chatId) {
				if chatId == coordinatorId && haveCoord {
					msg.Text = "Вы являетесь координатором! Сначало отключите координатора введя команду /stopcoord после этого вы сможете стать участником поддержки!"
				} else {
					if _, err := strconv.Atoi(update.Message.Text); err == nil && len(update.Message.Text) < 5 {
						if update.Message.Text == config["pin_code"] {
							if addToScreamersList(chatId) {
								msg.Text = "Вы уже были добавлены в список получения!"
							} else {
								msg.Text = fmt.Sprintf("Вы добавлены в список получения! Ваш номер: %d", screamersCount)
							}
						} else {

							msg.Text = "Не корректный пин-код!"
						}
					} else {
						msg.Text = "Не корректный пин-код!!!!"
					}
				}

			} else {
				msg.Text = "Вы уже были добавлены в список получения!"
			}

			switch update.Message.Command() {
			case "start":
				if findInScreamersList(chatId) {
					msg.Text = "Вы уже были добавлены в список получения!"
				} else {
					msg.Text = "Вы не авторизованны! Пожалуйста введите Пин-код: "
				}
			case "stop":
				if delFromScreamersList(chatId) {
					msg.Text = "Вы удалены из списка на получение!"
				} else {
					msg.Text = "Вы уже были удалены из списка на получение!"
				}
			case "coord":
				if findInScreamersList(chatId) {
					if haveCoord {
						textMsgToCoordinator := fmt.Sprintf("Координатор был изменён!")
						msg := tgbotapi.NewMessage(coordinatorId, textMsgToCoordinator)

						if _, err := bot.Send(msg); err != nil {
							panic(err)
						}
					}
					coordinatorId = update.Message.Chat.ID
					haveCoord = true
					delFromScreamersList(update.Message.Chat.ID)
					msg.Text = "Вы назначены координатором!"
				} else {
					msg.Text = "Вы не авторизованы!"
				}
			case "stopcoord":
				haveCoord = false
				msg.Text = "Координатор отключен!"
			}

			logger.Printf("[%s %d] %s", update.Message.From.FirstName, chatId, update.Message.Text)
			log.Printf("[%s %d] %s", update.Message.From.FirstName, chatId, update.Message.Text)

			bot.Send(msg)
		} else if update.CallbackQuery != nil {

			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)

			if _, err := bot.Request(callback); err != nil {
				panic(err)
			}

			logger.Printf("Confirm message chatId ID: %d %d", update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)

			msg := tgbotapi.NewEditMessageReplyMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, tgbotapi.InlineKeyboardMarkup{
				InlineKeyboard: make([][]tgbotapi.InlineKeyboardButton, 0),
			})
			if _, err := bot.Request(msg); err != nil {
				panic(err)
			}
		}
	}
}
