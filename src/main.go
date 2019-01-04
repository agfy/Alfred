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
	Id   int
	Name string
	//Class    string
	//Shop     string
	Volume string
	Price  int
	//FoodType string
}

type order struct {
	Id              int
	OwnerTelegramId int
	GoodId          int
	Amount          int
	//CreateTime uint32
}

type qualification struct {
	Request  string
	Shop     string
	FoodType string
	Volume   string
	Class    string
	Amount   string
	OrderIds []int
	GoodId   string
}

const garageChatId = -1001245213385

var (
	requests = []string{
		"Сделать заказ",
		"Забрать все заказы",
		"Изменить мой заказ",
		//"Добавить товар",
	}

	shops     = make([]string, 0)
	foodTypes = make([]string, 0)
	volumes   = make([]string, 0)
	classes   = make([]string, 0)
	users     = make(map[int]string)

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
	qual.OrderIds = qual.OrderIds[:0]
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

func askOrdersToExclude(chatId int64, messageId, userId int, bot *tgbotapi.BotAPI, orderIds *[]int, messageText string) {
	if _, err := bot.Send(tgbotapi.NewEditMessageText(chatId, messageId, messageText)); err != nil {
		log.Println(err)
	}
	orders, _ := getOrders(userId)
	keyboard := tgbotapi.InlineKeyboardMarkup{}

	var row []tgbotapi.InlineKeyboardButton
	btn := tgbotapi.NewInlineKeyboardButtonData("Готово", "Готово")
	row = append(row, btn)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)

	for _, order := range orders {
		*orderIds = append(*orderIds, order.Id)
		gd, _ := getGood(order.GoodId)

		var row []tgbotapi.InlineKeyboardButton
		btn := tgbotapi.NewInlineKeyboardButtonData(gd.Name+" "+gd.Volume+" "+strconv.Itoa(order.Amount)+" шт", strconv.Itoa(order.Id))
		row = append(row, btn)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	}

	if _, err := bot.Send(tgbotapi.NewEditMessageReplyMarkup(chatId, messageId, keyboard)); err != nil {
		log.Println(err)
	}
}

func deleteOrder(orderId int, chatId int64, messageId int, bot *tgbotapi.BotAPI, orderIds *[]int, permanently bool) {
	for i, id := range *orderIds {
		if id == orderId {
			*orderIds = append((*orderIds)[:i], (*orderIds)[i+1:]...)
		}
	}
	if permanently {
		deleteOrderfunc(orderId)
	}

	if _, err := bot.Send(tgbotapi.NewEditMessageText(chatId, messageId, "Нажмите на заказ если хотите искоючить его из списка")); err != nil {
		log.Println(err)
	}
	keyboard := tgbotapi.InlineKeyboardMarkup{}

	var row []tgbotapi.InlineKeyboardButton
	btn := tgbotapi.NewInlineKeyboardButtonData("Готово", "Готово")
	row = append(row, btn)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)

	for _, orderId := range *orderIds {
		order, _ := getOrder(orderId)
		gd, _ := getGood(order.GoodId)

		var row []tgbotapi.InlineKeyboardButton
		btn := tgbotapi.NewInlineKeyboardButtonData(gd.Name+" "+gd.Volume+" "+strconv.Itoa(order.Amount)+" шт", strconv.Itoa(orderId))
		row = append(row, btn)
		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
	}

	if _, err := bot.Send(tgbotapi.NewEditMessageReplyMarkup(chatId, messageId, keyboard)); err != nil {
		log.Println(err)
	}
}

func orderReady(chatId int64, messageId, userId int, bot *tgbotapi.BotAPI, orderIds *[]int) {
	if len(*orderIds) > 0 {
		purchases := make(map[int]int)

		for _, orderId := range *orderIds {
			order, _ := getOrder(orderId)
			gd, _ := getGood(order.GoodId)

			purchases[order.OwnerTelegramId] += gd.Price * order.Amount
		}

		ordersMessage := users[userId] + " забрал все заказы: "
		for id, value := range purchases {
			ordersMessage += "@" + users[id] + " " + strconv.Itoa(value) + "руб, "
		}

		markOrdersBought(orderIds, userId)
		bot.Send(tgbotapi.NewMessage(garageChatId, ordersMessage))
	}

	askForBasicFunctions(chatId, messageId, bot, "Чем я могу еще помочь?")
}

func processCreatingOrder(query string, chatId int64, messageId, userId int, bot *tgbotapi.BotAPI, qualification *qualification) {
	if include(requests, query) {
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
		askForBasicFunctions(chatId, messageId, bot, "Ваш заказ сделан, он будет доступен в течение 12 часов. Чем я могу еще помочь?")
	}
}

func processTakingOrders(query string, chatId int64, messageId, userId int, bot *tgbotapi.BotAPI, qualification *qualification) {
	if include(requests, query) {
		askOrdersToExclude(chatId, messageId, 0, bot, &qualification.OrderIds,
			"Нажмите на заказ если хотите исключить его из списка")
	} else if orderId, err := strconv.Atoi(query); err == nil {
		deleteOrder(orderId, chatId, messageId, bot, &qualification.OrderIds, false)
	} else if query == "Готово" {
		orderReady(chatId, messageId, userId, bot, &qualification.OrderIds)
		qualification.clear()
	}
}

func processEditingOrders(query string, chatId int64, messageId, userId int, bot *tgbotapi.BotAPI, qualification *qualification) {
	if include(requests, query) {
		askOrdersToExclude(chatId, messageId, userId, bot, &qualification.OrderIds,
			"Нажмите на заказ если хотите удалить его")
	} else if orderId, err := strconv.Atoi(query); err == nil {
		deleteOrder(orderId, chatId, messageId, bot, &qualification.OrderIds, true)
	} else if query == "Готово" {
		askForBasicFunctions(chatId, messageId, bot, "Чем я могу еще помочь?")
		qualification.clear()
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
	users, err = getUsers()

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
		userQualifications[id] = &qualification{OrderIds: make([]int, 0)}
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
				if update.Message.Chat.ID == garageChatId {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Бот доступен только в приват режиме"))
					continue
				}

				println(update.Message.Chat.ID)
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
			case requests[1]:
				processTakingOrders(query, chatId, messageId, userId, bot, userQualifications[userId])
			case requests[2]:
				processEditingOrders(query, chatId, messageId, userId, bot, userQualifications[userId])
			}
		} else {
			log.Println("else")
		}
	}
}
