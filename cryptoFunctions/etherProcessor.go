package cryptoFunctions

import (
	"encoding/json"
	"gitlab.com/gameraccoon/telegram-accountant-bot/currencies"
	"io/ioutil"
	"net/http"
	"log"
	"math/big"
)

const etherscanApiKey string = "KBT56RI9SUTF2GR1TNN41W48FUQ4YAK3GK"

type EtherProcessor struct {
}

type EtherRespData struct {
	Account string `json:"account"`
	Balance string `json:"balance"`
}

type EtherResp struct {
	Result string `json:"result"`
}

type EtherMultiResp struct {
	Result []EtherRespData `json:"result"`
}

func (processor *EtherProcessor) GetBalance(address currencies.AddressData) *big.Int {
	resp, err := http.Get("http://api.etherscan.io/api?module=account&action=balance&address=" + address.Address + "&tag=latest&apikey=" + etherscanApiKey)
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

	var parsedResp = new(EtherResp)
	err = json.Unmarshal(body, &parsedResp)
	if(err != nil){
		log.Print(string(body[:]))
		log.Print(err)
		return nil
	}

	intValue := new(big.Int)
	_, ok := intValue.SetString(parsedResp.Result, 10)

	if ok {
		return intValue
	} else {
		log.Print(string(body[:]))
		log.Print("Int parse problem")
		return nil
	}
}

func (processor *EtherProcessor) GetBalanceBunch(addresses []currencies.AddressData) []*big.Int {
	if len(addresses) == 1 {
		return []*big.Int {
			processor.GetBalance(addresses[0]),
		}
	}

	balances := make([]*big.Int, len(addresses))

	resp, err := http.Get("http://api.etherscan.io/api?module=account&action=balancemulti&address=" + joinAddresses(addresses) + "&tag=latest&apikey=" + etherscanApiKey)
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

	var parsedResp = new(EtherMultiResp)
	err = json.Unmarshal(body, &parsedResp)
	if(err != nil){
		log.Print(string(body[:]))
		log.Print(err)
		return balances
	}

	// I'm not sure if it's more time efficient
	addressesIndexes := map[string]int{}
	for i, address := range addresses {
		addressesIndexes[address.Address] = i
	}

	for _, data := range parsedResp.Result {
		if intValue, ok := new(big.Int).SetString(data.Balance, 10); ok {
			if i, ok := addressesIndexes[data.Account]; ok {
				balances[i] = intValue
			}
		}
	}

	return balances
}

func (processor *EtherProcessor) GetToUsdRate() *big.Float {
	return getCurrencyToUsdRate("ethereum")
}
