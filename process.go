package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"gitlab.com/gameraccoon/telegram-accountant-bot/dialogManager"
	"gitlab.com/gameraccoon/telegram-accountant-bot/processing"
	"strings"
)

type ProcessorFunc func(*processing.ProcessData)

type ProcessorFuncMap map[string]ProcessorFunc

func startCommand(data *processing.ProcessData) {
	data.SendMessage(data.Trans("disclaimer_message"))
	data.SendDialog(data.Static.MakeDialogFn("wl", data.UserId, data.Trans, data.Static))
}

func walletsCommand(data *processing.ProcessData) {
	data.SendDialog(data.Static.MakeDialogFn("wl", data.UserId, data.Trans, data.Static))
}

func createWalletCommand(data *processing.ProcessData) {
	data.SendDialog(data.Static.MakeDialogFn("cw", data.UserId, data.Trans, data.Static))
}

func settingsCommand(data *processing.ProcessData) {
	data.SendDialog(data.Static.MakeDialogFn("lc", data.UserId, data.Trans, data.Static))
}

func helpCommand(data *processing.ProcessData) {
	data.SendMessage(data.Trans("help_info"))
}

func makeUserCommandProcessors() ProcessorFuncMap {
	return map[string]ProcessorFunc{
		"start":      startCommand,
		"wallets":    walletsCommand,
		"new_wallet": createWalletCommand,
		"settings":   settingsCommand,
		"help":       helpCommand,
	}
}

func processCommandByProcessors(data *processing.ProcessData, processors *ProcessorFuncMap) bool {
	processor, ok := (*processors)[data.Command]
	if ok {
		processor(data)
	}

	return ok
}

func processCommand(data *processing.ProcessData, dialogManager *dialogManager.DialogManager, processors *ProcessorFuncMap) (succeeded bool) {
	// drop any text processors for the case wi will process a command
	data.Static.SetUserStateTextProcessor(data.UserId, nil)
	// process dialogs
	ids := strings.Split(data.Command, "_")
	if len(ids) >= 2 {
		dialogId := ids[0]
		variantId := ids[1]
		var additionalId string
		if len(ids) > 2 {
			additionalId = ids[2]
		}

		processed := dialogManager.ProcessVariant(dialogId, variantId, additionalId, data)
		if processed {
			return true
		}
	}

	// process static command
	processed := processCommandByProcessors(data, processors)
	if processed {
		return true
	}

	// if we here that means that no command was processed
	data.SendMessage(data.Trans("help_info"))
	data.SendDialog(data.Static.MakeDialogFn("wl", data.UserId, data.Trans, data.Static))
	return false
}

func processPlainMessage(data *processing.ProcessData, dialogManager *dialogManager.DialogManager) {
	success := dialogManager.ProcessText(data)

	if !success {
		data.SendMessage(data.Trans("help_info"))
	}
}

func processMessageUpdate(update *tgbotapi.Update, staticData *processing.StaticProccessStructs, dialogManager *dialogManager.DialogManager, processors *ProcessorFuncMap) {
	userId := staticData.Db.GetUserId(update.Message.Chat.ID, strings.ToLower(update.Message.From.LanguageCode))
	data := processing.ProcessData{
		Static: staticData,
		ChatId: update.Message.Chat.ID,
		UserId: userId,
		Trans:  staticData.FindTransFunction(userId),
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

		processCommand(&data, dialogManager, processors)
	} else {
		data.Message = message
		processPlainMessage(&data, dialogManager)
	}
}

func processCallbackUpdate(update *tgbotapi.Update, staticData *processing.StaticProccessStructs, dialogManager *dialogManager.DialogManager, processors *ProcessorFuncMap) {
	userId := staticData.Db.GetUserId(int64(update.CallbackQuery.From.ID), strings.ToLower(update.CallbackQuery.From.LanguageCode))
	data := processing.ProcessData{
		Static:            staticData,
		ChatId:            int64(update.CallbackQuery.From.ID),
		UserId:            userId,
		Trans:             staticData.FindTransFunction(userId),
		AnsweredMessageId: int64(update.CallbackQuery.Message.MessageID),
	}

	message := update.CallbackQuery.Data
	commandLen := strings.Index(message, " ")
	if commandLen != -1 {
		data.Command = message[1:commandLen]
		data.Message = message[commandLen+1:]
	} else {
		data.Command = message[1:]
	}

	processCommand(&data, dialogManager, processors)
}
