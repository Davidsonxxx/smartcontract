package currencies

import (
	"math/big"
	"time"
)

type TransactionsHistoryItem struct {
	From string
	To string
	Amount *big.Int
	Time time.Time
}
