package currencies

import (
	"math/big"
)

type BalanceNotify struct {
	// filled from DB
	NotifyId int64
	WalletId int64
	UserId int64
	OldBalance *big.Int
	// filled in processing
	IsInitialChange bool // if true, we don't need to show this notification
	WalletAddress AddressData
	NewBalance *big.Int
}
