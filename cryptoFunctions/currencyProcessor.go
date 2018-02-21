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
	// get history of transactions sorted from new to old
	GetTransactionsHistory(address currencies.AddressData, limit int) []currencies.TransactionsHistoryItem
	// check adress for validness
	IsAddressValid(address string) bool
}
