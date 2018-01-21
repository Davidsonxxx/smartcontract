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
}

func (cache *dataCache) Init() {
	if cache.rates.toUsd == nil {
		cache.rates.toUsd = map[currencies.Currency]*big.Float{}
	}

	if cache.balances == nil {
		cache.balances = map[currencies.AddressData]*big.Int{}
	}
}

func (cache *dataCache) GetBalance(address currencies.AddressData) *big.Int {
	cache.balancesMutex.Lock()
	defer cache.balancesMutex.Unlock()

	balance, ok := cache.balances[address]
	if ok {
		return balance
	} else {
		return nil
	}
}

func (cache *dataCache) GetRateToUsd(currency currencies.Currency) *big.Float {
	cache.ratesMutex.Lock()
	defer cache.ratesMutex.Unlock()

	rateToUsd, ok := cache.rates.toUsd[currency]

	if ok {
		return rateToUsd
	} else {
		return nil
	}
}