package staticFunctions

import (
	"github.com/gameraccoon/telegram-bot-skeleton/processing"
	"gitlab.com/gameraccoon/telegram-accountant-bot/database"
	"github.com/nicksnyder/go-i18n/i18n"
	"log"
	"math/big"
)

func FindTransFunction(userId int64, staticData *processing.StaticProccessStructs) i18n.TranslateFunc {
	// ToDo: cache user's lang
	lang := database.GetUserLanguage(staticData.Db, userId)

	// replace empty language to default one (some clients don't send user's language)
	if len(lang) <= 0 {
		log.Printf("User %d has empty language. Setting to default.", userId)
		lang = staticData.Config.DefaultLanguage
		database.SetUserLanguage(staticData.Db, userId, lang)
	}

	if foundTrans, ok := staticData.Trans[lang]; ok {
		return foundTrans
	}

	// unknown language, use default instead
	if foundTrans, ok := staticData.Trans[staticData.Config.DefaultLanguage]; ok {
		log.Printf("User %d has unknown language (%s). Setting to default.", userId, lang)
		lang = staticData.Config.DefaultLanguage
		database.SetUserLanguage(staticData.Db, userId, lang)
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
	translator, _ := i18n.Tfunc(staticData.Config.DefaultLanguage)
	return translator
}

func FormatCurrencyAmount(intValue *big.Int, digits int) string {
	// isn't it a piece of shit?
	// balanceFloat = balance / (10.0 ** currencyDigits)
	var balanceFloat *big.Float = new(big.Float).Quo(new(big.Float).SetInt(intValue), new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(digits)), big.NewInt(0))))

	if balanceFloat == nil {
		return "Error"
	}

	return balanceFloat.Text('f', digits)
}