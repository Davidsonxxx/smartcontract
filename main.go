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
	"gitlab.com/gameraccoon/telegram-accountant-bot/serverData"
	static "gitlab.com/gameraccoon/telegram-accountant-bot/staticData"
	"io/ioutil"
	"log"
	"strings"
	"sync"
	"time"
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

func loadConfig(path string) (config static.StaticConfiguration, err error) {
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
		Config: config,
		Trans:  translators,
		MakeDialogFn: func(id string, userId int64, trans i18n.TranslateFunc, staticData *processing.StaticProccessStructs) *dialog.Dialog {
			return dialogManager.MakeDialog(id, userId, trans, staticData)
		},
	}

	staticData.Init()

	dbMutex := &sync.Mutex{}

	serverDataManager := serverData.ServerDataManager{}
	serverDataManager.RegisterServerDataInterface(staticData)
	serverDataManager.InitialUpdate(db, dbMutex)

	go updateTimer(staticData, &serverDataManager, config.UpdateIntervalSec, dbMutex)
	updateBot(chat, staticData, dialogManager, dbMutex)
}

func updateTimer(staticData *processing.StaticProccessStructs, serverDataManager *serverData.ServerDataManager, updateIntervalSec int, dbMutex *sync.Mutex) {
	if updateIntervalSec <= 0 {
		log.Fatal("Wrong time interval. Add updateIntervalSec to config")
	}

	for {
		time.Sleep(time.Duration(updateIntervalSec) * time.Second)
		serverDataManager.TimerTick(staticData.Db, dbMutex)
	}
}

func updateBot(chat *telegramChat.TelegramChat, staticData *processing.StaticProccessStructs, dialogManager *dialogManager.DialogManager, dbMutex *sync.Mutex) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := chat.GetBot().GetUpdatesChan(u)

	if err != nil {
		log.Fatal(err.Error())
	}

	processors := makeUserCommandProcessors()

	for update := range updates {
		if update.Message != nil {
			processMessageUpdate(&update, staticData, dialogManager, &processors, dbMutex)
		}
		if update.CallbackQuery != nil {
			processCallbackUpdate(&update, staticData, dialogManager, &processors, dbMutex)
		}
	}
}

func processMessageUpdate(update *tgbotapi.Update, staticData *processing.StaticProccessStructs, dialogManager *dialogManager.DialogManager, processors *ProcessorFuncMap, dbMutex *sync.Mutex) {
	data := processing.ProcessData{
		Static: staticData,
		ChatId: update.Message.Chat.ID,
	}

	userLangCode := strings.ToLower(update.Message.From.LanguageCode)

	message := update.Message.Text

	if strings.HasPrefix(message, "/") {
		commandLen := strings.Index(message, " ")
		if commandLen != -1 {
			data.Command = message[1:commandLen]
			data.Message = message[commandLen+1:]
		} else {
			data.Command = message[1:]
		}

		dbMutex.Lock()
		processCommand(&data, dialogManager, processors, userLangCode)
		dbMutex.Unlock()
	} else {
		data.Message = message

		dbMutex.Lock()
		processPlainMessage(&data, dialogManager, userLangCode)
		dbMutex.Unlock()
	}
}

func processCallbackUpdate(update *tgbotapi.Update, staticData *processing.StaticProccessStructs, dialogManager *dialogManager.DialogManager, processors *ProcessorFuncMap, dbMutex *sync.Mutex) {
	data := processing.ProcessData{
		Static:            staticData,
		ChatId:            int64(update.CallbackQuery.From.ID),
		AnsweredMessageId: int64(update.CallbackQuery.Message.MessageID),
	}

	userLangCode := strings.ToLower(update.CallbackQuery.From.LanguageCode)

	message := update.CallbackQuery.Data

	commandLen := strings.Index(message, " ")
	if commandLen != -1 {
		data.Command = message[1:commandLen]
		data.Message = message[commandLen+1:]
	} else {
		data.Command = message[1:]
	}

	dbMutex.Lock()
	processCommand(&data, dialogManager, processors, userLangCode)
	dbMutex.Unlock()
}
