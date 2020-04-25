package database

import (
	"github.com/gameraccoon/telegram-accountant-bot/currencies"
)

type WalletAddressDbWrapper struct {
	Data currencies.AddressData
	WalletId int64
}
