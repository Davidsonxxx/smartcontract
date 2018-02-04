package cryptoFunctions

import (
	"math/big"
	"strings"
)

func GetFloatBalance(intValue *big.Int, digits int) *big.Float {
	if intValue == nil {
		return nil
	}

	// return balance / (10.0 ** currencyDecimals)
	return new(big.Float).Quo(new(big.Float).SetInt(intValue), new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(digits)), big.NewInt(0))))
}

func FormatFloatCurrencyAmount(floatValue *big.Float, digits int) string {
	resultText := floatValue.Text('f', digits)

	if strings.ContainsAny(resultText, ".") {
		return strings.TrimRight(strings.TrimRight(resultText, "0"), ".")
	} else {
		return resultText
	}
}

func FormatCurrencyAmount(intValue *big.Int, digits int) string {

	var floatValue *big.Float = GetFloatBalance(intValue, digits)

	if floatValue == nil {
		return "Error"
	}

	return FormatFloatCurrencyAmount(floatValue, digits)
}
