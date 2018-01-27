package cryptoFunctions

import (
	"gitlab.com/gameraccoon/telegram-accountant-bot/currencies"
)

var processorsList map[currencies.Currency]CurrencyProcessor = map[currencies.Currency]CurrencyProcessor{
	currencies.Bitcoin : &BitcoinProcessor{},
	currencies.BitcoinCash : &BitcoinCashProcessor{},
	currencies.BitcoinGold : &BitcoinGoldProcessor{},
	currencies.Ether : &EtherProcessor{},
	currencies.RippleXrp : &RippleXrpProcessor{},
}

var erc20Processor Erc20Processor = Erc20Processor{}

func GetProcessor(currency currencies.Currency) *CurrencyProcessor {
	processor, ok := processorsList[currency]

	if ok {
		return &processor
	} else {
		return nil
	}
}

func GetAllProcessors() map[currencies.Currency]CurrencyProcessor {
	processors := map[currencies.Currency]CurrencyProcessor {}

	// copy the map
	for currency, processor := range processorsList {
		processors[currency] = processor
	}

	return processors
}

func GetErc20Processor() *Erc20Processor {
	return &erc20Processor
}
