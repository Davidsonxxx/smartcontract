package cryptoFunctions

import (
	"gitlab.com/gameraccoon/telegram-accountant-bot/currencies"
)

var processorsList map[currencies.Currency]currencies.CurrencyProcessor = map[currencies.Currency]currencies.CurrencyProcessor{
	currencies.Bitcoin : &BitcoinProcessor{},
	currencies.BitcoinCash : &BitcoinCashProcessor{},
	currencies.BitcoinGold : &BitcoinGoldProcessor{},
}

func GetProcessor(currency currencies.Currency) *currencies.CurrencyProcessor {
	processor, ok := processorsList[currency]

	if ok {
		return &processor
	} else {
		return nil
	}
}
