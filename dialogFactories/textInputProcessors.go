package dialogFactories

import (
	"gitlab.com/gameraccoon/telegram-accountant-bot/processing"
	"gitlab.com/gameraccoon/telegram-accountant-bot/dialogManager"
)

func GetTextInputProcessorManager() dialogManager.TextInputProcessorManager {
	return dialogManager.TextInputProcessorManager {
		Processors : dialogManager.TextProcessorsMap {
			"newWalletName" : processAddWallet,
		},
	}
}

func processAddWallet(additionalId string, data *processing.ProcessData) bool {
	// newWalletId := data.Static.Db.CreateWatchOnlyWallet(data.UserId, data.Message)
	// data.Static.SetUserStateTextProcessor(data.UserId, &processing.AwaitingTextProcessorData{
	// 	ProcessorId: "whatsnext", // ToDo:
	// 	AdditionalId: strconv.FormatInt(newWalletId, 10),
	// })
	// data.Static.Chat.SendMessage(data.ChatId, data.Static.Trans("say_wait_items"))
	return true
}
