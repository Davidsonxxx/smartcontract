package currencies

type CurrencyProcessor interface {
	GetBalance(address string) int
}
