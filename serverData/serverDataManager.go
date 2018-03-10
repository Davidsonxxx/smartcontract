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

	newNotifies := []currencies.BalanceNotify{}

	for _, notify := range oldNotifiesData {
		balance := balanceChanges[notify.WalletId]
		if balance != nil && balance.Cmp(notify.OldBalance) != 0 {
			notify.NewBalance = balance

			newNotifies = append(newNotifies, notify)
		}
	}

	db.UpdateBalanceNotifies(newNotifies)

	return newNotifies
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
