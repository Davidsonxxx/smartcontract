package database

import (
	"github.com/gameraccoon/telegram-bot-skeleton/database"
	"github.com/stretchr/testify/require"
	"gitlab.com/gameraccoon/telegram-accountant-bot/currencies"
	"os"
	"testing"
)

const (
	testDbPath = "./testDb.db"
)

func dropDatabase(fileName string) {
	os.Remove(fileName)
}

func clearDb() {
	dropDatabase(testDbPath)
}

func connectDb(t *testing.T) *database.Database {
	assert := require.New(t)
	db, err := Init(testDbPath)

	if err != nil {
		assert.Fail("Problem with creation db connection:" + err.Error())
		return nil
	}
	return db
}

func createDbAndConnect(t *testing.T) *database.Database {
	clearDb()
	return connectDb(t)
}

func TestConnection(t *testing.T) {
	assert := require.New(t)
	dropDatabase(testDbPath)

	db, err := Init(testDbPath)

	defer dropDatabase(testDbPath)
	if err != nil {
		assert.Fail("Problem with creation db connection:" + err.Error())
		return
	}

	assert.True(db.IsConnectionOpened())

	db.Disconnect()

	assert.False(db.IsConnectionOpened())
}

func TestSanitizeString(t *testing.T) {
	assert := require.New(t)
	db := createDbAndConnect(t)
	defer clearDb()
	if db == nil {
		t.Fail()
		return
	}
	defer db.Disconnect()

	testText := "text'test''test\"test\\"

	SetDatabaseVersion(db, testText)
	assert.Equal(testText, GetDatabaseVersion(db))
}

func TestDatabaseVersion(t *testing.T) {
	assert := require.New(t)
	db := createDbAndConnect(t)
	defer clearDb()
	if db == nil {
		t.Fail()
		return
	}

	{
		version := GetDatabaseVersion(db)
		assert.Equal(latestVersion, version)
	}

	{
		SetDatabaseVersion(db, "1.0")
		version := GetDatabaseVersion(db)
		assert.Equal("1.0", version)
	}

	db.Disconnect()

	{
		db = connectDb(t)
		version := GetDatabaseVersion(db)
		assert.Equal("1.0", version)
		db.Disconnect()
	}

	{
		db = connectDb(t)
		SetDatabaseVersion(db, "1.2")
		db.Disconnect()
	}

	{
		db = connectDb(t)
		version := GetDatabaseVersion(db)
		assert.Equal("1.2", version)
		db.Disconnect()
	}
}

func TestGetUserId(t *testing.T) {
	assert := require.New(t)
	db := createDbAndConnect(t)
	defer clearDb()
	if db == nil {
		t.Fail()
		return
	}
	defer db.Disconnect()

	var chatId1 int64 = 321
	var chatId2 int64 = 123

	id1 := GetUserId(db, chatId1, "")
	id2 := GetUserId(db, chatId1, "")
	id3 := GetUserId(db, chatId2, "")

	assert.Equal(id1, id2)
	assert.NotEqual(id1, id3)

	assert.Equal(chatId1, GetUserChatId(db, id1))
	assert.Equal(chatId2, GetUserChatId(db, id3))
}

func TestCreateAndRemoveWallet(t *testing.T) {
	assert := require.New(t)
	db := createDbAndConnect(t)
	defer clearDb()
	if db == nil {
		t.Fail()
		return
	}
	defer db.Disconnect()

	var chatId int64 = 123
	userId := GetUserId(db, chatId, "")

	{
		ids, names := GetUserWallets(db, userId)
		assert.Equal(0, len(ids))
		assert.Equal(0, len(names))
	}

	walletId := CreateWatchOnlyWallet(db, userId, "testwallet", currencies.Bitcoin, "key")
	assert.True(IsWalletBelongsToUser(db, walletId, userId))
	{
		ids, names := GetUserWallets(db, userId)
		assert.Equal(1, len(ids))
		assert.Equal(1, len(names))
		if len(ids) > 0 && len(names) > 0 {
			assert.Equal(walletId, ids[0])
			assert.Equal("testwallet", names[0])
			assert.Equal("testwallet", GetWalletName(db, ids[0]))
		}
	}

	DeleteWallet(db, walletId)
	assert.False(IsWalletBelongsToUser(db, walletId, userId))
	{
		ids, names := GetUserWallets(db, userId)
		assert.Equal(0, len(ids))
		assert.Equal(0, len(names))
	}
}

func TestWalletBelongsToUser(t *testing.T) {
	assert := require.New(t)
	db := createDbAndConnect(t)
	defer clearDb()
	if db == nil {
		t.Fail()
		return
	}
	defer db.Disconnect()

	userId1 := GetUserId(db, 123, "")
	userId2 := GetUserId(db, 321, "")

	wallet1Id := CreateWatchOnlyWallet(db, userId1, "testwalt", currencies.Bitcoin, "key1")
	wallet2Id := CreateWatchOnlyWallet(db, userId2, "123", currencies.Bitcoin, "key2")

	assert.True(IsWalletBelongsToUser(db, userId1, wallet1Id))
	assert.True(IsWalletBelongsToUser(db, userId2, wallet2Id))
	assert.False(IsWalletBelongsToUser(db, userId1, wallet2Id))
	assert.False(IsWalletBelongsToUser(db, userId2, wallet1Id))
	// nonexistent wallet
	assert.False(IsWalletBelongsToUser(db, userId1, -1))

	walletAddresses := GetAllWalletAddresses(db)
	assert.Equal(2, len(walletAddresses))
}

func TestUsersLanguage(t *testing.T) {
	assert := require.New(t)
	db := createDbAndConnect(t)
	defer clearDb()
	if db == nil {
		t.Fail()
		return
	}
	defer db.Disconnect()

	userId1 := GetUserId(db, 123, "")
	userId2 := GetUserId(db, 321, "")

	SetUserLanguage(db, userId1, "en-US")

	{
		lang1 := GetUserLanguage(db, userId1)
		lang2 := GetUserLanguage(db, userId2)
		assert.Equal("en-US", lang1)
		assert.Equal("", lang2)
	}

	// in case of some side-effects
	{
		lang1 := GetUserLanguage(db, userId1)
		lang2 := GetUserLanguage(db, userId2)
		assert.Equal("en-US", lang1)
		assert.Equal("", lang2)
	}
}

func TestWalletRenaming(t *testing.T) {
	assert := require.New(t)
	db := createDbAndConnect(t)
	defer clearDb()
	if db == nil {
		t.Fail()
		return
	}
	defer db.Disconnect()

	userId := GetUserId(db, 123, "")

	walletId := CreateWatchOnlyWallet(db, userId, "testwallet", currencies.Bitcoin, "key1")

	{
		ids, names := GetUserWallets(db, userId)
		if len(ids) > 0 && len(names) > 0 {
			assert.Equal(walletId, ids[0])
			assert.Equal("testwallet", names[0])
		}
	}

	RenameWallet(db, walletId, "test2")

	{
		ids, names := GetUserWallets(db, userId)
		if len(ids) > 0 && len(names) > 0 {
			assert.Equal(walletId, ids[0])
			assert.Equal("test2", names[0])
		}
	}
}

func TestGettingWalletAddresses(t *testing.T) {
	assert := require.New(t)
	db := createDbAndConnect(t)
	defer clearDb()
	if db == nil {
		t.Fail()
		return
	}
	defer db.Disconnect()

	userId1 := GetUserId(db, 123, "")
	userId2 := GetUserId(db, 321, "")

	walletId1 := CreateWatchOnlyWallet(db, userId1, "testwallet1", currencies.Bitcoin, "adr1")
	walletId2 := CreateWatchOnlyWallet(db, userId1, "testwallet2", currencies.Ether, "adr2")
	walletId3 := CreateWatchOnlyWallet(db, userId2, "testwallet3", currencies.Bitcoin, "adr3")

	{
		addr1 := GetWalletAddress(db, walletId1)
		addr2 := GetWalletAddress(db, walletId2)
		addr3 := GetWalletAddress(db, walletId3)

		assert.Equal("adr1", addr1.Address)
		assert.Equal(currencies.Bitcoin, addr1.Currency)
		assert.Equal("adr2", addr2.Address)
		assert.Equal(currencies.Ether, addr2.Currency)
		assert.Equal("adr3", addr3.Address)
		assert.Equal(currencies.Bitcoin, addr3.Currency)
	}

	{
		addresses := GetUserWalletAddresses(db, userId1)
		assert.Equal(2, len(addresses))
		for _, address := range addresses {
			if address.Currency == currencies.Bitcoin {
				assert.Equal("adr1", address.Address)
			} else if address.Currency == currencies.Ether {
				assert.Equal("adr2", address.Address)
			} else {
				assert.Fail("Unexpected currency type")
			}
		}
	}
}
