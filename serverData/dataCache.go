package serverData

import (
	"gitlab.com/gameraccoon/telegram-accountant-bot/currencies"
	"math/big"
	"sync"
)

type ratesStruct struct {
	toUsd map[currencies.Currency]*big.Float
}

type erc20TokenData struct {
	name string
	symbol string
	decimals int64
}

type dataCache struct {
	rates ratesStruct
	ratesMutex sync.Mutex
	balances map[currencies.AddressData]*big.Int
	balancesMutex sync.Mutex
	erc20Tokens map[string]erc20TokenData
	erc20TokensMutex sync.Mutex
}

func (cache *dataCache) Init() {
	if cache.rates.toUsd == nil {
		cache.rates.toUsd = map[currencies.Currency]*big.Float{}
	}

	if cache.balances == nil {
		cache.balances = map[currencies.AddressData]*big.Int{}
	}
}

func (cache *dataCache) getBalance(address currencies.AddressData) *big.Int {
	cache.balancesMutex.Lock()
	defer cache.balancesMutex.Unlock()

	balance, ok := cache.balances[address]
	if ok {
		return balance
	} else {
		return nil
	}
}

func (cache *dataCache) getRateToUsd(currency currencies.Currency) *big.Float {
	cache.ratesMutex.Lock()
	defer cache.ratesMutex.Unlock()

	rateToUsd, ok := cache.rates.toUsd[currency]

	if ok {
		return rateToUsd
	} else {
		return nil
	}
}

func (cache *dataCache) getErc20TokenData(address currencies.AddressData) *currencies.Erc20TokenData {
	cache.erc20TokensMutex.Lock()
	defer cache.erc20TokensMutex.Unlock()

	tokenData, tokenFound := cache.erc20Tokens[address.ContractId]
	if tokenFound {
		return &currencies.Erc20TokenData {
			Name: tokenData.name,
			Symbol: tokenData.symbol,
			Decimals: tokenData.decimals,
		}
	} else {
		return nil
	}
}
