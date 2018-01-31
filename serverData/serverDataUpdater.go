package serverData

import (
	"gitlab.com/gameraccoon/telegram-accountant-bot/currencies"
	"gitlab.com/gameraccoon/telegram-accountant-bot/cryptoFunctions"
	"math/big"
	"log"
)

type serverDataUpdater struct {
	cache dataCache
}

func (dataUpdater *serverDataUpdater) updateBalanceOneWallet(walletAddress currencies.AddressData) *big.Int {
	processor := cryptoFunctions.GetProcessor(walletAddress.Currency)

	if processor == nil {
		log.Print("No processor found")
		return nil
	}

	balance := (*processor).GetBalance(walletAddress)

	if balance != nil {
		dataUpdater.cache.balancesMutex.Lock()
		dataUpdater.cache.balances[walletAddress] = balance
		dataUpdater.cache.balancesMutex.Unlock()
	}

	return balance
}

func (dataUpdater *serverDataUpdater) updateBalance(walletAddresses []currencies.AddressData) {
	if len(walletAddresses) == 0 {
		return
	}

	groupedWallets := make(map[currencies.Currency] []currencies.AddressData)

	for _, walletAddress := range walletAddresses {
		walletsSlice, ok := groupedWallets[walletAddress.Currency]
		if ok {
			groupedWallets[walletAddress.Currency] = append(walletsSlice, walletAddress)
		} else {
			groupedWallets[walletAddress.Currency] = []currencies.AddressData{ walletAddress }
		}
	}

	for currency, addresses := range groupedWallets {
		processor := cryptoFunctions.GetProcessor(currency)

		if processor == nil {
			log.Print("No processor found")
			continue
		}

		balances := (*processor).GetBalanceBunch(addresses)

		if len(addresses) != len(balances) {
			log.Printf("return count doesn't match input count: %d != %d", len(addresses), len(balances))
			continue
		}

		dataUpdater.cache.balancesMutex.Lock()

		for i, address := range addresses {
			balance := balances[i]
			if balance != nil {
				dataUpdater.cache.balances[address] = balance
			}
		}

		dataUpdater.cache.balancesMutex.Unlock()
	}
}

func (dataUpdater *serverDataUpdater) updateRates() {
	processors := cryptoFunctions.GetAllProcessors()

	toUsdRates := map[currencies.Currency]*big.Float{}

	for currency, processor := range processors {
		toUsdRate := processor.GetToUsdRate()

		if toUsdRate != nil {
			toUsdRates[currency] = toUsdRate
		}
	}

	dataUpdater.cache.balancesMutex.Lock()
	dataUpdater.cache.rates.toUsd = toUsdRates
	dataUpdater.cache.balancesMutex.Unlock()
}

func (dataUpdater *serverDataUpdater) updateErc20TokensData(contracts []string) {
	//processor := cryptoFunctions.GetErc20TokenProcessor()
	
	//for 
}
