package wallettypes

type WalletType int8

const (
	// don't use iota to make it more explicit
	WatchOnly WalletType = 0
)
