package dialogFactories

import (
	"gitlab.com/gameraccoon/telegram-accountant-bot/currencies"
	"gitlab.com/gameraccoon/telegram-accountant-bot/processing"
	"gitlab.com/gameraccoon/telegram-accountant-bot/dialogManager"
)

func GetTextInputProcessorManager() dialogManager.TextInputProcessorManager {
	return dialogManager.TextInputProcessorManager {
		Processors : dialogManager.TextProcessorsMap {
			"newWatchOnlyWalletName" : processNewWatchOnlyWalletName,
			"newWatchOnlyWalletKey" : processNewWatchOnlyWalletKey,
		},
	}
}

func processNewWatchOnlyWalletName(additionalId string, data *processing.ProcessData) bool {
	data.Static.SetUserStateNewWalletName(data.UserId, data.Message)
	data.Static.SetUserStateTextProcessor(data.UserId, &processing.AwaitingTextProcessorData{
		ProcessorId: "newWatchOnlyWalletKey",
	})
	data.SendMessage(data.Trans("send_public_key"))
	return true
}

func processNewWatchOnlyWalletKey(additionalId string, data *processing.ProcessData) bool {
	walletName := data.Static.GetUserStateNewWalletName(data.UserId)
	data.Static.Db.CreateWatchOnlyWallet(data.UserId, walletName, currencies.Bitcoin, data.Message)
	data.SendMessage(data.Trans("wallet_created"))
	data.SendDialog(data.Static.MakeDialogFn("mn", data.UserId, data.Trans, data.Static))
	return true
}
