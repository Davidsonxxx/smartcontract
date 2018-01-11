package cryptoFunctions

import (
	"math/big"
)

func GetFloatBalance(intValue *big.Int, digits int) *big.Float {
	if intValue == nil {
		return nil
	}
	
	// balanceFloat = balance / (10.0 ** currencyDigits)
	return new(big.Float).Quo(new(big.Float).SetInt(intValue), new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(digits)), big.NewInt(0))))
}

func FormatCurrencyAmount(intValue *big.Int, digits int) string {
	
	var balanceFloat *big.Float = GetFloatBalance(intValue, digits)

	if balanceFloat == nil {
		return "Error"
	}

	return balanceFloat.Text('f', digits)
}
