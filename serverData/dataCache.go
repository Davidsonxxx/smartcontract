package serverData

import (
	"gitlab.com/gameraccoon/telegram-accountant-bot/currencies"
	"math/big"
	"sync"
)

type ratesStruct struct {
	toUsd map[currencies.Currency]*big.Float
}

type dataCache struct {
	rates ratesStruct
	ratesMutex sync.Mutex
	balances map[currencies.AddressData]*big.Int
	balancesMutex sync.Mutex
	erc20Tokens map[string]currencies.Erc20TokenData
	erc20TokensMutex sync.Mutex
}

func (cache *dataCache) Init() {
	if cache.rates.toUsd == nil {
		cache.rates.toUsd = make(map[currencies.Currency]*big.Float)
	}

	if cache.balances == nil {
		cache.balances = make(map[currencies.AddressData]*big.Int)
	}

	if cache.erc20Tokens == nil {
		cache.erc20Tokens = make(map[string]currencies.Erc20TokenData)
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

func (cache *dataCache) getErc20TokenData(contractId string) *currencies.Erc20TokenData {
	cache.erc20TokensMutex.Lock()
	defer cache.erc20TokensMutex.Unlock()

	tokenData, tokenFound := cache.erc20Tokens[contractId]
	if tokenFound {
		return &currencies.Erc20TokenData {
			Name: tokenData.Name,
			Symbol: tokenData.Symbol,
			Decimals: tokenData.Decimals,
		}
	} else {
		return nil
	}
}
