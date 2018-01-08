package database

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/gameraccoon/telegram-bot-skeleton/database"
	"gitlab.com/gameraccoon/telegram-accountant-bot/currencies"
	"gitlab.com/gameraccoon/telegram-accountant-bot/wallettypes"
	"log"
)


func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func Init(path string) (db *database.Database, err error) {
	db = &database.Database{}

	err = db.Connect(path)

	if err != nil {
		return
	}

	db.Exec("PRAGMA foreign_keys = ON")

	db.Exec("CREATE TABLE IF NOT EXISTS" +
		" global_vars(name TEXT PRIMARY KEY" +
		",integer_value INTEGER" +
		",string_value TEXT" +
		")")

	db.Exec("CREATE TABLE IF NOT EXISTS" +
		" users(id INTEGER NOT NULL PRIMARY KEY" +
		",chat_id INTEGER UNIQUE NOT NULL" +
		",language TEXT NOT NULL" +
		")")

	db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS" +
		" chat_id_index ON users(chat_id)")

	db.Exec("CREATE TABLE IF NOT EXISTS" +
		" wallets(id INTEGER NOT NULL PRIMARY KEY" +
		",is_removed INTEGER" + // NULL for alive wallets
		",user_id INTEGER NOT NULL" +
		",name STRING NOT NULL" +
		",currency INTEGER NOT NULL" +
		",address TEXT NOT NULL" +
		",type INTEGER NOT NULL" +
		",balance INTEGER" + // only for virtual wallets
		",private_key_storage TEXT" + // NULL for watch-only wallets
		",FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE SET NULL" +
		")")

	db.Exec("CREATE TABLE IF NOT EXISTS" +
		" rates(id INTEGER NOT NULL PRIMARY KEY" +
		",rate_to_usd REAL NOT NULL" +
		",time TIME NOT NULL" +
		")")

	return
}

func GetDatabaseVersion(db *database.Database) (version string) {
	rows, err := db.Query("SELECT string_value FROM global_vars WHERE name='version'")

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

func SetDatabaseVersion(db *database.Database, version string) {
	db.Exec("DELETE FROM global_vars WHERE name='version'")

	safeVersion := database.SanitizeString(version)
	db.Exec(fmt.Sprintf("INSERT INTO global_vars (name, string_value) VALUES ('version', '%s')", safeVersion))
}

func GetUserId(db *database.Database, chatId int64, userLangCode string) (userId int64) {
	db.Exec(fmt.Sprintf("INSERT OR IGNORE INTO users(chat_id, language) "+
		"VALUES (%d, '%s')", chatId, userLangCode))

	rows, err := db.Query(fmt.Sprintf("SELECT id FROM users WHERE chat_id=%d", chatId))
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

func GetUserChatId(db *database.Database, userId int64) (chatId int64) {
	rows, err := db.Query(fmt.Sprintf("SELECT chat_id FROM users WHERE id=%d", userId))
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

func GetUserWallets(db *database.Database, userId int64) (ids []int64, names []string) {
	rows, err := db.Query(fmt.Sprintf("SELECT id, name FROM wallets WHERE user_id=%d AND is_removed IS NULL", userId))
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

func GetWalletName(db *database.Database, walletId int64) (name string) {
	rows, err := db.Query(fmt.Sprintf("SELECT name FROM wallets WHERE id=%d AND is_removed IS NULL", walletId))
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

func GetLastInsertedItemId(db *database.Database) (id int64) {
	rows, err := db.Query("SELECT last_insert_rowid()")
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

func CreateWatchOnlyWallet(db *database.Database, userId int64, name string, currency currencies.Currency, address string) (newWalletId int64) {
	db.Exec(fmt.Sprintf(
		"INSERT INTO wallets(" +
		"user_id" +
		",name" +
		",currency" +
		",address" +
		",type" +
		")VALUES(%d,'%s',%d,'%s',%d)",
		userId,
		database.SanitizeString(name),
		currency,
		database.SanitizeString(address),
		wallettypes.WatchOnly,
	))

	return GetLastInsertedItemId(db)
}

func CreateFullWallet(db *database.Database, userId int64, name string, currency currencies.Currency, address string, privateKey string) (newWalletId int64) {
	db.Exec(fmt.Sprintf(
		"INSERT INTO wallets(" +
		"user_id" +
		",name" +
		",currency" +
		",address" +
		",type" +
		",private_key_storage" +
		")VALUES(%d,'%s',%d,'%s',%d,'%s')",
		userId,
		database.SanitizeString(name),
		currency,
		database.SanitizeString(address),
		wallettypes.Full,
		database.SanitizeString(privateKey),
	))

	return GetLastInsertedItemId(db)
}

func CreateVirtualWallet(db *database.Database, userId int64, name string, currency currencies.Currency, address string) (newWalletId int64) {
	db.Exec(fmt.Sprintf(
		"INSERT INTO wallets(" +
		"user_id" +
		",name" +
		",currency" +
		",address" +
		",type" +
		",balance" +
		")VALUES(%d,'%s',%s,'%s',%d,%d)",
		userId,
		database.SanitizeString(name),
		currency,
		database.SanitizeString(address),
		wallettypes.Virtual,
		0, // init with zero balance
	))

	return GetLastInsertedItemId(db)
}

func DeleteWallet(db *database.Database, walletId int64) {
	// give a way to recover things (don't delete completely)
	db.Exec(fmt.Sprintf("UPDATE OR ROLLBACK wallets SET is_removed=1 WHERE id=%d",  walletId))
}

func IsWalletBelongsToUser(db *database.Database, userId int64, walletId int64) bool {
	rows, err := db.Query(fmt.Sprintf("SELECT COUNT(*) FROM wallets WHERE id=%d AND user_id=%d AND is_removed IS NULL", walletId, userId))
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

func SetUserLanguage(db *database.Database, userId int64, language string) {
	db.Exec(fmt.Sprintf("UPDATE OR ROLLBACK users SET language='%s' WHERE id=%d", language, userId))
}

func GetUserLanguage(db *database.Database, userId int64) (language string) {
	rows, err := db.Query(fmt.Sprintf("SELECT language FROM users WHERE id=%d AND language IS NOT NULL", userId))
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

func RenameWallet(db *database.Database, walletId int64, newName string) {
	db.Exec(fmt.Sprintf("UPDATE OR ROLLBACK wallets SET name='%s' WHERE id=%d AND is_removed IS NULL", newName, walletId))
}

func GetWalletAddress(db *database.Database, walletId int64) (addressData currencies.AddressData) {
	rows, err := db.Query(fmt.Sprintf("SELECT currency, address FROM wallets WHERE id=%d AND is_removed IS NULL LIMIT 1", walletId))
	if err != nil {
		log.Fatal(err.Error())
	}
	defer rows.Close()

	if rows.Next() {
		var currency int64
		var address string

		err := rows.Scan(&currency, &address)
		if err != nil {
			log.Fatal(err.Error())
		}

		addressData = currencies.AddressData{
			Currency: currencies.Currency(currency),
			Address: address,
		}
	} else {
		log.Fatalf("No wallet found with id %d", walletId)
	}

	return
}

func GetUserWalletAddresses(db *database.Database, userId int64) (addresses []currencies.AddressData) {
	rows, err := db.Query(fmt.Sprintf("SELECT currency, address FROM wallets WHERE user_id=%d AND is_removed IS NULL", userId))
	if err != nil {
		log.Fatal(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var currency int64
		var address string

		err := rows.Scan(&currency, &address)
		if err != nil {
			log.Fatal(err.Error())
		}

		addresses = append(
			addresses,
			currencies.AddressData{
				Currency: currencies.Currency(currency),
				Address: address,
			},
		)
	}

	return
}
