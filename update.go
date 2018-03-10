package main

import (
	"github.com/gameraccoon/telegram-bot-skeleton/dialogManager"
	"github.com/gameraccoon/telegram-bot-skeleton/processing"
	"github.com/gameraccoon/telegram-bot-skeleton/telegramChat"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"gitlab.com/gameraccoon/telegram-accountant-bot/cryptoFunctions"
	"gitlab.com/gameraccoon/telegram-accountant-bot/serverData"
	"gitlab.com/gameraccoon/telegram-accountant-bot/staticFunctions"
	"log"
	"strings"
	"time"
)

type userChannel chan *processing.ProcessData

type userChannelData struct {
	channel userChannel
	// to be able to close old channels
	lastUpdateTime time.Time
}

type userChannelsData map[int64]*userChannelData

func startUpdating(chat *telegramChat.TelegramChat, dialogManager *dialogManager.DialogManager, staticData *processing.StaticProccessStructs, serverDataManager *serverData.ServerDataManager, updateIntervalSec int) {
	go updateTimer(staticData, serverDataManager, updateIntervalSec)
	updateBot(chat, staticData, dialogManager)
}

func updateTimer(staticData *processing.StaticProccessStructs, serverDataManager *serverData.ServerDataManager, updateIntervalSec int) {
	if updateIntervalSec <= 0 {
		log.Fatal("Wrong time interval. Add updateIntervalSec to config")
	}

	for {
		time.Sleep(time.Duration(updateIntervalSec) * time.Second)
		tickUpdateData := serverDataManager.TimerTick(staticFunctions.GetDb(staticData))
		tickAfterupdate(staticData, tickUpdateData)
	}
}

func tickAfterupdate(staticData *processing.StaticProccessStructs, tickUpdateData serverData.TickUpdateData) {
	db := staticFunctions.GetDb(staticData)

	serverData := serverData.GetServerData(staticData)

	if serverData == nil {
		log.Print("ServerData is nil")
		return
	}

	for _, balanceNotify := range tickUpdateData.BalanceNotifies {
		userChatId := db.GetUserChatId(balanceNotify.UserId)

		walletName := db.GetWalletName(balanceNotify.WalletId)

		walletAddress := db.GetWalletAddress(balanceNotify.WalletId)

		currencySymbol, currencyDecimals := staticFunctions.GetCurrencySymbolAndDecimals(serverData, walletAddress.Currency, walletAddress.ContractAddress)

		var oldBalanceStr string
		if balanceNotify.OldBalance != nil {
			oldBalanceStr = cryptoFunctions.FormatCurrencyAmount(balanceNotify.OldBalance, currencyDecimals)
		} else {
			oldBalanceStr = "0"
		}

		var newBalanceStr string
		if balanceNotify.NewBalance != nil {
			newBalanceStr = cryptoFunctions.FormatCurrencyAmount(balanceNotify.NewBalance, currencyDecimals)
		} else {
			newBalanceStr = "0"
		}

		translateMap := map[string]interface{}{
			"Name":   walletName,
			"Sign":   currencySymbol,
			"OldBal": oldBalanceStr,
			"NewBal": newBalanceStr,
		}

		translateFn := staticFunctions.FindTransFunction(balanceNotify.UserId, staticData)

		staticData.Chat.SendMessage(userChatId,
			translateFn("balance_notify_template", translateMap),
			0,
		)
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

	userChans := make(userChannelsData)

	for update := range updates {
		if update.Message != nil {
			processMessageUpdate(userChans, &update, staticData, dialogManager, &processors)
		}
		if update.CallbackQuery != nil {
			processCallbackUpdate(userChans, &update, staticData, dialogManager, &processors)
		}
	}
}

func processMessageUpdate(userChans userChannelsData, update *tgbotapi.Update, staticData *processing.StaticProccessStructs, dialogManager *dialogManager.DialogManager, processors *ProcessorFuncMap) {
	data := processing.ProcessData{
		Static:         staticData,
		ChatId:         update.Message.Chat.ID,
		UserSystemLang: strings.ToLower(update.Message.From.LanguageCode),
	}

	message := update.Message.Text

	if strings.HasPrefix(message, "/") {
		commandLen := strings.Index(message, " ")
		if commandLen != -1 {
			data.Command = message[1:commandLen]
			data.Message = message[commandLen+1:]
		} else {
			data.Command = message[1:]
		}
	} else {
		data.Message = message
	}

	processUpdate(userChans, &data, dialogManager, processors)
}

func processCallbackUpdate(userChans userChannelsData, update *tgbotapi.Update, staticData *processing.StaticProccessStructs, dialogManager *dialogManager.DialogManager, processors *ProcessorFuncMap) {
	data := processing.ProcessData{
		Static:            staticData,
		ChatId:            int64(update.CallbackQuery.From.ID),
		AnsweredMessageId: int64(update.CallbackQuery.Message.MessageID),
		UserSystemLang:    strings.ToLower(update.CallbackQuery.From.LanguageCode),
	}

	message := update.CallbackQuery.Data

	commandLen := strings.Index(message, " ")
	if commandLen != -1 {
		data.Command = message[1:commandLen]
		data.Message = message[commandLen+1:]
	} else {
		data.Command = message[1:]
	}

	processUpdate(userChans, &data, dialogManager, processors)
}

func processUpdate(userChans userChannelsData, data *processing.ProcessData, dialogManager *dialogManager.DialogManager, processors *ProcessorFuncMap) {
	userChanData, found := userChans[data.ChatId]

	if !found || userChanData == nil {
		userChanData = &userChannelData{
			channel: make(userChannel),
		}
		userChans[data.ChatId] = userChanData

		// start updates for a user
		go processUserUpdatesParallel(userChanData.channel, dialogManager, processors)
	}

	userChanData.lastUpdateTime = time.Now()

	// send parallel to not to wait
	go sendUpdate(userChanData.channel, data)
}

func sendUpdate(userChan userChannel, data *processing.ProcessData) {
	userChan <- data
}

func processUserUpdatesParallel(userChan userChannel, dialogManager *dialogManager.DialogManager, processors *ProcessorFuncMap) {
	for {
		updateData, chanIsOk := <-userChan

		if !chanIsOk {
			log.Print("Close channel")
			return
		}

		if len(updateData.Command) > 0 {
			processCommand(updateData, dialogManager, processors)
		} else {
			processPlainMessage(updateData, dialogManager)
		}
	}
}
