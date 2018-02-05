package staticFunctions

import (
	"github.com/gameraccoon/telegram-bot-skeleton/processing"
	"gitlab.com/gameraccoon/telegram-accountant-bot/database"
	"gitlab.com/gameraccoon/telegram-accountant-bot/currencies"
	"gitlab.com/gameraccoon/telegram-accountant-bot/serverData"
	static "gitlab.com/gameraccoon/telegram-accountant-bot/staticData"
	"github.com/nicksnyder/go-i18n/i18n"
	"log"
	"time"
)

func GetDb(staticData *processing.StaticProccessStructs) *database.AccountDb {
	if staticData == nil {
		log.Fatal("staticData is nil")
		return nil
	}

	db, ok := staticData.Db.(*database.AccountDb)
	if ok && db != nil {
		return db
	} else {
		log.Fatal("database is not set properly")
		return nil
	}
}

func FindTransFunction(userId int64, staticData *processing.StaticProccessStructs) i18n.TranslateFunc {
	// ToDo: cache user's lang
	lang := GetDb(staticData).GetUserLanguage(userId)

	config, configCastSuccess := staticData.Config.(static.StaticConfiguration)

	if !configCastSuccess {
		config = static.StaticConfiguration{}
	}

	// replace empty language to default one (some clients don't send user's language)
	if len(lang) <= 0 {
		log.Printf("User %d has empty language. Setting to default.", userId)
		lang = config.DefaultLanguage
		GetDb(staticData).SetUserLanguage(userId, lang)
	}

	if foundTrans, ok := staticData.Trans[lang]; ok {
		return foundTrans
	}

	// unknown language, use default instead
	if foundTrans, ok := staticData.Trans[config.DefaultLanguage]; ok {
		log.Printf("User %d has unknown language (%s). Setting to default.", userId, lang)
		lang = config.DefaultLanguage
		GetDb(staticData).SetUserLanguage(userId, lang)
		return foundTrans
	}

	// something gone wrong
	log.Printf("Translator didn't found: %s", lang)
	// fall to the first available translator
	for lang, trans := range staticData.Trans {
		log.Printf("Using first available translator: %s", lang)
		return trans
	}

	// something gone completely wrong
	log.Fatal("There are no available translators")
	// we will probably crash but there is nothing else we can do
	translator, _ := i18n.Tfunc(config.DefaultLanguage)
	return translator
}

func GetCurrencySymbolAndDecimals(serverData serverData.ServerDataInterface, currency currencies.Currency, contractAddress string) (currencySymbol string, currencyDecimals int) {
	if currency != currencies.Erc20Token {
		currencySymbol = currencies.GetCurrencySymbol(currency)
		currencyDecimals = currencies.GetCurrencyDecimals(currency)
	} else {
		if contractAddress == "" {
			log.Print("No contractAddress for token")
			return
		}

		tokenData := serverData.GetErc20TokenData(contractAddress)
		if tokenData != nil {
			currencySymbol = tokenData.Symbol
			currencyDecimals = tokenData.Decimals
		}
	}
	return
}

func FormatTimestamp(timestamp time.Time) string {
	loc, err := time.LoadLocation("EST")
	if err == nil {
		return timestamp.In(loc).Format(time.UnixDate)
	} else {
		return timestamp.Format(time.UnixDate)
	}
}
