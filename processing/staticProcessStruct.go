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
	AdditionalId string
}

type UserState struct {
	awaitingTextProcessor *AwaitingTextProcessorData
	currentPage int
	lastMessages []int64
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

func (staticData *StaticProccessStructs) SetUserStateTextProcessor(userId int64, proessor *AwaitingTextProcessorData) {
	state := staticData.userStates[userId]
	state.awaitingTextProcessor = proessor
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

func (staticData *StaticProccessStructs) FindTransFunction(userId int64) i18n.TranslateFunc {
	// ToDo: cache user's lang
	lang := staticData.Db.GetUserLanguage(userId)

	if len(lang) > 0 {
		foundTrans, ok := staticData.Trans[lang]
		if ok {
			return foundTrans
		} else {
			log.Printf("Translator didn't found: %s", lang)
		}
	} else {
		log.Printf("User %d has empty language", userId)
	}

	// fall to default translation
	for lang, trans := range staticData.Trans {
		log.Printf("Found default translator: %s", lang)
		return trans
	}

	// there are no translators
	log.Fatal("No available translations")
	// we will probably crash but there is nothing else we can do
	translator, _ := i18n.Tfunc(staticData.Config.DefaultLanguage)
	return translator
}
