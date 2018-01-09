package cryptoFunctions

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"log"
)

type BitcoinProcessor struct {
}

type BitcoinRespData struct {
	Balance int64 `json:"balance"`
}

type BitcoinResp struct {
	Data BitcoinRespData `json:"data"`
}

func (processor *BitcoinProcessor) GetBalance(address string) int64 {
	resp, err := http.Get("https://chain.api.btc.com/v3/address/" + address)
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

	var parsedResp = new(BitcoinResp)
	err = json.Unmarshal(body, &parsedResp)
	if(err != nil){
		log.Print(err)
		return -1
	}

	return parsedResp.Data.Balance
}
