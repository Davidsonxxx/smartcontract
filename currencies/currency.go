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
	Decimals int // how mady decimal digits after zero can it have
}

var currencyStaticDataMap map[Currency]currencyStaticData

func init() {
	currencyStaticDataMap = map[Currency]currencyStaticData {
		Bitcoin : {
			FullName: "Bitcoin",
			Symbol: "BTC",
			Decimals: 8,
		},
		Ether : {
			FullName: "Ethereum",
			Symbol: "ETH",
			Decimals: 18,
		},
		BitcoinCash : {
			FullName: "Bitcoin Cash",
			Symbol: "BCH",
			Decimals: 8,
		},
		BitcoinGold : {
			FullName: "Bitcoin Gold",
			Symbol: "BTG",
			Decimals: 8,
		},
		RippleXrp : {
			FullName: "Ripple",
			Symbol: "XRP",
			Decimals: 6,
		},
		Erc20Token : {
			FullName: "ERC20 Token",
			Symbol: "",
			Decimals: 18,
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

func GetCurrencyDecimals(currency Currency) int {
	currencyData, ok := currencyStaticDataMap[currency]

	if !ok {
		log.Printf("Unknown currency: %d ", int8(currency))
		return 0
	}

	return currencyData.Decimals
}
