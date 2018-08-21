package main

import (
	"github.com/gameraccoon/telegram-bot-skeleton/processing"
	"gitlab.com/gameraccoon/telegram-accountant-bot/cryptoFunctions"
	"gitlab.com/gameraccoon/telegram-accountant-bot/currencies"
	"gitlab.com/gameraccoon/telegram-accountant-bot/database"
	"gitlab.com/gameraccoon/telegram-accountant-bot/serverData"
	"gitlab.com/gameraccoon/telegram-accountant-bot/staticFunctions"
	"log"
	"math/big"
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

		SendBalanceChangeNotification(staticData, db, serverData, &balanceNotify)
	}
}

func SendBalanceChangeNotification(staticData *processing.StaticProccessStructs, db *database.AccountDb, serverData serverData.ServerDataInterface, balanceNotify *currencies.BalanceNotify) {
	if balanceNotify.OldBalance != nil && balanceNotify.NewBalance != nil {
		userChatId := db.GetUserChatId(balanceNotify.UserId)
		walletName := db.GetWalletName(balanceNotify.WalletId)

		currencySymbol, currencyDecimals := staticFunctions.GetCurrencySymbolAndDecimals(serverData, balanceNotify.WalletAddress.Currency, balanceNotify.WalletAddress.ContractAddress)

		var balanceDiff = new(big.Int)
		var balanceNotifyTemplate string

		if balanceNotify.NewBalance.Cmp(balanceNotify.OldBalance) > 0 {
			balanceDiff.Sub(balanceNotify.NewBalance, balanceNotify.OldBalance)
			balanceNotifyTemplate = "balance_notify_inc_template"
		} else {
			balanceDiff.Sub(balanceNotify.OldBalance, balanceNotify.NewBalance)
			balanceNotifyTemplate = "balance_notify_dec_template"
		}

		var balanceDiffStr = cryptoFunctions.FormatCurrencyAmount(balanceDiff, currencyDecimals)
		var newBalanceStr = cryptoFunctions.FormatCurrencyAmount(balanceNotify.NewBalance, currencyDecimals)

		translateMap := map[string]interface{}{
			"Name":   walletName,
			"Sign":   currencySymbol,
			"Diff":   balanceDiffStr,
			"NewBal": newBalanceStr,
		}

		translateFn := staticFunctions.FindTransFunction(balanceNotify.UserId, staticData)

		staticData.Chat.SendMessage(userChatId,
			translateFn(balanceNotifyTemplate, translateMap),
			0,
		)
	}
}
