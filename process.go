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
	data.Static.Chat.SendMessage(data.ChatId, data.Trans("disclaimer_message"))
	data.Static.Chat.SendDialog(data.ChatId, data.Static.MakeDialogFn("mn", data.UserId, data.Trans, data.Static))
}

func settingsCommand(data *processing.ProcessData) {
	data.Static.Chat.SendDialog(data.ChatId, data.Static.MakeDialogFn("ls", data.UserId, data.Trans, data.Static))
}

func createWalletCommand(data *processing.ProcessData) {

}

func makeUserCommandProcessors() ProcessorFuncMap {
	return map[string]ProcessorFunc{
		"start":      startCommand,
		"settings":   settingsCommand,
		"new_wallet": createWalletCommand,
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
	data.Static.Chat.SendMessage(data.ChatId, data.Trans("warn_unknown_command"))
	data.Static.Chat.SendDialog(data.ChatId, data.Static.MakeDialogFn("mn", data.UserId, data.Trans, data.Static))
	return false
}

func processPlainMessage(data *processing.ProcessData, dialogManager *dialogManager.DialogManager) {
	success := dialogManager.ProcessText(data)

	if !success {
		data.Static.Chat.SendMessage(data.ChatId, data.Trans("warn_unknown_command"))
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
		Static: staticData,
		ChatId: int64(update.CallbackQuery.From.ID),
		UserId: userId,
		Trans:  staticData.FindTransFunction(userId),
	}

	message := update.CallbackQuery.Data
	commandLen := strings.Index(message, " ")
	if commandLen != -1 {
		data.Command = message[1:commandLen]
		data.Message = message[commandLen+1:]
	} else {
		data.Command = message[1:]
	}

	isSucceeded := processCommand(&data, dialogManager, processors)
	if isSucceeded {
		// deleteConfig := tgbotapi.DeleteMessageConfig{
		// 	ChatID: data.ChatId,
		// 	MessageID:,
		// }

		// _, err := tgbotapi.DeleteMessage(deleteConfig)
		// if err != nil {

		// }
	}
}
