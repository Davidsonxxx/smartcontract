package serverData

import (
	"github.com/gameraccoon/telegram-bot-skeleton/processing"
	"gitlab.com/gameraccoon/telegram-accountant-bot/currencies"
	"math/big"
	"log"
	"sync"
)

type ratesStruct struct {
	toUsd map[currencies.Currency]*big.Float
}

type DataCache struct {
	rates ratesStruct
	ratesMutex sync.Mutex
	balances map[currencies.AddressData]*big.Int
	balancesMutex sync.Mutex
}

func GetServerDataCache(staticData *processing.StaticProccessStructs) *DataCache {
	if staticData == nil {
		log.Print("staticData is nil")
		return nil
	}

	dataCache, ok := staticData.GetCustomValue("serverDataCache").(*DataCache)
	if ok {
		return dataCache
	} else {
		log.Fatal("")
		return nil
	}
}

func (cache *DataCache) Init() {
	cache.rates.toUsd = map[currencies.Currency]*big.Float{}
	cache.balances = map[currencies.AddressData]*big.Int{}
}

func (cache *DataCache) GetBalance(address currencies.AddressData) *big.Int {
	cache.balancesMutex.Lock()
	defer cache.balancesMutex.Unlock()
	
	balance, ok := cache.balances[address]
	if ok {
		return balance
	} else {
		return nil
	}
}

func (cache *DataCache) GetRateToUsd(currency currencies.Currency) *big.Float {
	cache.ratesMutex.Lock()
	defer cache.ratesMutex.Unlock()
	
	rateToUsd, ok := cache.rates.toUsd[currency]

	if ok {
		return rateToUsd
	} else {
		return nil
	}
}
