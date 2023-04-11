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
var chatsList [50]int64
var screamersCount = 0
var sendedNow = 1
var senedPhrase = 0

var phrases [10]string

var collection *mongo.Collection
var ctx = context.TODO()

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
			text_msg = fmt.Sprintf("%d:\n %s", numb_id, fmt.Sprintf(phrases[senedPhrase], result.Name, result.Name))
		} else if senedPhrase == 4 {
			text_msg = fmt.Sprintf("%d:\n %s", numb_id, fmt.Sprintf(phrases[senedPhrase], result.Name, result.Name, result.Name, result.Name))
		} else if senedPhrase > 4 && senedPhrase < 8 {
			text_msg = fmt.Sprintf("%d:\n %s", numb_id, fmt.Sprintf(phrases[senedPhrase], result.Name, result.Name, result.Name, result.Name, result.Name, result.Name, result.Name, result.Name))
		} else {
			text_msg = fmt.Sprintf("%d:\n %s", numb_id, phrases[senedPhrase])
			senedPhrase = -1
		}

		senedPhrase++
		msg := tgbotapi.NewMessage(chatsList[sendedNow], text_msg)

		if sendedNow == screamersCount {
			sendedNow = 0
		}
		sendedNow++

		bot.Send(msg)
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
	phrases[0] = "Бегать быстро %s умеет, никого не пожалеет!\n\nЖылдам жүгіріңіз %s қалай біледі,ешкімге өкініш!"
	phrases[1] = "Беги всегда, беги везде — %s, беги к своей мечте!\n\nӘрқашан жүгіріңіз, барлық жерде жүгіріңіз %s, арманыңа жүгір!"
	phrases[2] = "За победу надо драться — %s, тебе придётся постараться!\n\nЖеңіс үшін күресу керек%s, тырысу керек!"
	phrases[3] = "Побеждать — твоя судьба! %s, победи всех как всегда!\n\nЖеңіс - сіздің тағдырыңыз!%s, әдеттегідей бәрін жең!"
	phrases[4] = "Эй, %s, ты не на прогулке!\nДавай %s, шевели булками!\n\nЭй %s, сіз серуендеуге шыққан жоқсыз!\nКеліңіздер %s, орамдарды жылжытыңыз!"
	phrases[5] = "Оббеги хоть всю планету, \nБыстрее %s в мире нету.\n %s - вперед!\n %s - давай!\n %s - беги, не засыпай!\n\nБүкіл планетаны айналып жүгіріңіз\nӘлемде жылдам %s жоқ.\n %s - алға!\n %s - келіңіз!\n %s - жүгір, ұйықтап қалма!"
	phrases[6] = "Кто забыл фразу «ой, всё, не могу…»?! %s\nКто шепчет себе «я добегу!»?! %s\nКто победит на раз и два! %s\nВыше нос %s! Ура! Ура!\n\n«Ой, болды, мен алмаймын...» дегенді ұмытқан кім?! %s\n«Мен жүгіремін!» деп сыбырлайтын кім бар?! %s\nБір және екеуі кім жеңеді! %s\nМұрын жоғары %s! Ура! Ура!"
	phrases[7] = "Кто всегда вперед стремится?! %s\nКто летит быстрее птицы?! %s\nСпортивный дух у у кого в крови?! %s\nДавай, %s! Всех, всех порви!\n\nКім әрқашан алға ұмтылады? %s\nҚұстан да жылдам ұшатын кім?! %s\nБіреудің қанында спорттық рух бар ма?! %s\nКеліңіздер, %s! Барлығын құртыңдар!"
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
		if update.Message != nil { // If we got a message

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

			switch update.Message.Command() {
			case "start":
				if addToScreamersList(update.Message.Chat.ID) {
					msg.Text = "Вы уже были добавлены в список получения!"
				} else {
					msg.Text = "Вы добавлены в список получения!"
				}
			case "stop":
				if delFromScreamersList(update.Message.Chat.ID) {
					msg.Text = "Вы удалены из списка на получение!"
				} else {
					msg.Text = "Вы уже были удалены из списка на получение!"
				}

			default:
				msg.Text = "Ошибка! Команда не найдена!"
			}

			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			bot.Send(msg)
		}
	}
}
