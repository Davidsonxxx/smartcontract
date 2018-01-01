package database

import (
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

func connectDb(t *testing.T) *Database {
	assert := require.New(t)
	db := &Database{}

	err := db.Connect(testDbPath)
	if err != nil {
		assert.Fail("Problem with creation db connection:" + err.Error())
		return nil
	}
	return db
}

func createDbAndConnect(t *testing.T) *Database {
	clearDb()
	return connectDb(t)
}

func TestConnection(t *testing.T) {
	assert := require.New(t)
	dropDatabase(testDbPath)

	db := &Database{}

	assert.False(db.IsConnectionOpened())

	err := db.Connect(testDbPath)
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

	db.SetDatabaseVersion(testText)
	assert.Equal(testText, db.GetDatabaseVersion())
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
		version := db.GetDatabaseVersion()
		assert.Equal(latestVersion, version)
	}

	{
		db.SetDatabaseVersion("1.0")
		version := db.GetDatabaseVersion()
		assert.Equal("1.0", version)
	}

	db.Disconnect()

	{
		db = connectDb(t)
		version := db.GetDatabaseVersion()
		assert.Equal("1.0", version)
		db.Disconnect()
	}

	{
		db = connectDb(t)
		db.SetDatabaseVersion("1.2")
		db.Disconnect()
	}

	{
		db = connectDb(t)
		version := db.GetDatabaseVersion()
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

	id1 := db.GetUserId(chatId1, "")
	id2 := db.GetUserId(chatId1, "")
	id3 := db.GetUserId(chatId2, "")

	assert.Equal(id1, id2)
	assert.NotEqual(id1, id3)

	assert.Equal(chatId1, db.GetUserChatId(id1))
	assert.Equal(chatId2, db.GetUserChatId(id3))
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
	userId := db.GetUserId(chatId, "")

	{
		ids, names := db.GetUserWallets(userId)
		assert.Equal(0, len(ids))
		assert.Equal(0, len(names))
	}

	walletId := db.CreateWatchOnlyWallet(userId, "testwallet", currencies.Bitcoin, "key")
	assert.True(db.IsWalletBelongsToUser(walletId, userId))
	{
		ids, names := db.GetUserWallets(userId)
		assert.Equal(1, len(ids))
		assert.Equal(1, len(names))
		if len(ids) > 0 && len(names) > 0 {
			assert.Equal(walletId, ids[0])
			assert.Equal("testwallet", names[0])
			assert.Equal("testwallet", db.GetWalletName(ids[0]))
		}
	}

	db.DeleteWallet(walletId)
	assert.False(db.IsWalletBelongsToUser(walletId, userId))
	{
		ids, names := db.GetUserWallets(userId)
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

	userId1 := db.GetUserId(123, "")
	userId2 := db.GetUserId(321, "")

	wallet1Id := db.CreateWatchOnlyWallet(userId1, "testwalt", currencies.Bitcoin, "key1")
	wallet2Id := db.CreateWatchOnlyWallet(userId2, "123", currencies.Bitcoin, "key2")

	assert.True(db.IsWalletBelongsToUser(userId1, wallet1Id))
	assert.True(db.IsWalletBelongsToUser(userId2, wallet2Id))
	assert.False(db.IsWalletBelongsToUser(userId1, wallet2Id))
	assert.False(db.IsWalletBelongsToUser(userId2, wallet1Id))
	// nonexistent wallet
	assert.False(db.IsWalletBelongsToUser(userId1, -1))
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

	userId1 := db.GetUserId(123, "")
	userId2 := db.GetUserId(321, "")

	db.SetUserLanguage(userId1, "en-US")

	{
		lang1 := db.GetUserLanguage(userId1)
		lang2 := db.GetUserLanguage(userId2)
		assert.Equal("en-US", lang1)
		assert.Equal("", lang2)
	}

	// in case of some side-effects
	{
		lang1 := db.GetUserLanguage(userId1)
		lang2 := db.GetUserLanguage(userId2)
		assert.Equal("en-US", lang1)
		assert.Equal("", lang2)
	}
}
