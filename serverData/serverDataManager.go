package serverData

import (
	"github.com/gameraccoon/telegram-bot-skeleton/processing"
	"github.com/gameraccoon/telegram-bot-skeleton/database"
	ourDb "gitlab.com/gameraccoon/telegram-accountant-bot/database"
	"gitlab.com/gameraccoon/telegram-accountant-bot/currencies"
	"log"
	"math/big"
	"sync"
)

type ServerDataManager struct {
	dataUpdater serverDataUpdater
}

func GetServerData(staticData *processing.StaticProccessStructs) ServerDataInterface {
	if staticData == nil {
		log.Print("staticData is nil")
		return nil
	}

	dataCache, ok := staticData.GetCustomValue("serverDataInterface").(ServerDataInterface)
	if ok {
		return dataCache
	} else {
		log.Fatal("serverDataInterface is not set properly")
		return nil
	}
}

func (serverDataManager *ServerDataManager) RegisterServerDataInterface(staticData *processing.StaticProccessStructs) {
	if staticData == nil {
		log.Fatal("staticData is nil")
	}

	serverDataManager.dataUpdater.cache.Init()

	var serverDataInterface ServerDataInterface = serverDataManager

	if serverDataInterface != nil {
		staticData.SetCustomValue("serverDataInterface", serverDataInterface)
	} else {
		log.Fatal("ServerDataManager does not implement ServerDataInterface")
	}
}

func (serverDataManager *ServerDataManager) updateAll(db *database.Database, dbMutex *sync.Mutex) {
	dbMutex.Lock()
	walletAddresses := ourDb.GetAllWalletAddresses(db)
	dbMutex.Unlock()

	serverDataManager.dataUpdater.updateBalance(walletAddresses)
	serverDataManager.dataUpdater.updateRates()
}

func (serverDataManager *ServerDataManager) InitialUpdate(db *database.Database, dbMutex *sync.Mutex) {
	if db == nil {
		log.Fatal("database is nil")
		return
	}

	serverDataManager.updateAll(db, dbMutex)
}

func (serverDataManager *ServerDataManager) TimerTick(db *database.Database, dbMutex *sync.Mutex) {
	if db == nil {
		log.Print("database is nil, skip update")
		return
	}

	serverDataManager.updateAll(db, dbMutex)
}

func (serverDataManager *ServerDataManager) GetBalance(address currencies.AddressData) *big.Int {
	balance := serverDataManager.dataUpdater.cache.GetBalance(address)

	if balance != nil {
		return balance
	} else {
		return serverDataManager.dataUpdater.updateBalanceOneWallet(address)
	}
}

func (serverDataManager *ServerDataManager) GetRateToUsd(currency currencies.Currency) *big.Float {
	return serverDataManager.dataUpdater.cache.GetRateToUsd(currency)
}
