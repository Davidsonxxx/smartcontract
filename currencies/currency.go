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
)

type currencyStaticData struct {
	FullName string
	Code string
}

var currencyStaticDataMap map[Currency]currencyStaticData

func init() {
	currencyStaticDataMap = map[Currency]currencyStaticData {
		Bitcoin : {
			FullName: "Bitcoin",
			Code: "BTC",
		},
		Ether : {
			FullName: "Ether",
			Code: "ETH",
		},
		BitcoinCash : {
			FullName: "Bitcoin Cash",
			Code: "BCH",
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
