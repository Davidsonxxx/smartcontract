package serverData

import (
	"gitlab.com/gameraccoon/telegram-accountant-bot/currencies"
	"gitlab.com/gameraccoon/telegram-accountant-bot/cryptoFunctions"
	"gitlab.com/gameraccoon/telegram-accountant-bot/database"
	"math/big"
	"log"
)

type serverDataUpdater struct {
	cache dataCache
}

type balanceChangesData map[int64]*big.Int

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

func (dataUpdater *serverDataUpdater) updateBalance(walletAddresses []database.WalletAddressDbWrapper) (balanceChanges balanceChangesData) {
	if len(walletAddresses) == 0 {
		return
	}

	balanceChanges = make(balanceChangesData)

	// group wallets to process in groups
	groupedWallets := make(map[currencies.Currency] []database.WalletAddressDbWrapper)

	for _, walletAddressWrapper := range walletAddresses {
		walletsSlice, ok := groupedWallets[walletAddressWrapper.Data.Currency]
		if ok {
			groupedWallets[walletAddressWrapper.Data.Currency] = append(walletsSlice, walletAddressWrapper)
		} else {
			groupedWallets[walletAddressWrapper.Data.Currency] = []database.WalletAddressDbWrapper{ walletAddressWrapper }
		}
	}

	for currency, addressWrappers := range groupedWallets {
		processor := cryptoFunctions.GetProcessor(currency)

		if processor == nil {
			log.Print("No processor found")
			continue
		}

		addresses := []currencies.AddressData{}
		for _, addressWrapper := range addressWrappers {
			addresses = append(addresses, addressWrapper.Data)
		}

		// request and get balances
		balances := (*processor).GetBalanceBunch(addresses)

		if len(addressWrappers) != len(balances) {
			log.Printf("return count doesn't match input count: %d != %d", len(addressWrappers), len(balances))
			continue
		}

		dataUpdater.cache.balancesMutex.Lock()

		for i, addressWrapper := range addressWrappers {
			balance := balances[i]
			if balance != nil {
				oldBalance := dataUpdater.cache.balances[addressWrapper.Data]
				if oldBalance == nil || balance.Cmp(oldBalance) != 0 {
					// create a record to trigger notifies
					balanceChanges[addressWrapper.WalletId] =  new(big.Int).Set(balance)

					// change cached balance value
					dataUpdater.cache.balances[addressWrapper.Data] = balance
				}
			}
		}

		dataUpdater.cache.balancesMutex.Unlock()
	}
	return
}

func (dataUpdater *serverDataUpdater) updateRates(priceIds []string) {

	toUsdRates := map[string]*big.Float{}

	for _, priceId := range priceIds {
		toUsdRate := cryptoFunctions.GetCurrencyToUsdRate(priceId)

		if toUsdRate != nil {
			toUsdRates[priceId] = toUsdRate
		}
	}

	dataUpdater.cache.balancesMutex.Lock()

	for priceId, rate := range toUsdRates {
		if rate != nil {
			dataUpdater.cache.rates.toUsd[priceId] = rate
		}
	}

	dataUpdater.cache.balancesMutex.Unlock()
}

func (dataUpdater *serverDataUpdater) updateErc20TokensData(contractAddresses []string) {
	if len(contractAddresses) <= 0 {
		return
	}

	processor := cryptoFunctions.GetErc20TokenProcessor()

	if processor == nil {
		log.Print("no ERC20 Token processor")
		return
	}

	tokenDatas := make(map[string]*currencies.Erc20TokenData)

	for _, contractAddress := range contractAddresses {
		tokenDatas[contractAddress] = processor.GetTokenData(contractAddress)
	}

	dataUpdater.cache.balancesMutex.Lock()
	for contractAddress, contractData := range tokenDatas {
		if contractData != nil {
			dataUpdater.cache.erc20Tokens[contractAddress] = *contractData
		}
	}

	dataUpdater.cache.balancesMutex.Unlock()
}

func (dataUpdater *serverDataUpdater) updateOneErc20TokensData(contractAddress string) *currencies.Erc20TokenData {
	processor := cryptoFunctions.GetErc20TokenProcessor()

	if processor == nil {
		log.Print("no ERC20 Token processor")
		return nil
	}

	tokenData := processor.GetTokenData(contractAddress)

	if tokenData != nil {
		dataUpdater.cache.balancesMutex.Lock()
		dataUpdater.cache.erc20Tokens[contractAddress] = *tokenData
		dataUpdater.cache.balancesMutex.Unlock()
	}

	return tokenData
}
