package main

import (
	"github.com/gameraccoon/telegram-bot-skeleton/dialogManager"
	"github.com/gameraccoon/telegram-bot-skeleton/processing"
	"gitlab.com/gameraccoon/telegram-accountant-bot/staticFunctions"
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
	data.SendDialog(data.Static.MakeDialogFn("cc", data.UserId, data.Trans, data.Static))
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

func UpdateProcessData(data *processing.ProcessData) {
	userId := staticFunctions.GetDb(data.Static).GetUserId(data.ChatId, data.UserSystemLang)
	data.UserId = userId
	data.Trans = staticFunctions.FindTransFunction(userId, data.Static)
}

func processCommand(data *processing.ProcessData, dialogManager *dialogManager.DialogManager, processors *ProcessorFuncMap) (succeeded bool) {
	UpdateProcessData(data)

	// drop any text processors for the case we will process a command
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
	UpdateProcessData(data)

	success := dialogManager.ProcessText(data)

	if !success {
		data.SendMessage(data.Trans("help_info"))
	}
}
