package serverData

import (
	"github.com/gameraccoon/telegram-bot-skeleton/processing"
	"gitlab.com/gameraccoon/telegram-accountant-bot/database"
	"gitlab.com/gameraccoon/telegram-accountant-bot/cryptoFunctions"
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

func fillLastTransactionsData(notify *currencies.BalanceNotify, transactions []currencies.TransactionsHistoryItem) {
	if notify == nil {
		return
	}

	// fill transactions only if we're going to show the notify
	if !notify.IsInitialChange {
		lastKnownItemIdx := -1

		for i, transaction := range transactions {
			if transaction.Time.Equal(notify.OldTransactionTime) {
				lastKnownItemIdx = i
				break
			}
		}

		if lastKnownItemIdx != -1 {
			notify.LastTransactions = transactions[:lastKnownItemIdx]
		} else {
			notify.LastTransactions = transactions
		}
	}

	if len(transactions) > 0 {
		notify.NewTransactionTime = transactions[0].Time
	}
}

func (serverDataManager *ServerDataManager) updateBalanceNotifications(db *database.AccountDb, balanceChanges balanceChangesData) []currencies.BalanceNotify {
	walletIds := []int64{}

	for walletId, _ := range balanceChanges {
		walletIds = append(walletIds, walletId)
	}

	oldNotifiesData := db.GetBalanceNotifies(walletIds)

	notifiesToProcess := []currencies.BalanceNotify{}

	// check all the notifies
	for _, notifyData := range oldNotifiesData {
		balance := balanceChanges[notifyData.WalletId]
		// if we have new balance
		if balance != nil {
			notifyData.NewBalance = balance

			// if the balance is changed from the last DB record
			if notifyData.OldBalance == nil || balance.Cmp(notifyData.OldBalance) != 0 {
				if notifyData.OldBalance == nil {
					// if the record is new, mark that we don't want to show the notify
					notifyData.IsInitialChange = true;
				}

				walletAddress := db.GetWalletAddress(notifyData.WalletId)
				notifyData.WalletAddress = walletAddress

				processor := cryptoFunctions.GetProcessor(walletAddress.Currency)

				if processor != nil {
					lastTransactions := (*processor).GetTransactionsHistory(walletAddress, 10)

					fillLastTransactionsData(&notifyData, lastTransactions)
				}

				notifiesToProcess = append(notifiesToProcess, notifyData)
			}
		}
	}

	// write new values to the DB
	if len(notifiesToProcess) > 0 {
		db.UpdateBalanceNotifies(notifiesToProcess)
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
