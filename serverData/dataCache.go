package serverData

import (
	"github.com/gameraccoon/telegram-bot-skeleton/processing"
	"gitlab.com/gameraccoon/telegram-accountant-bot/currencies"
	"math/big"
	"log"
)

type ratesStruct struct {
	toUsd map[currencies.Currency]*big.Float
}

type DataCache struct {
	rates ratesStruct
	balances map[currencies.AddressData]*big.Int
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
	balance, ok := cache.balances[address]
	if ok {
		return balance
	} else {
		return nil
	}
}

func (cache *DataCache) GetRateToUsd(currency currencies.Currency) *big.Float {
	rateToUsd, ok := cache.rates.toUsd[currency]

	if ok {
		return rateToUsd
	} else {
		return nil
	}
}
