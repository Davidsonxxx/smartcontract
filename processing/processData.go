package processing

import (
	"github.com/nicksnyder/go-i18n/i18n"
	"gitlab.com/gameraccoon/telegram-accountant-bot/dialog"
)

type ProcessData struct {
	Static            *StaticProccessStructs
	Command           string // first part of command without slash(/)
	Message           string // parameters of command or plain message
	ChatId            int64
	UserId            int64
	Trans             i18n.TranslateFunc
	AnsweredMessageId int64
}

func (data *ProcessData) SendMessage(message string) int64 {
	return data.Static.Chat.SendMessage(data.ChatId, message, data.AnsweredMessageId)
}

func (data *ProcessData) SendDialog(dialog *dialog.Dialog) int64 {
	return data.Static.Chat.SendDialog(data.ChatId, dialog, data.AnsweredMessageId)
}
