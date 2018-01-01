package chat

import (
	//"gitlab.com/gameraccoon/telegram-accountant-bot/database"
	"gitlab.com/gameraccoon/telegram-accountant-bot/dialog"
)

type Chat interface {
	SendMessage(chatId int64, message string) int64
	SendDialog(chatId int64, dialog *dialog.Dialog) int64
}
