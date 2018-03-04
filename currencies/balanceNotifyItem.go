package currencies

import (
	"math/big"
)

type BalanceNotify struct {
	NotifyId int64
	WalletId int64
	LastBalance *big.Int
}