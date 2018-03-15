package main

import (
	"github.com/gameraccoon/telegram-bot-skeleton/processing"
	"gitlab.com/gameraccoon/telegram-accountant-bot/cryptoFunctions"
	"gitlab.com/gameraccoon/telegram-accountant-bot/currencies"
	"gitlab.com/gameraccoon/telegram-accountant-bot/database"
	"gitlab.com/gameraccoon/telegram-accountant-bot/serverData"
	"gitlab.com/gameraccoon/telegram-accountant-bot/staticFunctions"
	"log"
)

func updateBalanceNotifies(staticData *processing.StaticProccessStructs, balanceNotifies []currencies.BalanceNotify) {
	db := staticFunctions.GetDb(staticData)

	serverData := serverData.GetServerData(staticData)

	if serverData == nil {
		log.Print("ServerData is nil")
		return
	}

	for _, balanceNotify := range balanceNotifies {
		if balanceNotify.IsInitialChange {
			continue
		}

		if len(balanceNotify.LastTransactions) > 0 {
			SendTransactionNotifications(staticData, db, serverData, &balanceNotify)
		} else {
			// send info just about balance change
			SendBalanceChangeNotification(staticData, db, serverData, &balanceNotify)
		}
	}
}

func SendBalanceChangeNotification(staticData *processing.StaticProccessStructs, db *database.AccountDb, serverData serverData.ServerDataInterface, balanceNotify *currencies.BalanceNotify) {
	userChatId := db.GetUserChatId(balanceNotify.UserId)
	walletName := db.GetWalletName(balanceNotify.WalletId)

	currencySymbol, currencyDecimals := staticFunctions.GetCurrencySymbolAndDecimals(serverData, balanceNotify.WalletAddress.Currency, balanceNotify.WalletAddress.ContractAddress)

	var oldBalanceStr string
	if balanceNotify.OldBalance != nil {
		oldBalanceStr = cryptoFunctions.FormatCurrencyAmount(balanceNotify.OldBalance, currencyDecimals)
	} else {
		oldBalanceStr = "0"
	}

	var newBalanceStr string
	if balanceNotify.NewBalance != nil {
		newBalanceStr = cryptoFunctions.FormatCurrencyAmount(balanceNotify.NewBalance, currencyDecimals)
	} else {
		newBalanceStr = "0"
	}

	translateMap := map[string]interface{}{
		"Name":   walletName,
		"Sign":   currencySymbol,
		"OldBal": oldBalanceStr,
		"NewBal": newBalanceStr,
	}

	translateFn := staticFunctions.FindTransFunction(balanceNotify.UserId, staticData)

	staticData.Chat.SendMessage(userChatId,
		translateFn("balance_notify_template", translateMap),
		0,
	)
}

func SendTransactionNotifications(staticData *processing.StaticProccessStructs, db *database.AccountDb, serverData serverData.ServerDataInterface, balanceNotify *currencies.BalanceNotify) {
	userChatId := db.GetUserChatId(balanceNotify.UserId)
	walletName := db.GetWalletName(balanceNotify.WalletId)

	currencySymbol, currencyDecimals := staticFunctions.GetCurrencySymbolAndDecimals(serverData, balanceNotify.WalletAddress.Currency, balanceNotify.WalletAddress.ContractAddress)

	// for each notification backwards
	for i := len(balanceNotify.LastTransactions) - 1; i >= 0; i-- {
		transaction := balanceNotify.LastTransactions[i]

		amountText := cryptoFunctions.FormatCurrencyAmount(transaction.Amount, currencyDecimals)

		var translateTemplate string
		if transaction.From == balanceNotify.WalletAddress.Address {
			translateTemplate = "sent_transaction_notify_template"
		} else if transaction.To == balanceNotify.WalletAddress.Address {
			translateTemplate = "recieved_transaction_notify_template"
		}

		translateMap := map[string]interface{}{
			"Name":   walletName,
			"From":   transaction.From,
			"To":     transaction.To,
			"Amount": amountText,
			"Sign":   currencySymbol,
			"Time":   staticFunctions.FormatTimestamp(transaction.Time),
		}

		translateFn := staticFunctions.FindTransFunction(balanceNotify.UserId, staticData)

		staticData.Chat.SendMessage(userChatId,
			translateFn(translateTemplate, translateMap),
			0,
		)
	}
}
