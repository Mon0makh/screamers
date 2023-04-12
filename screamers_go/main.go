package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
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

var phrases [10]string

var collection *mongo.Collection
var ctx = context.TODO()

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

		filter := bson.D{{Key: "number", Value: numb_id}}

		var result Runner
		err := collection.FindOne(context.TODO(), filter).Decode(&result)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				log.Printf("No document %d", numb_id)
				return
			}
			panic(err)
		}
		var text_msg = ""
		if senedPhrase < 4 {
			text_msg = fmt.Sprintf("<b><u>%d</u></b>:\n %s", numb_id, fmt.Sprintf(phrases[senedPhrase], result.Name, result.Name))
		} else if senedPhrase == 4 {
			text_msg = fmt.Sprintf("<b><u>%d</u></b>:\n %s", numb_id, fmt.Sprintf(phrases[senedPhrase], result.Name, result.Name, result.Name, result.Name))
		} else if senedPhrase > 4 && senedPhrase < 8 {
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
			panic(err)
		}

		if haveCoord {
			textMsgToCoordinator := fmt.Sprintf("Группа: %d\nБегун: %d %s", sendedNow, numb_id, result.Name)
			msg := tgbotapi.NewMessage(coordinatorId, textMsgToCoordinator)

			if _, err := bot.Send(msg); err != nil {
				panic(err)
			}

		}

		if sendedNow == screamersCount {
			sendedNow = 0
		}
		sendedNow++
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
	return findFlag
}

func delFromScreamersList(id int64) bool {
	findFlag := false
	idIndex := 0
	for i := 1; i <= screamersCount; i++ {
		if chatsList[i] == id {
			findFlag = true
		}
	}
	if findFlag == true {
		chatsList[idIndex] = chatsList[screamersCount]
		screamersCount--

	}

	return findFlag
}

func initPhrases() {
	phrases[0] = "Бегать быстро <b><u> %s </u></b> умеет, никого не пожалеет!\n\nЖылдам жүгіріңіз <b><u> %s </u></b> қалай біледі,ешкімге өкініш!"
	phrases[1] = "Беги всегда, беги везде — <b><u> %s </u></b>, беги к своей мечте!\n\nӘрқашан жүгіріңіз, барлық жерде жүгіріңіз <b><u> %s </u></b>, арманыңа жүгір!"
	phrases[2] = "За победу надо драться — <b><u> %s </u></b>, тебе придётся постараться!\n\nЖеңіс үшін күресу керек<b><u> %s </u></b>, тырысу керек!"
	phrases[3] = "Побеждать — твоя судьба! <b><u> %s </u></b>, победи всех как всегда!\n\nЖеңіс - сіздің тағдырыңыз!<b><u> %s </u></b>, әдеттегідей бәрін жең!"
	phrases[4] = "Эй, <b><u> %s </u></b>, ты не на прогулке!\nДавай <b><u> %s </u></b>, шевели булками!\n\nЭй <b><u> %s </u></b>, сіз серуендеуге шыққан жоқсыз!\nКеліңіздер <b><u> %s </u></b>, орамдарды жылжытыңыз!"
	phrases[5] = "Оббеги хоть всю планету, \nБыстрее <b><u> %s </u></b> в мире нету.\n <b><u> %s </u></b> - вперед!\n <b><u> %s </u></b> - давай!\n <b><u> %s </u></b> - беги, не засыпай!\n\nБүкіл планетаны айналып жүгіріңіз\nӘлемде жылдам <b><u> %s </u></b> жоқ.\n <b><u> %s </u></b> - алға!\n <b><u> %s </u></b> - келіңіз!\n <b><u> %s </u></b> - жүгір, ұйықтап қалма!"
	phrases[6] = "Кто забыл фразу «ой, всё, не могу…»?! <b><u> %s </u></b>\nКто шепчет себе «я добегу!»?! <b><u> %s </u></b>\nКто победит на раз и два! <b><u> %s </u></b>\nВыше нос <b><u> %s </u></b>! Ура! Ура!\n\n«Ой, болды, мен алмаймын...» дегенді ұмытқан кім?! <b><u> %s </u></b>\n«Мен жүгіремін!» деп сыбырлайтын кім бар?! <b><u> %s </u></b>\nБір және екеуі кім жеңеді! <b><u> %s </u></b>\nМұрын жоғары <b><u> %s </u></b>! Ура! Ура!"
	phrases[7] = "Кто всегда вперед стремится?! <b><u> %s </u></b>\nКто летит быстрее птицы?! <b><u> %s </u></b>\nСпортивный дух у у кого в крови?! <b><u> %s </u></b>\nДавай, <b><u> %s </u></b>! Всех, всех порви!\n\nКім әрқашан алға ұмтылады? <b><u> %s </u></b>\nҚұстан да жылдам ұшатын кім?! <b><u> %s </u></b>\nБіреудің қанында спорттық рух бар ма?! <b><u> %s </u></b>\nКеліңіздер, <b><u> %s </u></b>! Барлығын құртыңдар!"
	phrases[8] = "Раз, и два — бежать пора!\nТри, четыре — лучшим в мире!\nПять и шесть — в ком сила есть!\nСемь и восемь — темп не бросим!\nДевять, десять – победу разделим все вместе!\n\nБір, екі, жүгіретін уақыт келді!\nҮш, төрт - әлемдегі ең жақсы!\nБес пен алты – кімде билік бар!\nЖеті мен сегіз - қарқыннан бас тартпайық!\nТоғыз, он – жеңісті бірге бөлісейік!"
}

func main() {
	initPhrases()
	// BOT
	var err error
	bot, err = tgbotapi.NewBotAPI("5707889035:AAFQ1tEPMMUUfnp5cN4IWF68iwcIBOOUK6A")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	// Mongo
	clientOptions := options.Client().ApplyURI("mongodb+srv://admin-v2:dTEJ8hum0jfTH8bH@testcluster.prws5fu.mongodb.net/test")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("screamers").Collection("runners")

	// HTTP Server
	http.HandleFunc("/numb/", numb)
	http.HandleFunc("/headers", headers)

	go http.ListenAndServe(":8090", nil)

	for update := range updates {
		if update.Message != nil {

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

			switch update.Message.Command() {
			case "start":
				if addToScreamersList(update.Message.Chat.ID) {
					msg.Text = "Вы уже были добавлены в список получения!"
				} else {
					msg.Text = fmt.Sprintf("Вы добавлены в список получения! Ваш номер: %d", screamersCount)
				}
			case "stop":
				if delFromScreamersList(update.Message.Chat.ID) {
					msg.Text = "Вы удалены из списка на получение!"
				} else {
					msg.Text = "Вы уже были удалены из списка на получение!"
				}
			case "coord":
				coordinatorId = update.Message.Chat.ID
				haveCoord = true
				delFromScreamersList(update.Message.Chat.ID)
				msg.Text = "Вы назначены координатором!"
			case "stopcoord":
				haveCoord = false
				msg.Text = "Координатор отключен!"
			default:
				msg.Text = "Ошибка! Команда не найдена!"
			}

			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			bot.Send(msg)
		} else if update.CallbackQuery != nil {

			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)

			if _, err := bot.Request(callback); err != nil {
				panic(err)
			}

			log.Printf("Message ID %d", update.CallbackQuery.Message.MessageID)

			msg := tgbotapi.NewEditMessageReplyMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, tgbotapi.InlineKeyboardMarkup{
				InlineKeyboard: make([][]tgbotapi.InlineKeyboardButton, 0),
			})
			if _, err := bot.Request(msg); err != nil {
				panic(err)
			}
		}
	}
}
