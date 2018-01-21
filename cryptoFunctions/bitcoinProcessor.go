package cryptoFunctions

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"log"
	"strings"
	"math/big"
)

type BitcoinProcessor struct {
}

type BitcoinRespData struct {
	Address string `json:"address"`
	Balance int64 `json:"balance"`
}

type BitcoinResp struct {
	Data BitcoinRespData `json:"data"`
}

type BitcoinMultiResp struct {
	Data []BitcoinRespData `json:"data"`
}

func (processor *BitcoinProcessor) GetBalance(address string) *big.Int {
	resp, err := http.Get("https://chain.api.btc.com/v3/address/" + address)
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

	var parsedResp = new(BitcoinResp)
	err = json.Unmarshal(body, &parsedResp)
	if(err != nil){
		log.Print(string(body[:]))
		log.Print(err)
		return nil
	}

	return big.NewInt(parsedResp.Data.Balance)
}

func (processor *BitcoinProcessor) GetBalanceBunch(addresses []string) []*big.Int {
	if len(addresses) == 1 {
		return []*big.Int {
			processor.GetBalance(addresses[0]),
		}
	}

	balances := make([]*big.Int, len(addresses))

	resp, err := http.Get("https://chain.api.btc.com/v3/address/" + strings.Join(addresses, ","))
	defer resp.Body.Close()
	if err != nil {
		log.Print(err)
		return balances
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
		return balances
	}

	var parsedResp = new(BitcoinMultiResp)
	err = json.Unmarshal(body, &parsedResp)
	if(err != nil){
		log.Print(string(body[:]))
		log.Print(err)
		return balances
	}

	// I'm not sure if it's more time efficient
	addressesIndexes := map[string]int{}
	for i, address := range addresses {
		addressesIndexes[address] = i
	}

	for _, data := range parsedResp.Data {
		if i, ok := addressesIndexes[data.Address]; ok {
			balances[i] = big.NewInt(data.Balance)
		}
	}

	return balances
}

func (processor *BitcoinProcessor) GetToUsdRate() *big.Float {
	return getCurrencyToUsdRate("bitcoin")
}
