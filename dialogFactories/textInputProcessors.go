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
			"renamingWallet" : processRenamingWallet,
		},
	}
}

func processNewWatchOnlyWalletName(additionalId int64, data *processing.ProcessData) bool {
	data.Static.SetUserStateNewWalletName(data.UserId, data.Message)
	data.Static.SetUserStateTextProcessor(data.UserId, &processing.AwaitingTextProcessorData{
		ProcessorId: "newWatchOnlyWalletKey",
	})
	data.SendMessage(data.Trans("send_public_key"))
	return true
}

func processNewWatchOnlyWalletKey(additionalId int64, data *processing.ProcessData) bool {
	walletName := data.Static.GetUserStateNewWalletName(data.UserId)
	walletId := data.Static.Db.CreateWatchOnlyWallet(data.UserId, walletName, currencies.Bitcoin, data.Message)
	data.SendMessage(data.Trans("wallet_created"))
	data.SendDialog(data.Static.MakeDialogFn("ws", walletId, data.Trans, data.Static))
	return true
}

func processRenamingWallet(walletId int64, data *processing.ProcessData) bool {
	if walletId == 0 {
		return false
	}

	data.Static.Db.RenameWallet(walletId, data.Message)
	data.SendDialog(data.Static.MakeDialogFn("wa", walletId, data.Trans, data.Static))
	return true
}
