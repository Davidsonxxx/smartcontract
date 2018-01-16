package serverData

import (
	"github.com/gameraccoon/telegram-bot-skeleton/processing"
	"github.com/gameraccoon/telegram-bot-skeleton/database"
	ourDb "gitlab.com/gameraccoon/telegram-accountant-bot/database"
	"log"
)

type ServerDataManager struct {
	dataUpdater serverDataUpdater
}

func (serverDataManager *ServerDataManager) RegisterServerDataCache(staticData *processing.StaticProccessStructs) {
	if staticData == nil {
		log.Fatal("staticData is nil")
	}

	serverDataManager.dataUpdater.cache.Init()

	staticData.SetCustomValue("serverDataCache", &serverDataManager.dataUpdater.cache)
}

func (serverDataManager *ServerDataManager) updateAll(db *database.Database) {
	walletAddresses := ourDb.GetAllWalletAddresses(db)

	serverDataManager.dataUpdater.updateBalance(walletAddresses)
	serverDataManager.dataUpdater.updateRates()
}

func (serverDataManager *ServerDataManager) InitialUpdate(db *database.Database) {
	if db == nil {
		log.Fatal("database is nil")
		return
	}

	serverDataManager.updateAll(db)
}

func (serverDataManager *ServerDataManager) TimerTick(db *database.Database) {
	if db == nil {
		log.Print("database is nil, skip update")
		return
	}

	serverDataManager.updateAll(db)
}
