package main

import (
	"log"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"os"
)

func Index(vs []string, t string) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}
	return -1
}

func Include(vs []string, t string) bool {
	return Index(vs, t) >= 0
}

type good struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Class string `json:"class"`
	Shop string `json:"shop"`
	Volume int `json:"volume"`
	Price int `json:"price"`
	FoodType string `json:"food_type"`
}

type order struct {
	Id int `json:"id"`
	TelegramId int `json:"telegram_id"`
	GoodsId int `json:"goods_id"`
	Amount int `json:"amount"`
	CreateTime uint32 `json:"create_time"`
}

type qualification struct {
	Request string `json:"request"`
	Shop string `json:"shop"`
	FoodType string `json:"food_type"`
	Volume string `json:"volume"`
	Class string `json:"class"`
	Amount string `json:"amount"`
	OrderId string `json:"order_id"`
}

var (
	Requests = []string{
		"Сделать заказ",
		"Забрать все заказы",
		"Удалить заказ",
		"Добавить товар",
	}

	shops = []string{
		"Открывашка",
		"Литра",
		"Все равно",
	}

	foodTypes = []string{
		"Напиток",
		"Еда",
	}

	volumes = []string{
		"0.5 л",
		"1 л",
		"1.5 л",
	}

	amounts = []string{
		"1 шт",
		"2 шт",
		"3 шт",
		"4 шт",
		"5 шт",
	}
/*
	users = []int{
		123,
		234,
		345,
		456,
		345,
		123,
		234,
		345,
		456,
		345,
		234,
	}
*/
)

func createMessage(arr []string, chatId int64, text string) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(chatId, text)

	keyboard := tgbotapi.InlineKeyboardMarkup{}
	for _, class := range arr {
		var row []tgbotapi.InlineKeyboardButton
		btn := tgbotapi.NewInlineKeyboardButtonData(class, class)
		row = append(row, btn)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	}

	msg.ReplyMarkup = keyboard
	return msg
}

func createReplyMarkup(arr []string, chatId int64, messageId int) tgbotapi.EditMessageReplyMarkupConfig {
	keyboard := tgbotapi.InlineKeyboardMarkup{}

	for _, class := range arr {
		var row []tgbotapi.InlineKeyboardButton
		btn := tgbotapi.NewInlineKeyboardButtonData(class, class)
		row = append(row, btn)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	}

	return tgbotapi.NewEditMessageReplyMarkup(chatId, messageId, keyboard)
}

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	//bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	users, err := getTelegramIds()

	for _, id := range users {
		log.Println(id)
	}
	//orders := make(map[int]order)
	//goods := make(map[int]good)

	for update := range updates {
		if update.Message == nil && update.CallbackQuery != nil {
			query := update.CallbackQuery.Data
			chatId := update.CallbackQuery.Message.Chat.ID
			messageId := update.CallbackQuery.Message.MessageID
			log.Println("query")
			log.Println(query)
			if query == Requests[0] {
				bot.Send(tgbotapi.NewEditMessageText(chatId, messageId, "Какой магазин вы предпочитаете?"))
				bot.Send(createReplyMarkup(shops, chatId, messageId))
			} else if Include(shops, query){
				//goods[update.CallbackQuery.From.ID].Shop = query
				bot.Send(tgbotapi.NewEditMessageText(chatId, messageId, "Что вы желаете?"))
				bot.Send(createReplyMarkup(foodTypes, chatId, messageId))
			} else if query == foodTypes[0] {
				//goods[update.CallbackQuery.From.ID].Type = query
				bot.Send(tgbotapi.NewEditMessageText(chatId, messageId, "Какой объем бутылки вы хотите?"))
				bot.Send(createReplyMarkup(volumes, chatId, messageId))
			} else if Include(volumes, query) || query == foodTypes[1] {
				//goods[update.CallbackQuery.From.ID].Volume = query
				bot.Send(tgbotapi.NewEditMessageText(chatId, messageId, "Сколько штук?"))
				bot.Send(createReplyMarkup(amounts, chatId, messageId))
			} else if Include(amounts, query) {
				//goods[update.CallbackQuery.From.ID].Volume = query
				bot.Send(tgbotapi.NewEditMessageText(chatId, messageId, "Ваш заказ сделан. Чем я могу еще помочь?"))
				bot.Send(createReplyMarkup(Requests, chatId, messageId))
			}

		} else if(update.Message != nil){
			log.Printf("[%s] [%s] %s", update.Message.From.UserName, update.Message.From.ID, update.Message.Text)

			command := update.Message.Command()
			log.Println(command)
			if command == "start" {
				bot.Send(createMessage(Requests, update.Message.Chat.ID, "Чем я могу помочь?"))
			}
		} else{
			log.Println("else")
			//log.Println(update.Message.Text)
			//msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			//msg.ReplyToMessageID = update.Message.MessageID

			//bot.Send(msg)
		}


	}
}