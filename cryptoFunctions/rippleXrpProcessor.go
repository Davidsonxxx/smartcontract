package cryptoFunctions

import (
	"gitlab.com/gameraccoon/telegram-accountant-bot/currencies"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"math/big"
)

type RippleXrpProcessor struct {
}

type RippleXrpRespData struct {
	Value string `json:"value"`
}

type RippleXrpResp struct {
	Balances []RippleXrpRespData `json:"balances"`
}

func (processor *RippleXrpProcessor) GetBalance(address currencies.AddressData) *big.Int {
	resp, err := http.Get("https://data.ripple.com/v2/accounts/" + address.Address + "/balances?currency=XRP&limit=1")
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

	var parsedResp = new(RippleXrpResp)
	err = json.Unmarshal(body, &parsedResp)
	if(err != nil){
		log.Print(string(body[:]))
		log.Print(err)
		return nil
	}

	if len(parsedResp.Balances) > 0 {
		floatValue, _, err := new(big.Float).Parse(parsedResp.Balances[0].Value, 10)

		if err == nil {
			intValue, _ := new(big.Float).Mul(floatValue, new(big.Float).SetFloat64(1000000)).Int(nil)
			return intValue
		} else {
			return nil
		}
	} else {
		return nil
	}
}

func (processor *RippleXrpProcessor) GetBalanceBunch(addresses []currencies.AddressData) []*big.Int {
	balances := make([]*big.Int, len(addresses))

	for i, walletAddress := range addresses {
		balances[i] = processor.GetBalance(walletAddress)
	}

	return balances
}

func (processor *RippleXrpProcessor) GetToUsdRate() *big.Float {
	return getCurrencyToUsdRate("ripple")
}

func (processor *RippleXrpProcessor) GetTransactionsHistory(address currencies.AddressData, limit int) []currencies.TransactionsHistoryItem {
	return make([]currencies.TransactionsHistoryItem, 0)
}
