package cryptoFunctions

import (
	"math/big"
	"net/http"
	"io/ioutil"
	"log"
	"encoding/json"
	"strings"
)

type RateData struct {
	PriceUsd string `json:"price_usd"`
}

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

// currencyId see here https://coinmarketcap.com/api/
func GetCurrencyToUsdRate(currencyId string) *big.Float {
	resp, err := http.Get("https://api.coinmarketcap.com/v1/ticker/" + currencyId + "/")
	if err != nil {
		log.Print(err)
		return nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
		return nil
	}

	jsonStr := string(body[:])

	// remove trailing brackets
	jsonStr = strings.Replace(jsonStr, "[", "", -1)
	jsonStr = strings.Replace(jsonStr, "]", "", -1)

	var parsedResp = new(RateData)
	err = json.Unmarshal([]byte(jsonStr), &parsedResp)
	if(err != nil){
		log.Print(string(body[:]))
		log.Print(err)
		return nil
	}

	rate, _, err := new(big.Float).Parse(parsedResp.PriceUsd, 10)

	if err == nil {
		return rate
	} else {
		log.Print(string(body[:]))
		log.Print(err)
		return nil
	}
}
