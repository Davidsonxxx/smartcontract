package currencies

import (
	"math/big"
)

type Erc20TokenBalanceData struct {
	Name string
	Symbol string
	Balance *big.Int
	Decimals int64
}