package currencies

import (
	"math/big"
)

type BalanceNotify struct {
	NotifyId int64
	WalletId int64
	UserId int64
	OldBalance *big.Int
	NewBalance *big.Int
}
