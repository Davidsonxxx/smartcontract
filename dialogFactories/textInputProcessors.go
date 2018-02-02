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
			"newWalletName" : processNewWalletName,
			"newWalletKey" : processNewWalletKey,
			"newWalletContractId" : processNewWalletContractId,
			"renamingWallet" : processRenamingWallet,
		},
	}
}

func processNewWalletName(additionalId int64, data *processing.ProcessData) bool {
	walletCurrency, ok := data.Static.GetUserStateValue(data.UserId, "walletCurrency").(currencies.Currency)
	if !ok {
		return false
	}

	data.Static.SetUserStateValue(data.UserId, "walletName", data.Message)

	if walletCurrency != currencies.Erc20Token {
		data.Static.SetUserStateTextProcessor(data.UserId, &processing.AwaitingTextProcessorData{
			ProcessorId: "newWalletKey",
		})
		data.SendMessage(data.Trans("send_address"))
	} else {
		// ERC20 Token
		data.Static.SetUserStateTextProcessor(data.UserId, &processing.AwaitingTextProcessorData{
			ProcessorId: "newWalletContractId",
		})
		data.SendMessage(data.Trans("send_contract_id"))
	}
	return true
}

func processNewWalletContractId(additionalId int64, data *processing.ProcessData) bool {
	data.Static.SetUserStateValue(data.UserId, "walletContractId", data.Message)
	data.Static.SetUserStateTextProcessor(data.UserId, &processing.AwaitingTextProcessorData{
		ProcessorId: "newWalletKey",
	})
	data.SendMessage(data.Trans("send_address"))
	return true
}

func processNewWalletKey(additionalId int64, data *processing.ProcessData) bool {
	walletName, ok := data.Static.GetUserStateValue(data.UserId, "walletName").(string)
	if !ok {
		return false
	}

	walletCurrency, ok := data.Static.GetUserStateValue(data.UserId, "walletCurrency").(currencies.Currency)
	if !ok {
		return false
	}

	walletContractId, ok := data.Static.GetUserStateValue(data.UserId, "walletContractId").(string)
	if !ok {
		walletContractId = ""
	}

	walletAddress := currencies.AddressData{
		Currency: walletCurrency,
		ContractId: walletContractId,
		Address: data.Message,
	}

	walletId := staticFunctions.GetDb(data.Static).CreateWatchOnlyWallet(data.UserId, walletName, walletAddress)
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
