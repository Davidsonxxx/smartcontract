package cryptoFunctions

import (
	"encoding/json"
	"fmt"
	"gitlab.com/gameraccoon/telegram-accountant-bot/currencies"
	"io/ioutil"
	"net/http"
	"log"
	"math/big"
	"strconv"
	"time"
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

type EtherHistoryRespItem struct {
	From string `json:"from"`
	To string `json:"to"`
	Value string `json:"value"`
	ContractAddress string `json:"contractAddress"`
	TimeStamp string `json:"timeStamp"`
}

type EtherHistoryResp struct {
	Result []EtherHistoryRespItem `json:"result"`
}

func (processor *EtherProcessor) GetBalance(address currencies.AddressData) *big.Int {
	resp, err := http.Get("http://api.etherscan.io/api?module=account&action=balance&address=" + address.Address + "&tag=latest&apikey=" + etherscanApiKey)
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

func (processor *EtherProcessor) GetTransactionsHistory(address currencies.AddressData, limit int) []currencies.TransactionsHistoryItem {
	var requestText string
	if limit > 0 {
		requestText = fmt.Sprintf(
			"http://api.etherscan.io/api?module=account&action=txlist&address=%s&startblock=0&endblock=99999999&sort=asc&apikey=%s",
			address.Address,
			limit,
			etherscanApiKey,
		)
	} else {
		requestText = fmt.Sprintf(
			"https://api.etherscan.io/api?module=account&action=txlist&address=%s&startblock=0&endblock=99999999&page=1&offset=%d&sort=asc&apikey=%s",
			address.Address,
			etherscanApiKey,
		)
	}

	resp, err := http.Get(requestText)
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

	var parsedResp = new(EtherHistoryResp)
	err = json.Unmarshal(body, &parsedResp)
	if(err != nil){
		log.Print(string(body[:]))
		log.Print(err)
		return nil
	}

	history := make([]currencies.TransactionsHistoryItem, 0, len(parsedResp.Result))

	for _, historyItem := range parsedResp.Result {
		amount, ok := new(big.Int).SetString(historyItem.Value, 10)

		if !ok {
			amount = big.NewInt(0)
			log.Printf("Wrong amount value: %s", historyItem.Value)
		}

		var to string
		if historyItem.To != "" {
			to = historyItem.To
		} else {
			to = historyItem.ContractAddress
		}

		intTime, err := strconv.ParseInt(historyItem.TimeStamp, 10, 64)
		if err != nil {
			log.Print(err.Error())
			intTime = int64(0)
		}

		history = append(history, currencies.TransactionsHistoryItem {
				From: historyItem.From,
				To: to,
				Amount: amount,
				Time: time.Unix(intTime, int64(0)),
			})
	}

	return history
}

func (processor *EtherProcessor) IsAddressValid(address string) bool {
	return true
}
