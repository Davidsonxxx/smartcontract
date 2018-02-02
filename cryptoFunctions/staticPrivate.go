package cryptoFunctions

import (
	"bytes"
	"encoding/json"
	"gitlab.com/gameraccoon/telegram-accountant-bot/currencies"
	"io/ioutil"
	"net/http"
	"log"
	"strings"
	"math/big"
)

type RateData struct {
	PriceUsd string `json:"price_usd"`
}

// currencyId see here https://coinmarketcap.com/api/
func getCurrencyToUsdRate(currencyId string) *big.Float {
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

func joinAddresses(addresses []currencies.AddressData) string {
	var b bytes.Buffer

	for i, addressData := range addresses {
		if i > 0 {
			b.WriteString(",")
		}
		b.WriteString(addressData.Address)
	}

	return b.String()
}
