package serverData

import (
	"github.com/gameraccoon/telegram-bot-skeleton/processing"
	"gitlab.com/gameraccoon/telegram-accountant-bot/currencies"
	"math/big"
	"log"
)

type ServerDataManager struct {
	cache dataCache
}

func GetServerDataManager(staticData *processing.StaticProccessStructs) ServerDataManager {
	if staticData == nil {
		log.Fatal("staticData is nil")
	}

	serverDataManager, ok := staticData.GetCustomValue("serverDataManager").(ServerDataManager)
	if ok {
		return serverDataManager
	} else {
		serverDataManager := ServerDataManager{}
		staticData.SetCustomValue("serverDataManager", serverDataManager)
		return serverDataManager
	}
}

func (serverDataManager *ServerDataManager) GetBalance(address currencies.AddressData) *big.Int {
	balance, ok := serverDataManager.cache.balances[address]
	if ok {
		return balance
	} else {
		return nil
	}
}

func (serverDataManager *ServerDataManager) GetRateToUsd(currency currencies.Currency) *big.Float {
	rateToUsd, ok := serverDataManager.cache.rates.toUsd[currency]

	if ok {
		return rateToUsd
	} else {
		return nil
	}
}

func (serverDataManager *ServerDataManager) CalcUsdBalance(address currencies.AddressData) *big.Int {
	balance, ok := serverDataManager.cache.balances[address]
	if ok {
		return balance
	} else {
		return nil
	}
}
