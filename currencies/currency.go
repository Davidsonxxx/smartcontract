package currencies

import (
	"log"
)

type Currency int8

const (
	// don't use iota to make it more explicit
	// don't change already assigned numbers (only add new ones)
	Bitcoin Currency = 0
	Ether Currency = 1
	BitcoinCash Currency = 2
	BitcoinGold Currency = 3
	RippleXrp Currency = 4
	Erc20Token Currency = 5
)

type currencyStaticData struct {
	FullName string
	Symbol string
	Decimals int // how many decimal digits after zero can it have
	PriceId string
	// feature flags
	IsHistoryEnabled bool
}

var currencyStaticDataMap map[Currency]currencyStaticData

func init() {
	currencyStaticDataMap = map[Currency]currencyStaticData {
		Bitcoin : {
			FullName: "Bitcoin",
			Symbol: "BTC",
			Decimals: 8,
			PriceId: "bitcoin",
			IsHistoryEnabled: false,
		},
		Ether : {
			FullName: "Ethereum",
			Symbol: "ETH",
			Decimals: 18,
			PriceId: "ethereum",
			IsHistoryEnabled: true,
		},
		BitcoinCash : {
			FullName: "Bitcoin Cash",
			Symbol: "BCH",
			Decimals: 8,
			PriceId: "bitcoin-cash",
			IsHistoryEnabled: false,
		},
		BitcoinGold : {
			FullName: "Bitcoin Gold",
			Symbol: "BTG",
			Decimals: 8,
			PriceId: "bitcoin-gold",
			IsHistoryEnabled: false,
		},
		RippleXrp : {
			FullName: "Ripple",
			Symbol: "XRP",
			Decimals: 6,
			PriceId: "ripple",
			IsHistoryEnabled: false,
		},
		Erc20Token : {
			FullName: "ERC20 Token",
			Symbol: "",
			Decimals: 18,
			PriceId: "",
			IsHistoryEnabled: false,
		},
	}
}

func GetCurrencySymbol(currency Currency) string {
	currencyData, ok := currencyStaticDataMap[currency]

	if !ok {
		log.Printf("Unknown currency: %d ", int8(currency))
		return "unknown"
	}

	return currencyData.Symbol
}

func GetCurrencyFullName(currency Currency) string {
	currencyData, ok := currencyStaticDataMap[currency]

	if !ok {
		log.Printf("Unknown currency: %d ", int8(currency))
		return "unknown"
	}

	return currencyData.FullName
}

func GetCurrencyPriceId(currency Currency) string {
	currencyData, ok := currencyStaticDataMap[currency]

	if !ok {
		log.Printf("Unknown currency: %d ", int8(currency))
		return "unknown"
	}

	return currencyData.PriceId
}

func GetCurrencyDecimals(currency Currency) int {
	currencyData, ok := currencyStaticDataMap[currency]

	if !ok {
		log.Printf("Unknown currency: %d ", int8(currency))
		return 0
	}

	return currencyData.Decimals
}

func IsHistoryEnabled(currency Currency) bool {
	currencyData, ok := currencyStaticDataMap[currency]

	if !ok {
		log.Printf("Unknown currency: %d ", int8(currency))
		return false
	}

	return currencyData.IsHistoryEnabled
}

func GetAllCurrencies() (availableCurrencies []Currency) {
	for currency, _ := range currencyStaticDataMap {
		availableCurrencies = append(availableCurrencies, currency)
	}
	return
}