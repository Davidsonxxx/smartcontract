package main

import (
	"github.com/gameraccoon/telegram-bot-skeleton/dialogManager"
	"github.com/gameraccoon/telegram-bot-skeleton/processing"
	"github.com/gameraccoon/telegram-bot-skeleton/telegramChat"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"gitlab.com/gameraccoon/telegram-accountant-bot/serverData"
	"gitlab.com/gameraccoon/telegram-accountant-bot/staticFunctions"
	"log"
	"strings"
	"time"
)

func updateTimer(staticData *processing.StaticProccessStructs, serverDataManager *serverData.ServerDataManager, updateIntervalSec int) {
	if updateIntervalSec <= 0 {
		log.Fatal("Wrong time interval. Add updateIntervalSec to config")
	}

	for {
		time.Sleep(time.Duration(updateIntervalSec) * time.Second)
		serverDataManager.TimerTick(staticFunctions.GetDb(staticData))
	}
}

func updateBot(chat *telegramChat.TelegramChat, staticData *processing.StaticProccessStructs, dialogManager *dialogManager.DialogManager) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := chat.GetBot().GetUpdatesChan(u)

	if err != nil {
		log.Fatal(err.Error())
	}

	processors := makeUserCommandProcessors()

	for update := range updates {
		if update.Message != nil {
			processMessageUpdate(&update, staticData, dialogManager, &processors)
		}
		if update.CallbackQuery != nil {
			processCallbackUpdate(&update, staticData, dialogManager, &processors)
		}
	}
}

func processMessageUpdate(update *tgbotapi.Update, staticData *processing.StaticProccessStructs, dialogManager *dialogManager.DialogManager, processors *ProcessorFuncMap) {
	data := processing.ProcessData{
		Static: staticData,
		ChatId: update.Message.Chat.ID,
	}

	userLangCode := strings.ToLower(update.Message.From.LanguageCode)

	message := update.Message.Text

	if strings.HasPrefix(message, "/") {
		commandLen := strings.Index(message, " ")
		if commandLen != -1 {
			data.Command = message[1:commandLen]
			data.Message = message[commandLen+1:]
		} else {
			data.Command = message[1:]
		}

		processCommand(&data, dialogManager, processors, userLangCode)
	} else {
		data.Message = message

		processPlainMessage(&data, dialogManager, userLangCode)
	}
}

func processCallbackUpdate(update *tgbotapi.Update, staticData *processing.StaticProccessStructs, dialogManager *dialogManager.DialogManager, processors *ProcessorFuncMap) {
	data := processing.ProcessData{
		Static:            staticData,
		ChatId:            int64(update.CallbackQuery.From.ID),
		AnsweredMessageId: int64(update.CallbackQuery.Message.MessageID),
	}

	userLangCode := strings.ToLower(update.CallbackQuery.From.LanguageCode)

	message := update.CallbackQuery.Data

	commandLen := strings.Index(message, " ")
	if commandLen != -1 {
		data.Command = message[1:commandLen]
		data.Message = message[commandLen+1:]
	} else {
		data.Command = message[1:]
	}

	processCommand(&data, dialogManager, processors, userLangCode)
}
