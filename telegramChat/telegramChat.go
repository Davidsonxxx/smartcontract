package telegramChat

import (
	"bytes"
	"fmt"
	"gitlab.com/gameraccoon/telegram-accountant-bot/dialog"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type TelegramChat struct {
	bot *tgbotapi.BotAPI
}

func MakeTelegramChat(apiToken string) (bot *TelegramChat, outErr error) {
	newBot, err := tgbotapi.NewBotAPI(apiToken)
	if err != nil {
		outErr = err
		return
	}

	bot = &TelegramChat{
		bot: newBot,
	}

	return
}

func (telegramChat *TelegramChat) GetBot() *tgbotapi.BotAPI {
	return telegramChat.bot
}

func (telegramChat *TelegramChat) SetDebugModeEnabled(isEnabled bool) {
	telegramChat.bot.Debug = isEnabled
}

func (telegramChat *TelegramChat) GetBotUsername() string {
	return telegramChat.bot.Self.UserName
}

func (telegramChat *TelegramChat) SendMessage(chatId int64, message string) (messageId int64) {
	msg := tgbotapi.NewMessage(chatId, message)
	msg.ParseMode = "HTML"
	telegramChat.bot.Send(msg)
	return
}

func getCommand(dialogId string, variantId string, additionalId string) string {
	if additionalId == "" {
		return fmt.Sprintf("/%s_%s", dialogId, variantId)
	} else {
		return fmt.Sprintf("/%s_%s_%s", dialogId, variantId, additionalId)
	}
}

func (telegramChat *TelegramChat) SendDialog(chatId int64, dialog *dialog.Dialog) (messageId int64) {
	var buffer bytes.Buffer

	buffer.WriteString(dialog.Text)

	markup := tgbotapi.NewInlineKeyboardMarkup()

	currentRow := []tgbotapi.InlineKeyboardButton{}
	currentRowId := 0
	for _, variant := range dialog.Variants {
			if currentRowId != variant.RowId {
				if len(currentRow) > 0 {
					markup.InlineKeyboard = append(markup.InlineKeyboard, currentRow)
				}
				currentRow = []tgbotapi.InlineKeyboardButton{}
				currentRowId = variant.RowId
			}
		currentRow = append(currentRow, tgbotapi.NewInlineKeyboardButtonData(
					variant.Text,
					getCommand(dialog.Id, variant.Id, variant.AdditionalId),
		))
	}
	markup.InlineKeyboard = append(markup.InlineKeyboard, currentRow)

	msg := tgbotapi.NewMessage(chatId, buffer.String())
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = &markup
	telegramChat.bot.Send(msg)
	return
}
