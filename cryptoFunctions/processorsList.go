package cryptoFunctions

import (
	"gitlab.com/gameraccoon/telegram-accountant-bot/currencies"
)

var processorsList map[currencies.Currency]CurrencyProcessor = map[currencies.Currency]CurrencyProcessor{
	currencies.Bitcoin : &BitcoinProcessor{},
	currencies.BitcoinCash : &BitcoinCashProcessor{},
	currencies.BitcoinGold : &BitcoinGoldProcessor{},
	currencies.Ether : &EtherProcessor{},
}

func GetProcessor(currency currencies.Currency) *CurrencyProcessor {
	processor, ok := processorsList[currency]

	if ok {
		return &processor
	} else {
		return nil
	}
}
