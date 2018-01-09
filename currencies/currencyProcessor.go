package currencies

type CurrencyProcessor interface {
	GetBalance(address string) int64
	GetSumBalance(addresses []string) int64
}
