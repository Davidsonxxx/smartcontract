package cryptoFunctions

import (
	"gitlab.com/gameraccoon/telegram-accountant-bot/currencies"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"math/big"
)

type BitcoinGoldProcessor struct {
}

func (processor *BitcoinGoldProcessor) GetBalance(address currencies.AddressData) *big.Int {
	resp, err := http.Get("http://btgexp.com/ext/getbalance/" + address.Address)
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

	floatValue, err := strconv.ParseFloat(string(body[:]), 64)

	if err == nil {
		return big.NewInt(int64(floatValue * 1.0E8))
	} else {
		log.Print(err)
		return nil
	}
}

func (processor *BitcoinGoldProcessor) GetBalanceBunch(addresses []currencies.AddressData) []*big.Int {
	balances := make([]*big.Int, len(addresses))

	for i, walletAddress := range addresses {
		balances[i] = processor.GetBalance(walletAddress)
	}

	return balances
}

func (processor *BitcoinGoldProcessor) GetTransactionsHistory(address currencies.AddressData, limit int) (isAvailable bool, history []currencies.TransactionsHistoryItem) {
	return false, make([]currencies.TransactionsHistoryItem, 0)
}

func (processor *BitcoinGoldProcessor) IsAddressValid(address string) bool {
	return isBitcoinAddressValid(address)
}
