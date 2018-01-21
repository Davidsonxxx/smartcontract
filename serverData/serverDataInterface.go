package serverData

import (
	"gitlab.com/gameraccoon/telegram-accountant-bot/currencies"
	"math/big"
)

type ServerDataInterface interface {
	GetBalance(address currencies.AddressData) *big.Int
	GetRateToUsd(currency currencies.Currency) *big.Float
}
