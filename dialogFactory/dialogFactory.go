package dialogFactory

import (
	"gitlab.com/gameraccoon/telegram-accountant-bot/dialog"
	"gitlab.com/gameraccoon/telegram-accountant-bot/processing"
	"github.com/nicksnyder/go-i18n/i18n"
)

type DialogFactory interface {
	MakeDialog(id int64, trans i18n.TranslateFunc, staticData *processing.StaticProccessStructs) *dialog.Dialog
	ProcessVariant(variantId string, additionalId string, data *processing.ProcessData) bool
}
