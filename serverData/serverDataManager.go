package serverData

import (
	"github.com/gameraccoon/telegram-bot-skeleton/processing"
	"gitlab.com/gameraccoon/telegram-accountant-bot/database"
	"gitlab.com/gameraccoon/telegram-accountant-bot/currencies"
	"log"
	"math/big"
)

type TickUpdateData struct {
	BalanceNotifies []currencies.BalanceNotify
}

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

func (serverDataManager *ServerDataManager) updateBalanceNotifications(db *database.AccountDb, balanceChanges balanceChangesData) []currencies.BalanceNotify {
	walletIds := []int64{}

	for walletId, _ := range balanceChanges {
		walletIds = append(walletIds, walletId)
	}

	oldNotifiesData := db.GetBalanceNotifies(walletIds)

	notifiesToProcess := []currencies.BalanceNotify{}
	notifiesToInit := []currencies.BalanceNotify{}

	for _, notifyData := range oldNotifiesData {
		balance := balanceChanges[notifyData.WalletId]
		if balance != nil && notifyData.OldBalance != nil && balance.Cmp(notifyData.OldBalance) != 0 {
			notifyData.NewBalance = balance
			notifiesToProcess = append(notifiesToProcess, notifyData)
		}

		if notifyData.OldBalance == nil && balance != nil {
			notifyData.NewBalance = balance
			notifiesToInit = append(notifiesToInit, notifyData)
		}
	}

	if len(notifiesToProcess) > 0 {
		db.UpdateBalanceNotifies(notifiesToProcess)
	}

	if len(notifiesToInit) > 0 {
		db.UpdateBalanceNotifies(notifiesToInit)
	}

	return notifiesToProcess
}

func (serverDataManager *ServerDataManager) updateAll(db *database.AccountDb) []currencies.BalanceNotify {
	walletAddresses := db.GetAllWalletAddresses()
	priceIds := db.GetAllPriceIds()

	changedWalletIds := serverDataManager.dataUpdater.updateBalance(walletAddresses)

	serverDataManager.dataUpdater.updateRates(priceIds)

	return serverDataManager.updateBalanceNotifications(db, changedWalletIds)
}

func (serverDataManager *ServerDataManager) InitialUpdate(db *database.AccountDb) TickUpdateData {
	if db == nil {
		log.Fatal("database is nil")
		return TickUpdateData{}
	}

	balanceNotifies := serverDataManager.updateAll(db)

	contractsIds := db.GetAllContractAddresses()
	serverDataManager.dataUpdater.updateErc20TokensData(contractsIds)

	return TickUpdateData {
		BalanceNotifies: balanceNotifies,
	}
}

func (serverDataManager *ServerDataManager) TimerTick(db *database.AccountDb) TickUpdateData {
	if db == nil {
		log.Print("database is nil, skip update")
		return TickUpdateData{}
	}

	balanceNotifies := serverDataManager.updateAll(db)

	return TickUpdateData {
		BalanceNotifies: balanceNotifies,
	}
}

func (serverDataManager *ServerDataManager) GetBalance(address currencies.AddressData) *big.Int {
	balance := serverDataManager.dataUpdater.cache.getBalance(address)

	if balance != nil {
		return balance
	} else {
		return serverDataManager.dataUpdater.updateBalanceOneWallet(address)
	}
}

func (serverDataManager *ServerDataManager) GetRateToUsd(priceId string) *big.Float {
	return serverDataManager.dataUpdater.cache.getRateToUsd(priceId)
}

func (serverDataManager *ServerDataManager) GetErc20TokenData(contractAddress string) *currencies.Erc20TokenData {
	tokenData := serverDataManager.dataUpdater.cache.getErc20TokenData(contractAddress)
	if tokenData == nil {
		tokenData = serverDataManager.dataUpdater.updateOneErc20TokensData(contractAddress)
	}

	return tokenData
}
