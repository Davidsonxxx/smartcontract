package database

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"gitlab.com/gameraccoon/telegram-accountant-bot/currencies"
	"gitlab.com/gameraccoon/telegram-accountant-bot/wallettypes"
	"log"
	"strings"
)


func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

type Database struct {
	// connection
	conn *sql.DB
}

func sanitizeString(input string) string {
	return strings.Replace(input, "'", "''", -1)
}

func (database *Database) execQuery(query string) {
	_, err := database.conn.Exec(query)

	if err != nil {
		log.Fatal(err.Error())
	}
}

func (database *Database) Connect(fileName string) error {
	db, err := sql.Open("sqlite3", fileName)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}

	database.conn = db

	database.execQuery("PRAGMA foreign_keys = ON")

	database.execQuery("CREATE TABLE IF NOT EXISTS" +
		" global_vars(name TEXT PRIMARY KEY" +
		",integer_value INTEGER" +
		",string_value TEXT" +
		")")

	database.execQuery("CREATE TABLE IF NOT EXISTS" +
		" users(id INTEGER NOT NULL PRIMARY KEY" +
		",chat_id INTEGER UNIQUE NOT NULL" +
		",language TEXT NOT NULL" +
		")")

	database.execQuery("CREATE UNIQUE INDEX IF NOT EXISTS" +
		" chat_id_index ON users(chat_id)")

	database.execQuery("CREATE TABLE IF NOT EXISTS" +
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

	database.execQuery("CREATE TABLE IF NOT EXISTS" +
		" rates(id INTEGER NOT NULL PRIMARY KEY" +
		",rate_to_usd REAL NOT NULL" +
		",time TIME NOT NULL" +
		")")

	return nil
}

func (database *Database) Disconnect() {
	database.conn.Close()
	database.conn = nil
}

func (database *Database) IsConnectionOpened() bool {
	return database.conn != nil
}

func (database *Database) createUniqueRecord(table string, values string) int64 {
	var err error
	if len(values) == 0 {
		_, err = database.conn.Exec(fmt.Sprintf("INSERT INTO %s DEFAULT VALUES ", table))
	} else {
		_, err = database.conn.Exec(fmt.Sprintf("INSERT INTO %s VALUES (%s)", table, values))
	}

	if err != nil {
		log.Fatal(err.Error())
		return -1
	}

	rows, err := database.conn.Query(fmt.Sprintf("SELECT id FROM %s ORDER BY id DESC LIMIT 1", table))

	if err != nil {
		log.Fatal(err.Error())
		return -1
	}
	defer rows.Close()

	if rows.Next() {
		var id int64
		err := rows.Scan(&id)
		if err != nil {
			log.Fatal(err.Error())
			return -1
		}

		return id
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal("No record created")
	return -1
}

func (database *Database) GetDatabaseVersion() (version string) {
	rows, err := database.conn.Query("SELECT string_value FROM global_vars WHERE name='version'")

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

func (database *Database) SetDatabaseVersion(version string) {
	database.execQuery("DELETE FROM global_vars WHERE name='version'")

	safeVersion := sanitizeString(version)
	database.execQuery(fmt.Sprintf("INSERT INTO global_vars (name, string_value) VALUES ('version', '%s')", safeVersion))
}

func (database *Database) GetUserId(chatId int64, userLangCode string) (userId int64) {
	database.execQuery(fmt.Sprintf("INSERT OR IGNORE INTO users(chat_id, language) "+
		"VALUES (%d, '%s')", chatId, userLangCode))

	rows, err := database.conn.Query(fmt.Sprintf("SELECT id FROM users WHERE chat_id=%d", chatId))
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

func (database *Database) GetUserChatId(userId int64) (chatId int64) {
	rows, err := database.conn.Query(fmt.Sprintf("SELECT chat_id FROM users WHERE id=%d", userId))
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

func (database *Database) GetUserWallets(userId int64) (ids []int64, names []string) {
	rows, err := database.conn.Query(fmt.Sprintf("SELECT id, name FROM wallets WHERE user_id=%d AND is_removed IS NULL", userId))
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

func (database *Database) GetWalletName(walletId int64) (name string) {
	rows, err := database.conn.Query(fmt.Sprintf("SELECT name FROM wallets WHERE id=%d AND is_removed IS NULL", walletId))
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

func (database* Database) GetLastInsertedItemId() (id int64) {
	rows, err := database.conn.Query("SELECT last_insert_rowid()")
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

func (database *Database) CreateWatchOnlyWallet(userId int64, name string, currency currencies.Currency, address string) (newWalletId int64) {
	database.execQuery(fmt.Sprintf(
		"INSERT INTO wallets(" +
		"user_id" +
		",name" +
		",currency" +
		",address" +
		",type" +
		")VALUES(%d,'%s',%d,'%s',%d)",
		userId,
		sanitizeString(name),
		currency,
		sanitizeString(address),
		wallettypes.WatchOnly,
	))

	return database.GetLastInsertedItemId()
}

func (database *Database) CreateFullWallet(userId int64, name string, currency currencies.Currency, address string, privateKey string) (newWalletId int64) {
	database.execQuery(fmt.Sprintf(
		"INSERT INTO wallets(" +
		"user_id" +
		",name" +
		",currency" +
		",address" +
		",type" +
		",private_key_storage" +
		")VALUES(%d,'%s',%d,'%s',%d,'%s')",
		userId,
		sanitizeString(name),
		currency,
		sanitizeString(address),
		wallettypes.Full,
		sanitizeString(privateKey),
	))

	return database.GetLastInsertedItemId()
}

func (database *Database) CreateVirtualWallet(userId int64, name string, currency currencies.Currency, address string) (newWalletId int64) {
	database.execQuery(fmt.Sprintf(
		"INSERT INTO wallets(" +
		"user_id" +
		",name" +
		",currency" +
		",address" +
		",type" +
		",balance" +
		")VALUES(%d,'%s',%s,'%s',%d,%d)",
		userId,
		sanitizeString(name),
		currency,
		sanitizeString(address),
		wallettypes.Virtual,
		0, // init with zero balance
	))

	return database.GetLastInsertedItemId()
}

func (database *Database) DeleteWallet(walletId int64) {
	// give a way to recover things (don't delete completely)
	database.execQuery(fmt.Sprintf("UPDATE OR ROLLBACK wallets SET is_removed=1 WHERE id=%d",  walletId))
}

func (database *Database) IsWalletBelongsToUser(userId int64, walletId int64) bool {
	rows, err := database.conn.Query(fmt.Sprintf("SELECT COUNT(*) FROM wallets WHERE id=%d AND user_id=%d AND is_removed IS NULL", walletId, userId))
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

func (database *Database) SetUserLanguage(userId int64, language string) {
	database.execQuery(fmt.Sprintf("UPDATE OR ROLLBACK users SET language='%s' WHERE id=%d", language, userId))
}

func (database *Database) GetUserLanguage(userId int64) (language string) {
	rows, err := database.conn.Query(fmt.Sprintf("SELECT language FROM users WHERE id=%d AND language IS NOT NULL", userId))
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

func (database *Database) RenameWallet(walletId int64, newName string) {
	database.execQuery(fmt.Sprintf("UPDATE OR ROLLBACK wallets SET name='%s' WHERE id=%d AND is_removed IS NULL", newName, walletId))
}

func (database *Database) GetWalletAddress(walletId int64) (addressData currencies.AddressData) {
	rows, err := database.conn.Query(fmt.Sprintf("SELECT currency, address FROM wallets WHERE id=%d AND is_removed IS NULL LIMIT 1", walletId))
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

func (database *Database) GetUserWalletAddresses(userId int64) (addresses []currencies.AddressData) {
	rows, err := database.conn.Query(fmt.Sprintf("SELECT currency, address FROM wallets WHERE user_id=%d AND is_removed IS NULL", userId))
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
