package cryptoFunctions

import (
	"gitlab.com/gameraccoon/telegram-accountant-bot/currencies"
)

var processorsList map[currencies.Currency]currencies.CurrencyProcessor = map[currencies.Currency]currencies.CurrencyProcessor{
	currencies.Bitcoin : &BitcoinProcessor{},
}

func GetProcessor(currency currencies.Currency) *currencies.CurrencyProcessor {
	processor, ok := processorsList[currency]

	if ok {
		return &processor
	} else {
		return nil
	}
}
