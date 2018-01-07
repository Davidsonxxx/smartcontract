package processing

import (
	"gitlab.com/gameraccoon/telegram-accountant-bot/chat"
	"gitlab.com/gameraccoon/telegram-accountant-bot/database"
	"gitlab.com/gameraccoon/telegram-accountant-bot/dialog"
	"github.com/nicksnyder/go-i18n/i18n"
	"log"
	"time"
)

type LanguageData struct {
	Key string
	Name string
}

type StaticConfiguration struct {
	AvailableLanguages []LanguageData
	DefaultLanguage string
	ExtendedLog bool
}

type AwaitingTextProcessorData struct {
	ProcessorId string
	AdditionalId int64
}

type UserState struct {
	awaitingTextProcessor *AwaitingTextProcessorData
	currentPage int // temporary data for lists handling
	// temporary data for wallets creation
	newWalletName string
}

type StaticProccessStructs struct {
	Chat chat.Chat
	Db *database.Database
	Timers map[int64]time.Time
	Config *StaticConfiguration
	Trans map[string]i18n.TranslateFunc
	MakeDialogFn func(string, int64, i18n.TranslateFunc, *StaticProccessStructs)*dialog.Dialog
	userStates map[int64]UserState
}

func (staticData *StaticProccessStructs) Init() {
	staticData.userStates = make(map[int64]UserState)
}

func (staticData *StaticProccessStructs) SetUserStateTextProcessor(userId int64, processor *AwaitingTextProcessorData) {
	state := staticData.userStates[userId]
	state.awaitingTextProcessor = processor
	staticData.userStates[userId] = state
}

func (staticData *StaticProccessStructs) GetUserStateTextProcessor(userId int64) *AwaitingTextProcessorData {
	if state, ok := staticData.userStates[userId]; ok {
		return state.awaitingTextProcessor
	} else {
		return nil
	}
}

func (staticData *StaticProccessStructs) SetUserStateCurrentPage(userId int64, page int) {
	state := staticData.userStates[userId]
	state.currentPage = page
	staticData.userStates[userId] = state
}

func (staticData *StaticProccessStructs) GetUserStateCurrentPage(userId int64) int {
	if state, ok := staticData.userStates[userId]; ok {
		return state.currentPage
	} else {
		return 0
	}
}

func (staticData *StaticProccessStructs) SetUserStateNewWalletName(userId int64, walletName string) {
	state := staticData.userStates[userId]
	state.newWalletName = walletName
	staticData.userStates[userId] = state
}

func (staticData *StaticProccessStructs) GetUserStateNewWalletName(userId int64) string {
	if state, ok := staticData.userStates[userId]; ok {
		return state.newWalletName
	} else {
		return "Empty Name"
	}
}

func (staticData *StaticProccessStructs) FindTransFunction(userId int64) i18n.TranslateFunc {
	// ToDo: cache user's lang
	lang := staticData.Db.GetUserLanguage(userId)

	// replace empty language to default one (some clients don't send user's language)
	if len(lang) <= 0 {
		log.Printf("User %d has empty language. Setting to default.", userId)
		lang = staticData.Config.DefaultLanguage
		staticData.Db.SetUserLanguage(userId, lang)
	}

	if foundTrans, ok := staticData.Trans[lang]; ok {
		return foundTrans
	}

	// unknown language, use default instead
	if foundTrans, ok := staticData.Trans[staticData.Config.DefaultLanguage]; ok {
		log.Printf("User %d has unknown language (%s). Setting to default.", userId, lang)
		lang = staticData.Config.DefaultLanguage
		staticData.Db.SetUserLanguage(userId, lang)
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
