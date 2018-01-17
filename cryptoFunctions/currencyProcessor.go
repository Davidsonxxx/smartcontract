package cryptoFunctions

import (
	"math/big"
)

type CurrencyProcessor interface {
	// get account balance
	GetBalance(address string) *big.Int
	// get multiple accounts balance
	GetBalanceBunch(addresses []string) []*big.Int
	// Deprecated. Use GetBunchBalance instead
	GetSumBalance(addresses []string) *big.Int
	// get this currency to USD rate
	GetToUsdRate() *big.Float
}
