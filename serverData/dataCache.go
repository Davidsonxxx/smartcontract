package serverData

import (
	"gitlab.com/gameraccoon/telegram-accountant-bot/currencies"
	"gitlab.com/gameraccoon/telegram-accountant-bot/cryptoFunctions"
	"math/big"
	"log"
)

type ratesStruct struct {
	toUsd map[currencies.Currency]*big.Float
}

type dataCache struct {
	rates ratesStruct
	balances map[currencies.AddressData]*big.Int
}

func (cache *dataCache) updateBalance(walletAddresses []currencies.AddressData) {
	if len(walletAddresses) == 0 {
		return
	}

	groupedWallets := make(map[currencies.Currency] []string)

	for _, walletAddress := range walletAddresses {
		walletsSlice, ok := groupedWallets[walletAddress.Currency]
		if ok {
			groupedWallets[walletAddress.Currency] = append(walletsSlice, walletAddress.Address)
		} else {
			groupedWallets[walletAddress.Currency] = []string{ walletAddress.Address }
		}
	}

	for currency, addresses := range groupedWallets {
		processor := cryptoFunctions.GetProcessor(currency)

		if processor == nil {
			log.Print("No processor found")
			continue
		}

		balances := []*big.Int{}

		if processor != nil {
			balances = (*processor).GetBalanceBunch(addresses)
		}

		if len(addresses) != len(balances) {
			log.Printf("return count doesn't match input count: %d != %d", len(addresses), len(balances))
			continue
		}

		for i, address := range addresses {
			balance := balances[i]
			if balance != nil {
				addressData := currencies.AddressData{
					Currency: currency,
					Address: address,
				}
				cache.balances[addressData] = balance
			}
		}
	}
}

func (cache *dataCache) updateRates() {
	processors := cryptoFunctions.GetAllProcessors()

	for currency, processor := range processors {
		toUsdRate := processor.GetToUsdRate()

		if toUsdRate != nil {
			cache.rates.toUsd[currency] = toUsdRate
		}
	}
}
