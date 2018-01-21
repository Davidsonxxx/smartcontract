package dialogFactories

import (
	"github.com/gameraccoon/telegram-bot-skeleton/dialogManager"
	"github.com/gameraccoon/telegram-bot-skeleton/processing"
	"gitlab.com/gameraccoon/telegram-accountant-bot/currencies"
	"gitlab.com/gameraccoon/telegram-accountant-bot/staticFunctions"
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
	data.Static.SetUserStateValue(data.UserId, "walletName", data.Message)
	data.Static.SetUserStateTextProcessor(data.UserId, &processing.AwaitingTextProcessorData{
		ProcessorId: "newWatchOnlyWalletKey",
	})
	data.SendMessage(data.Trans("send_address"))
	return true
}

func processNewWatchOnlyWalletKey(additionalId int64, data *processing.ProcessData) bool {
	walletName, ok := data.Static.GetUserStateValue(data.UserId, "walletName").(string)
	if !ok {
		return false
	}

	walletCurrency, ok := data.Static.GetUserStateValue(data.UserId, "walletCurrency").(currencies.Currency)
	if !ok {
		return false
	}

	walletId := staticFunctions.GetDb(data.Static).CreateWatchOnlyWallet(data.UserId, walletName, walletCurrency, data.Message)
	data.SendMessage(data.Trans("wallet_created"))
	data.SendDialog(data.Static.MakeDialogFn("wa", walletId, data.Trans, data.Static))
	return true
}

func processRenamingWallet(walletId int64, data *processing.ProcessData) bool {
	if walletId == 0 {
		return false
	}

	staticFunctions.GetDb(data.Static).RenameWallet(walletId, data.Message)
	data.SendDialog(data.Static.MakeDialogFn("wa", walletId, data.Trans, data.Static))
	return true
}
