package main

import (
	"encoding/json"
	"github.com/gameraccoon/telegram-bot-skeleton/dialog"
	"github.com/gameraccoon/telegram-bot-skeleton/dialogManager"
	"github.com/gameraccoon/telegram-bot-skeleton/processing"
	"github.com/gameraccoon/telegram-bot-skeleton/telegramChat"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nicksnyder/go-i18n/i18n"
	"gitlab.com/gameraccoon/telegram-accountant-bot/database"
	"gitlab.com/gameraccoon/telegram-accountant-bot/dialogFactories"
	"io/ioutil"
	"log"
	"strings"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func getFileStringContent(filePath string) (content string, err error) {
	fileContent, err := ioutil.ReadFile(filePath)
	if err == nil {
		content = strings.TrimSpace(string(fileContent))
	}
	return
}

func getApiToken() (token string, err error) {
	return getFileStringContent("./telegramApiToken.txt")
}

func loadConfig(path string) (config processing.StaticConfiguration, err error) {
	jsonString, err := getFileStringContent(path)
	if err == nil {
		dec := json.NewDecoder(strings.NewReader(jsonString))
		err = dec.Decode(&config)
	}
	return
}

func main() {
	apiToken, err := getApiToken()
	if err != nil {
		log.Fatal(err.Error())
	}

	config, err := loadConfig("./config.json")
	if err != nil {
		log.Fatal(err.Error())
	}

	translators := make(map[string]i18n.TranslateFunc)

	for _, lang := range config.AvailableLanguages {
		i18n.MustLoadTranslationFile("./data/strings/" + lang.Key + ".all.json")

		trans, err := i18n.Tfunc(lang.Key)
		if err != nil {
			log.Fatal(err.Error())
		}
		translators[lang.Key] = trans
	}

	if len(translators) <= 0 {
		log.Fatal("Need at least one language available")
	}

	if _, ok := translators[config.DefaultLanguage]; !ok {
		log.Fatal("Default language should be in the list of available languages")
	}

	db, err := database.Init("./accounts-data.db")
	defer db.Disconnect()

	if err != nil {
		log.Fatal("Can't connect database")
	}

	database.UpdateVersion(db)

	chat, err := telegramChat.MakeTelegramChat(apiToken)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("Authorized on account %s", chat.GetBotUsername())

	chat.SetDebugModeEnabled(config.ExtendedLog)

	dialogManager := &(dialogManager.DialogManager{})
	dialogManager.RegisterDialogFactory("lc", dialogFactories.MakeLanguageSelectDialogFactory())
	dialogManager.RegisterDialogFactory("mn", dialogFactories.MakeMainMenuDialogFactory())
	dialogManager.RegisterDialogFactory("wt", dialogFactories.MakeWalletTypeDialogFactory())
	dialogManager.RegisterDialogFactory("wl", dialogFactories.MakeWalletsListDialogFactory())
	dialogManager.RegisterDialogFactory("wa", dialogFactories.MakeWalletDialogFactory())
	dialogManager.RegisterDialogFactory("ws", dialogFactories.MakeWalletSettingsDialogFactory())
	dialogManager.RegisterDialogFactory("rc", dialogFactories.MakeReceiveDialogFactory())
	dialogManager.RegisterDialogFactory("de", dialogFactories.MakeDeleteConfirmationDialogFactory())
	dialogManager.RegisterDialogFactory("hi", dialogFactories.MakeHistoryDialogFactory())
	dialogManager.RegisterDialogFactory("cc", dialogFactories.MakeChooseCurrencyDialogFactory())
	dialogManager.RegisterTextInputProcessorManager(dialogFactories.GetTextInputProcessorManager())

	staticData := &processing.StaticProccessStructs{
		Chat:   chat,
		Db:     db,
		Config: &config,
		Trans:  translators,
		MakeDialogFn: func(id string, userId int64, trans i18n.TranslateFunc, staticData *processing.StaticProccessStructs) *dialog.Dialog {
			return dialogManager.MakeDialog(id, userId, trans, staticData)
		},
	}

	staticData.Init()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := chat.GetBot().GetUpdatesChan(u)

	if err != nil {
		log.Fatal(err.Error())
	}

	processors := makeUserCommandProcessors()

	for update := range updates {
		if update.Message != nil {
			processMessageUpdate(&update, staticData, dialogManager, &processors)
		}
		if update.CallbackQuery != nil {
			processCallbackUpdate(&update, staticData, dialogManager, &processors)
		}
	}
}
