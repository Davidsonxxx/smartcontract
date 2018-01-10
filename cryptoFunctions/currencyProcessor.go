package cryptoFunctions

import (
	"math/big"
)

type CurrencyProcessor interface {
	GetBalance(address string) *big.Int
	GetSumBalance(addresses []string) *big.Int
}
