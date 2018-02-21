package cryptoFunctions

import (
	"gitlab.com/gameraccoon/telegram-accountant-bot/currencies"
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

func (processor *BitcoinCashProcessor) GetBalance(address currencies.AddressData) *big.Int {
	resp, err := http.Get("https://api.blockchair.com/bitcoin-cash/dashboards/address/" + address.Address)
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

func (processor *BitcoinCashProcessor) GetBalanceBunch(addresses []currencies.AddressData) []*big.Int {
	balances := make([]*big.Int, len(addresses))

	for i, walletAddress := range addresses {
		balances[i] = processor.GetBalance(walletAddress)
	}

	return balances
}

func (processor *BitcoinCashProcessor) GetTransactionsHistory(address currencies.AddressData, limit int) []currencies.TransactionsHistoryItem {
	return make([]currencies.TransactionsHistoryItem, 0)
}

func (processor *BitcoinCashProcessor) IsAddressValid(address string) bool {
	return isBitcoinAddressValid(address)
}
