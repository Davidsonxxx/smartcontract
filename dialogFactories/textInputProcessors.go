package dialogFactories

import (
	"github.com/gameraccoon/telegram-bot-skeleton/dialogManager"
	"github.com/gameraccoon/telegram-bot-skeleton/processing"
	"gitlab.com/gameraccoon/telegram-accountant-bot/database"
	"gitlab.com/gameraccoon/telegram-accountant-bot/currencies"
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
	data.SendMessage(data.Trans("send_address"))
	return true
}

func processNewWatchOnlyWalletKey(additionalId int64, data *processing.ProcessData) bool {
	walletName := data.Static.GetUserStateNewWalletName(data.UserId)
	walletId := database.CreateWatchOnlyWallet(data.Static.Db, data.UserId, walletName, currencies.Bitcoin, data.Message)
	data.SendMessage(data.Trans("wallet_created"))
	data.SendDialog(data.Static.MakeDialogFn("wa", walletId, data.Trans, data.Static))
	return true
}

func processRenamingWallet(walletId int64, data *processing.ProcessData) bool {
	if walletId == 0 {
		return false
	}

	database.RenameWallet(data.Static.Db, walletId, data.Message)
	data.SendDialog(data.Static.MakeDialogFn("wa", walletId, data.Trans, data.Static))
	return true
}
