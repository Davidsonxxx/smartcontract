package cryptoFunctions

import (
	"math/big"
)

type CurrencyProcessor interface {
	// get account balance
	GetBalance(address string) *big.Int
	// get sum balance of 
	GetSumBalance(addresses []string) *big.Int
	// get this currency to USD rate
	GetToUsdRate() *big.Float
}
