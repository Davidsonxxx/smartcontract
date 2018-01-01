package processing

import (
	"github.com/nicksnyder/go-i18n/i18n"
)

type ProcessData struct {
	Static  *StaticProccessStructs
	Command string // first part of command without slash(/)
	Message string // parameters of command or plain message
	ChatId  int64
	UserId  int64
	Trans   i18n.TranslateFunc
}
