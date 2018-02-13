package cryptoFunctions

import (
	"gitlab.com/gameraccoon/telegram-accountant-bot/currencies"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"log"
	"math/big"
	"regexp"
)

var bitcoinAddressRegex *regexp.Regexp

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

func init() {
	bitcoinAddressRegex = regexp.MustCompile("^[13][a-zA-Z0-9]{25,34}$")
	if bitcoinAddressRegex == nil {
		log.Fatal("Wrong regexp")
	}
}

func (processor *BitcoinProcessor) GetBalance(address currencies.AddressData) *big.Int {
	resp, err := http.Get("https://chain.api.btc.com/v3/address/" + address.Address)
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

	var parsedResp = new(BitcoinResp)
	err = json.Unmarshal(body, &parsedResp)
	if(err != nil){
		log.Print(string(body[:]))
		log.Print(err)
		return nil
	}

	return big.NewInt(parsedResp.Data.Balance)
}

func (processor *BitcoinProcessor) GetBalanceBunch(addresses []currencies.AddressData) []*big.Int {
	if len(addresses) == 1 {
		return []*big.Int {
			processor.GetBalance(addresses[0]),
		}
	}

	balances := make([]*big.Int, len(addresses))

	resp, err := http.Get("https://chain.api.btc.com/v3/address/" + joinAddresses(addresses))
	if err != nil {
		log.Print(err)
		return balances
	}
	defer resp.Body.Close()

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
		addressesIndexes[address.Address] = i
	}

	for _, data := range parsedResp.Data {
		if i, ok := addressesIndexes[data.Address]; ok {
			balances[i] = big.NewInt(data.Balance)
		}
	}

	return balances
}

func (processor *BitcoinProcessor) GetTransactionsHistory(address currencies.AddressData, limit int) []currencies.TransactionsHistoryItem {
	return make([]currencies.TransactionsHistoryItem, 0)
}

func (processor *BitcoinProcessor) IsAddressValid(address string) bool {
	return isBitcoinAddressValid(address)
}

func isBitcoinAddressValid(address string) bool {
	return bitcoinAddressRegex.MatchString(address)
}
