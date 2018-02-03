package cryptoFunctions

import (
	"gitlab.com/gameraccoon/telegram-accountant-bot/currencies"
	"math/big"
)

type CurrencyProcessor interface {
	// get account balance
	GetBalance(address currencies.AddressData) *big.Int
	// get multiple accounts balance
	GetBalanceBunch(addresses []currencies.AddressData) []*big.Int
	// get this currency to USD rate
	GetToUsdRate() *big.Float
}
