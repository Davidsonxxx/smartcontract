package cryptoFunctions

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type BitcoinCashProcessor struct {
}

type BitcoinCashRespData struct {
	SumValueUnspent string `json:"sum_value_unspent"`
}

type BitcoinCashResp struct {
	Data []BitcoinCashRespData `json:"data"`
}

func (processor *BitcoinCashProcessor) GetBalance(address string) int64 {
	resp, err := http.Get("https://api.blockchair.com/bitcoin-cash/dashboards/address/" + address)
	defer resp.Body.Close()
	if err != nil {
		log.Print(err)
		return -1
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
		return -1
	}

	var parsedResp = new(BitcoinCashResp)
	err = json.Unmarshal(body, &parsedResp)
	if(err != nil){
		log.Print(string(body[:]))
		log.Print(err)
		return -1
	}

	if len(parsedResp.Data) > 0 {
		intValue, err := strconv.ParseInt(parsedResp.Data[0].SumValueUnspent, 10, 64)

		if err == nil {
			return intValue
		} else {
			return 0
		}
	} else {
		return 0
	}
}

func (processor *BitcoinCashProcessor) GetSumBalance(addresses []string) int64 {
	var sumBalance int64 = 0

	for _, walletAddress := range addresses {
		sumBalance = sumBalance + processor.GetBalance(walletAddress)
	}

	return sumBalance
}
