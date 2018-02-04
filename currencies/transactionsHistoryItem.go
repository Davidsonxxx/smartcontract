package currencies

import (
	"math/big"
)

type TransactionsHistoryItem struct {
	From string
	To string
	Amount *big.Int
}
