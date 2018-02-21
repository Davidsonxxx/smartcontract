package database

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	dbBase "github.com/gameraccoon/telegram-bot-skeleton/database"
	"gitlab.com/gameraccoon/telegram-accountant-bot/currencies"
	"gitlab.com/gameraccoon/telegram-accountant-bot/wallettypes"
	"log"
	"sync"
)

type AccountDb struct {
	db dbBase.Database
	mutex sync.Mutex
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func ConnectDb(path string) (database *AccountDb, err error) {
	database = &AccountDb{}

	err = database.db.Connect(path)

	if err != nil {
		return
	}

	database.db.Exec("PRAGMA foreign_keys = ON")

	database.db.Exec("CREATE TABLE IF NOT EXISTS" +
		" global_vars(name TEXT PRIMARY KEY" +
		",integer_value INTEGER" +
		",string_value TEXT" +
		")")

	database.db.Exec("CREATE TABLE IF NOT EXISTS" +
		" users(id INTEGER NOT NULL PRIMARY KEY" +
		",chat_id INTEGER UNIQUE NOT NULL" +
		",language TEXT NOT NULL" +
		")")

	database.db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS" +
		" chat_id_index ON users(chat_id)")

	database.db.Exec("CREATE TABLE IF NOT EXISTS" +
		" wallets(id INTEGER NOT NULL PRIMARY KEY" +
		",is_removed INTEGER" + // NULL for alive wallets
		",user_id INTEGER NOT NULL" +
		",name STRING NOT NULL" +
		",currency INTEGER NOT NULL" +
		",address TEXT NOT NULL" +
		",type INTEGER NOT NULL" +
		",contract_address TEXT NOT NULL" + // not empty for ERC20 token wallets (currency == 5)
		",price_id TEXT NOT NULL" +
		",FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE SET NULL" +
		")")

	database.db.Exec("CREATE TABLE IF NOT EXISTS" +
		" rates(id INTEGER NOT NULL PRIMARY KEY" +
		",rate_to_usd REAL NOT NULL" +
		",time TIME NOT NULL" +
		")")

	return
}

func (database *AccountDb) IsConnectionOpened() bool {
	database.mutex.Lock()
	defer database.mutex.Unlock()

	return database.db.IsConnectionOpened()
}

func (database *AccountDb) Disconnect() {
	database.mutex.Lock()
	defer database.mutex.Unlock()

	database.db.Disconnect()
}

func (database *AccountDb) GetDatabaseVersion() (version string) {
	database.mutex.Lock()
	defer database.mutex.Unlock()

	rows, err := database.db.Query("SELECT string_value FROM global_vars WHERE name='version'")

	if err != nil {
		log.Fatal(err.Error())
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&version)
		if err != nil {
			log.Fatal(err.Error())
		}
	} else {
		// that means it's a new clean database
		version = latestVersion
	}

	return
}

func (database *AccountDb) SetDatabaseVersion(version string) {
	database.mutex.Lock()
	defer database.mutex.Unlock()

	database.db.Exec("DELETE FROM global_vars WHERE name='version'")

	safeVersion := dbBase.SanitizeString(version)
	database.db.Exec(fmt.Sprintf("INSERT INTO global_vars (name, string_value) VALUES ('version', '%s')", safeVersion))
}

func (database *AccountDb) GetUserId(chatId int64, userLangCode string) (userId int64) {
	database.mutex.Lock()
	defer database.mutex.Unlock()

	database.db.Exec(fmt.Sprintf("INSERT OR IGNORE INTO users(chat_id, language) "+
		"VALUES (%d, '%s')", chatId, userLangCode))

	rows, err := database.db.Query(fmt.Sprintf("SELECT id FROM users WHERE chat_id=%d", chatId))
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&userId)
		if err != nil {
			log.Fatal(err.Error())
		}
	} else {
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}
		log.Fatal("No user found")
	}

	return
}

func (database *AccountDb) GetUserChatId(userId int64) (chatId int64) {
	database.mutex.Lock()
	defer database.mutex.Unlock()

	rows, err := database.db.Query(fmt.Sprintf("SELECT chat_id FROM users WHERE id=%d", userId))
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&chatId)
		if err != nil {
			log.Fatal(err.Error())
		}
	} else {
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}
		log.Fatal("No user found")
	}

	return
}

func (database *AccountDb) GetUserWallets(userId int64) (ids []int64, names []string) {
	database.mutex.Lock()
	defer database.mutex.Unlock()

	rows, err := database.db.Query(fmt.Sprintf("SELECT id, name FROM wallets WHERE user_id=%d AND is_removed IS NULL", userId))
	if err != nil {
		log.Fatal(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var name string

		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err.Error())
		}

		ids = append(ids, id)
		names = append(names, name)
	}

	return
}

func (database *AccountDb) GetWalletName(walletId int64) (name string) {
	database.mutex.Lock()
	defer database.mutex.Unlock()

	rows, err := database.db.Query(fmt.Sprintf("SELECT name FROM wallets WHERE id=%d AND is_removed IS NULL", walletId))
	if err != nil {
		log.Fatal(err.Error())
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&name)
		if err != nil {
			log.Fatal(err.Error())
		}
	} else {
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}
		log.Fatal("No wallets found")
	}

	return
}

func (database *AccountDb) getLastInsertedItemId() (id int64) {
	rows, err := database.db.Query("SELECT last_insert_rowid()")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			log.Fatal(err.Error())
		}
		return
	} else {
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}
		log.Fatal("No item found")
	}
	return -1
}

func (database *AccountDb) CreateWatchOnlyWallet(userId int64, name string, address currencies.AddressData) (newWalletId int64) {
	database.mutex.Lock()
	defer database.mutex.Unlock()

	database.db.Exec(fmt.Sprintf(
		"INSERT INTO wallets(" +
		"user_id" +
		",name" +
		",currency" +
		",address" +
		",type" +
		",contract_address" +
		",price_id" +
		")VALUES(%d,'%s',%d,'%s',%d,'%s','%s')",
		userId,
		dbBase.SanitizeString(name),
		address.Currency,
		dbBase.SanitizeString(address.Address),
		wallettypes.WatchOnly,
		dbBase.SanitizeString(address.ContractAddress),
		dbBase.SanitizeString(address.PriceId),
	))

	return database.getLastInsertedItemId()
}

func (database *AccountDb) DeleteWallet(walletId int64) {
	database.mutex.Lock()
	defer database.mutex.Unlock()

	// give a way to recover things (don't delete completely)
	database.db.Exec(fmt.Sprintf("UPDATE OR ROLLBACK wallets SET is_removed=1 WHERE id=%d",  walletId))
}

func (database *AccountDb) IsWalletBelongsToUser(userId int64, walletId int64) bool {
	database.mutex.Lock()
	defer database.mutex.Unlock()

	rows, err := database.db.Query(fmt.Sprintf("SELECT COUNT(*) FROM wallets WHERE id=%d AND user_id=%d AND is_removed IS NULL", walletId, userId))
	if err != nil {
		log.Fatal(err.Error())
	}
	defer rows.Close()

	if rows.Next() {
		var count int
		err := rows.Scan(&count)
		if err != nil {
			log.Fatal(err.Error())
		} else {
			if count > 1 || count < 0 {
				log.Fatal("unique count of some walletId record is not 0 or 1")
			}

			if count >= 1 {
				return true
			}
		}
	}

	return false
}

func (database *AccountDb) SetUserLanguage(userId int64, language string) {
	database.mutex.Lock()
	defer database.mutex.Unlock()

	database.db.Exec(fmt.Sprintf("UPDATE OR ROLLBACK users SET language='%s' WHERE id=%d", language, userId))
}

func (database *AccountDb) GetUserLanguage(userId int64) (language string) {
	database.mutex.Lock()
	defer database.mutex.Unlock()

	rows, err := database.db.Query(fmt.Sprintf("SELECT language FROM users WHERE id=%d AND language IS NOT NULL", userId))
	if err != nil {
		log.Fatal(err.Error())
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&language)
		if err != nil {
			log.Fatal(err.Error())
		}
	} else {
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}
		// empty language
	}

	return
}

func (database *AccountDb) RenameWallet(walletId int64, newName string) {
	database.mutex.Lock()
	defer database.mutex.Unlock()

	database.db.Exec(fmt.Sprintf("UPDATE OR ROLLBACK wallets SET name='%s' WHERE id=%d AND is_removed IS NULL", newName, walletId))
}

func (database *AccountDb) GetWalletAddress(walletId int64) (addressData currencies.AddressData) {
	database.mutex.Lock()
	defer database.mutex.Unlock()

	rows, err := database.db.Query(fmt.Sprintf("SELECT currency, address, contract_address, price_id FROM wallets WHERE id=%d AND is_removed IS NULL LIMIT 1", walletId))
	if err != nil {
		log.Fatal(err.Error())
	}
	defer rows.Close()

	if rows.Next() {
		var currency int64
		var address string
		var contractAddress string
		var priceId string

		err := rows.Scan(&currency, &address, &contractAddress, &priceId)
		if err != nil {
			log.Fatal(err.Error())
		}

		addressData = currencies.AddressData{
			Currency: currencies.Currency(currency),
			Address: address,
			ContractAddress: contractAddress,
			PriceId: priceId,
		}
	} else {
		log.Fatalf("No wallet found with id %d", walletId)
	}

	return
}

func (database *AccountDb) GetUserWalletAddresses(userId int64) (addresses []currencies.AddressData) {
	database.mutex.Lock()
	defer database.mutex.Unlock()

	rows, err := database.db.Query(fmt.Sprintf("SELECT currency, address, contract_address, price_id FROM wallets WHERE user_id=%d AND is_removed IS NULL", userId))
	if err != nil {
		log.Fatal(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var currency int64
		var address string
		var contractAddress string
		var priceId string

		err := rows.Scan(&currency, &address, &contractAddress, &priceId)
		if err != nil {
			log.Fatal(err.Error())
		}

		addresses = append(
			addresses,
			currencies.AddressData{
				Currency: currencies.Currency(currency),
				Address: address,
				ContractAddress: contractAddress,
				PriceId: priceId,
			},
		)
	}

	return
}

func (database *AccountDb) GetAllWalletAddresses() (addresses []currencies.AddressData) {
	database.mutex.Lock()
	defer database.mutex.Unlock()

	rows, err := database.db.Query("SELECT currency, address, contract_address, price_id FROM wallets WHERE is_removed IS NULL")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var currency int64
		var address string
		var contractAddress string
		var priceId string

		err := rows.Scan(&currency, &address, &contractAddress, &priceId)
		if err != nil {
			log.Fatal(err.Error())
		}

		addresses = append(
			addresses,
			currencies.AddressData{
				Currency: currencies.Currency(currency),
				Address: address,
				ContractAddress: contractAddress,
				PriceId: priceId,
			},
		)
	}

	return
}

func (database *AccountDb) GetAllContractAddresses() (contractAddresses []string) {
	database.mutex.Lock()
	defer database.mutex.Unlock()

	rows, err := database.db.Query("SELECT DISTINCT contract_address FROM wallets WHERE is_removed IS NULL AND contract_address!=''")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var contractAddress string

		err := rows.Scan(&contractAddress)
		if err != nil {
			log.Fatal(err.Error())
		}

		contractAddresses = append(contractAddresses, contractAddress)
	}

	return
}

func (database *AccountDb) GetAllPriceIds() (priceIds []string) {
	database.mutex.Lock()
	defer database.mutex.Unlock()

	rows, err := database.db.Query("SELECT DISTINCT price_id FROM wallets WHERE is_removed IS NULL AND price_id!=''")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var priceId string

		err := rows.Scan(&priceId)
		if err != nil {
			log.Fatal(err.Error())
		}

		priceIds = append(priceIds, priceId)
	}

	return
}

func (database *AccountDb) SetWalletPriceId(walletId int64, priceId string) {
	database.mutex.Lock()
	defer database.mutex.Unlock()

	database.db.Exec(fmt.Sprintf("UPDATE OR ROLLBACK wallets SET price_id='%s' WHERE id=%d AND is_removed IS NULL", priceId, walletId))
}
