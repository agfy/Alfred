package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
	"strconv"
)

func index(vs []string, t string) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}
	return -1
}

func include(vs []string, t string) bool {
	return index(vs, t) >= 0
}

type good struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Class    string `json:"class"`
	Shop     string `json:"shop"`
	Volume   int    `json:"volume"`
	Price    int    `json:"price"`
	FoodType string `json:"food_type"`
}

type order struct {
	Id         int    `json:"id"`
	TelegramId int    `json:"telegram_id"`
	GoodsId    int    `json:"goods_id"`
	Amount     int    `json:"amount"`
	CreateTime uint32 `json:"create_time"`
}

type qualification struct {
	Request  string `json:"request"`
	Shop     string `json:"shop"`
	FoodType string `json:"food_type"`
	Volume   string `json:"volume"`
	Class    string `json:"class"`
	Amount   string `json:"amount"`
	OrderId  string `json:"order_id"`
	GoodId   string `json:"good_id"`
}

var (
	requests = []string{
		"Сделать заказ",
		"Забрать все заказы",
		"Удалить заказ",
		"Добавить товар",
	}

	shops     = make([]string, 0)
	foodTypes = make([]string, 0)
	volumes   = make([]string, 0)
	classes   = make([]string, 0)

	amounts = []string{
		"1 шт",
		"2 шт",
		"3 шт",
		"4 шт",
		"5 шт",
	}
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

func (qual *qualification) clear() {
	qual.Amount = ""
	qual.Class = ""
	qual.FoodType = ""
	qual.OrderId = ""
	qual.Request = ""
	qual.Shop = ""
	qual.Volume = ""
}

func askForShop(chatId int64, messageId int, bot *tgbotapi.BotAPI) {
	if _, err := bot.Send(tgbotapi.NewEditMessageText(chatId, messageId, "Укажите магазин")); err != nil {
		log.Println(err)
	}
	if _, err := bot.Send(createReplyMarkup(shops, chatId, messageId)); err != nil {
		log.Println(err)
	}
}

func askForFoodType(chatId int64, messageId int, bot *tgbotapi.BotAPI) {
	if _, err := bot.Send(tgbotapi.NewEditMessageText(chatId, messageId, "Что вы желаете?")); err != nil {
		log.Println(err)
	}
	if _, err := bot.Send(createReplyMarkup(foodTypes, chatId, messageId)); err != nil {
		log.Println(err)
	}
}

func askForVolume(chatId int64, messageId int, bot *tgbotapi.BotAPI) {
	if _, err := bot.Send(tgbotapi.NewEditMessageText(chatId, messageId, "Какой объем бутылки вы хотите?")); err != nil {
		log.Println(err)
	}
	if _, err := bot.Send(createReplyMarkup(volumes, chatId, messageId)); err != nil {
		log.Println(err)
	}
}

func askForNumber(chatId int64, messageId int, bot *tgbotapi.BotAPI) {
	if _, err := bot.Send(tgbotapi.NewEditMessageText(chatId, messageId, "Сколько штук?")); err != nil {
		log.Println(err)
	}
	if _, err := bot.Send(createReplyMarkup(amounts, chatId, messageId)); err != nil {
		log.Println(err)
	}
}

func askForClass(chatId int64, messageId int, bot *tgbotapi.BotAPI, shop, foodType string) {
	if _, err := bot.Send(tgbotapi.NewEditMessageText(chatId, messageId, "Что вы предпочитаете?")); err != nil {
		log.Println(err)
	}
	footTypeClasses, _ := getClasses(shop, foodType)
	if _, err := bot.Send(createReplyMarkup(footTypeClasses, chatId, messageId)); err != nil {
		log.Println(err)
	}
}

func askForGoods(chatId int64, messageId int, bot *tgbotapi.BotAPI, shop, foodType, class, volume string) {
	if _, err := bot.Send(tgbotapi.NewEditMessageText(chatId, messageId, "Что вы предпочитаете?")); err != nil {
		log.Println(err)
	}
	gds, _ := getGoods(shop, foodType, class, volume)
	keyboard := tgbotapi.InlineKeyboardMarkup{}

	for _, gd := range gds {
		var row []tgbotapi.InlineKeyboardButton
		btn := tgbotapi.NewInlineKeyboardButtonData(gd.Name+" "+strconv.Itoa(gd.Price)+" руб", strconv.Itoa(gd.Id))
		row = append(row, btn)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	}

	if _, err := bot.Send(tgbotapi.NewEditMessageReplyMarkup(chatId, messageId, keyboard)); err != nil {
		log.Println(err)
	}
}

func askForBasicFunctions(chatId int64, messageId int, bot *tgbotapi.BotAPI, message string) {
	if _, err := bot.Send(tgbotapi.NewEditMessageText(chatId, messageId, message)); err != nil {
		log.Println(err)
	}
	if _, err := bot.Send(createReplyMarkup(requests, chatId, messageId)); err != nil {
		log.Println(err)
	}
}

func processCreatingOrder(query string, chatId int64, messageId, userId int, bot *tgbotapi.BotAPI, qualification *qualification) {
	if include(requests, query) {
		qualification.Request = query
		askForShop(chatId, messageId, bot)
	} else if include(shops, query) {
		qualification.Shop = query
		askForFoodType(chatId, messageId, bot)
	} else if include(foodTypes, query) {
		qualification.FoodType = query
		askForClass(chatId, messageId, bot, qualification.Shop, qualification.FoodType)
	} else if include(classes, query) {
		qualification.Class = query
		if qualification.FoodType == "напиток" {
			askForVolume(chatId, messageId, bot)
		} else {
			askForGoods(chatId, messageId, bot, qualification.Shop, qualification.FoodType, qualification.Class, qualification.Volume)
		}
	} else if include(volumes, query) {
		qualification.Volume = query
		askForGoods(chatId, messageId, bot, qualification.Shop, qualification.FoodType, qualification.Class, qualification.Volume)
	} else if _, err := strconv.Atoi(query); err == nil {
		qualification.GoodId = query
		askForNumber(chatId, messageId, bot)
	} else if include(amounts, query) {
		qualification.Amount = query
		amount := index(amounts, qualification.Amount) + 1
		goodId, _ := strconv.Atoi(qualification.GoodId)
		err = createOrder(userId, amount, goodId)

		qualification.clear()
		askForBasicFunctions(chatId, messageId, bot, "Ваш заказ сделан. Чем я могу помочь?")
	}
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

	users, err := getUsers()
	if err = fillValues(&shops, "shop"); err != nil {
		log.Println(err)
	}
	if err = fillValues(&foodTypes, "foodtype"); err != nil {
		log.Println(err)
	}
	if err = fillValues(&volumes, "volume"); err != nil {
		log.Println(err)
	}
	if err = fillValues(&classes, "class"); err != nil {
		log.Println(err)
	}

	userQualifications := make(map[int]*qualification)
	for id := range users {
		userQualifications[id] = &qualification{}
	}

	for update := range updates {
		if update.Message != nil {
			//log.Printf("[%s] [%s] %s", update.Message.From.UserName, update.Message.From.ID, update.Message.Text)
			if _, ok := users[update.Message.From.ID]; !ok {
				if _, err = bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Ваш telegramId: "+
					strconv.Itoa(update.Message.From.ID)+" Обратитесь к MaximGanker чтобы пользоваться ботом")); err != nil {
					log.Println(err)
				}
				continue
			}

			command := update.Message.Command()
			if command == "start" {
				userQualifications[update.Message.From.ID].clear()
				if _, err = bot.Send(createMessage(requests, update.Message.Chat.ID, "Чем я могу помочь?")); err != nil {
					log.Println(err)
				}
			}
		} else if update.Message == nil && update.CallbackQuery != nil {
			query := update.CallbackQuery.Data
			chatId := update.CallbackQuery.Message.Chat.ID
			messageId := update.CallbackQuery.Message.MessageID
			userId := update.CallbackQuery.From.ID

			if userQualifications[userId].Request == "" {
				if !include(requests, query) {
					log.Panic("userQualifications[userId].Request == \"\" and !include(requests, query)")
					continue
				}

				userQualifications[userId].Request = query
			}

			switch userQualifications[userId].Request {
			case requests[0]:
				processCreatingOrder(query, chatId, messageId, userId, bot, userQualifications[userId])
			}
		} else {
			log.Println("else")
			//log.Println(update.Message.Text)
			//msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			//msg.ReplyToMessageID = update.Message.MessageID

			//bot.Send(msg)
		}
	}
}
