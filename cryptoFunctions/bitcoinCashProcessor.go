package cryptoFunctions

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"math/big"
)

type BitcoinCashProcessor struct {
}

type BitcoinCashRespData struct {
	SumValueUnspent string `json:"sum_value_unspent"`
}

type BitcoinCashResp struct {
	Data []BitcoinCashRespData `json:"data"`
}

func (processor *BitcoinCashProcessor) GetBalance(address string) *big.Int {
	resp, err := http.Get("https://api.blockchair.com/bitcoin-cash/dashboards/address/" + address)
	defer resp.Body.Close()
	if err != nil {
		log.Print(err)
		return nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
		return nil
	}

	var parsedResp = new(BitcoinCashResp)
	err = json.Unmarshal(body, &parsedResp)
	if(err != nil){
		log.Print(string(body[:]))
		log.Print(err)
		return nil
	}

	if len(parsedResp.Data) > 0 {
		intValue, err := strconv.ParseInt(parsedResp.Data[0].SumValueUnspent, 10, 64)

		if err == nil {
			return big.NewInt(intValue)
		} else {
			return nil
		}
	} else {
		return nil
	}
}

func (processor *BitcoinCashProcessor) GetSumBalance(addresses []string) *big.Int {
	sumBalance := big.NewInt(0)

	for _, walletAddress := range addresses {
		balance := processor.GetBalance(walletAddress)
		if balance != nil {
			sumBalance.Add(sumBalance, balance)
		}
	}

	return sumBalance
}

func (processor *BitcoinCashProcessor) GetToUsdRate() *big.Float {
	return getCurrencyToUsdRate("bitcoin-cash")
}
