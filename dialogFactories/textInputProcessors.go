package dialogFactories

import (
	"github.com/gameraccoon/telegram-bot-skeleton/dialogManager"
	"github.com/gameraccoon/telegram-bot-skeleton/processing"
	"gitlab.com/gameraccoon/telegram-accountant-bot/currencies"
	"gitlab.com/gameraccoon/telegram-accountant-bot/staticFunctions"
	"gitlab.com/gameraccoon/telegram-accountant-bot/cryptoFunctions"
	"log"
	"regexp"
	"time"
)

func GetTextInputProcessorManager() dialogManager.TextInputProcessorManager {
	return dialogManager.TextInputProcessorManager {
		Processors : dialogManager.TextProcessorsMap {
			"newWalletName" : processNewWalletName,
			"newWalletKey" : processNewWalletKey,
			"newWalletContractAddress" : processNewWalletContractAddress,
			"renamingWallet" : processRenamingWallet,
			"setWalletPriceId" : processSetWalletPriceId,
			"newTimezone" : processSetTimezone,
		},
	}
}

func processNewWalletName(additionalId int64, data *processing.ProcessData) bool {
	walletCurrency, ok := data.Static.GetUserStateValue(data.UserId, "walletCurrency").(currencies.Currency)
	if !ok {
		return false
	}
	
	if len(data.Message) == 0 {
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
			ProcessorId: "newWalletContractAddress",
		})
		data.SendMessage(data.Trans("send_contract_id"))
	}
	return true
}

func processNewWalletContractAddress(additionalId int64, data *processing.ProcessData) bool {
	if len(data.Message) == 0 {
		data.SendMessage(data.Trans("wrong_contract_address"))
		return true
	}
	
	erc20TokenProcessor := cryptoFunctions.GetErc20TokenProcessor()
	if erc20TokenProcessor == nil {
		return false
	}

	if !(*erc20TokenProcessor).IsContractAddressValid(data.Message) {
		data.SendMessage(data.Trans("wrong_contract_address"))
		return true
	}

	data.Static.SetUserStateValue(data.UserId, "walletContractAddress", data.Message)
	data.Static.SetUserStateTextProcessor(data.UserId, &processing.AwaitingTextProcessorData{
		ProcessorId: "newWalletKey",
	})
	data.SendMessage(data.Trans("send_address"))
	return true
}

func processNewWalletKey(additionalId int64, data *processing.ProcessData) bool {
	if len(data.Message) == 0 {
		data.SendMessage(data.Trans("wrong_wallet_address"))
		return true
	}
	
	walletName, ok := data.Static.GetUserStateValue(data.UserId, "walletName").(string)
	if !ok {
		return false
	}

	walletCurrency, ok := data.Static.GetUserStateValue(data.UserId, "walletCurrency").(currencies.Currency)
	if !ok {
		return false
	}

	currencyProcessor := cryptoFunctions.GetProcessor(walletCurrency)
	if currencyProcessor == nil {
		return false
	}

	if !(*currencyProcessor).IsAddressValid(data.Message) {
		data.SendMessage(data.Trans("wrong_wallet_address"))
		return true
	}

	walletContractAddress, ok := data.Static.GetUserStateValue(data.UserId, "walletContractAddress").(string)
	if !ok {
		walletContractAddress = ""
	}

	walletAddress := currencies.AddressData{
		Currency: walletCurrency,
		ContractAddress: walletContractAddress,
		Address: data.Message,
		PriceId: currencies.GetCurrencyPriceId(walletCurrency),
	}

	walletId := staticFunctions.GetDb(data.Static).CreateWatchOnlyWallet(data.UserId, walletName, walletAddress)
	staticFunctions.GetDb(data.Static).EnableBalanceNotifies(walletId)
	data.SendMessage(data.Trans("wallet_created"))
	data.SendDialog(data.Static.MakeDialogFn("wa", walletId, data.Trans, data.Static))
	return true
}

func processRenamingWallet(walletId int64, data *processing.ProcessData) bool {
	if walletId == 0 {
		return false
	}
	
	if len(data.Message) == 0 {
		return false
	}

	staticFunctions.GetDb(data.Static).RenameWallet(walletId, data.Message)
	data.SendDialog(data.Static.MakeDialogFn("wa", walletId, data.Trans, data.Static))
	return true
}

func processSetWalletPriceId(walletId int64, data *processing.ProcessData) bool {
	if walletId == 0 {
		return false
	}

	re := regexp.MustCompile("https?:\\/\\/coinmarketcap\\.com\\/currencies\\/([\\w-_]+).*")
	if re == nil {
		log.Print("Wrong regexp")
		return false
	}

	matches := re.FindStringSubmatch(data.Message)

	if len(matches) <= 1 {
		staticFunctions.GetDb(data.Static).SetWalletPriceId(walletId, "")
		data.SendMessage(data.Trans("wrong_coinmarketcap_link"))
		data.SendDialog(data.Static.MakeDialogFn("wa", walletId, data.Trans, data.Static))
		return true
	}

	staticFunctions.GetDb(data.Static).SetWalletPriceId(walletId, matches[1])
	data.SendDialog(data.Static.MakeDialogFn("wa", walletId, data.Trans, data.Static))
	return true
}

func processSetTimezone(additionalId int64, data *processing.ProcessData) bool {
	_, err := time.LoadLocation(data.Message)

	if err == nil {
		staticFunctions.GetDb(data.Static).SetUserTimezone(data.UserId, data.Message)
		data.SendDialog(data.Static.MakeDialogFn("us", data.UserId, data.Trans, data.Static))
		return true
	} else {
		data.Static.SetUserStateTextProcessor(data.UserId, &processing.AwaitingTextProcessorData{
			ProcessorId: "newTimezone",
		})
		data.SendMessage(data.Trans("wrong_timezone") + "\n" + data.Trans("send_timezone"))
		return true
	}
}
