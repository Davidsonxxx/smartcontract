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
	Code string
	Digits int // how mady decimal digits after zero can it have
}

var currencyStaticDataMap map[Currency]currencyStaticData

func init() {
	currencyStaticDataMap = map[Currency]currencyStaticData {
		Bitcoin : {
			FullName: "Bitcoin",
			Code: "BTC",
			Digits: 8,
		},
		Ether : {
			FullName: "Ethereum",
			Code: "ETH",
			Digits: 18,
		},
		BitcoinCash : {
			FullName: "Bitcoin Cash",
			Code: "BCH",
			Digits: 8,
		},
		BitcoinGold : {
			FullName: "Bitcoin Gold",
			Code: "BTG",
			Digits: 8,
		},
		RippleXrp : {
			FullName: "Ripple",
			Code: "XRP",
			Digits: 6,
		},
		Erc20Token : {
			FullName: "ERC20 Token",
			Code: "",
			Digits: 18,
		},
	}
}

func GetCurrencyCode(currency Currency) string {
	currencyData, ok := currencyStaticDataMap[currency]

	if !ok {
		log.Printf("Unknown currency: %d ", int8(currency))
		return "unknown"
	}

	return currencyData.Code
}

func GetCurrencyFullName(currency Currency) string {
	currencyData, ok := currencyStaticDataMap[currency]

	if !ok {
		log.Printf("Unknown currency: %d ", int8(currency))
		return "unknown"
	}

	return currencyData.FullName
}

func GetCurrencyDigits(currency Currency) int {
	currencyData, ok := currencyStaticDataMap[currency]

	if !ok {
		log.Printf("Unknown currency: %d ", int8(currency))
		return 0
	}

	return currencyData.Digits
}
