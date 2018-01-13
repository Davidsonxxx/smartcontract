package cryptoFunctions

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"math/big"
)

type BitcoinGoldProcessor struct {
}

func (processor *BitcoinGoldProcessor) GetBalance(address string) *big.Int {
	resp, err := http.Get("http://btgexp.com/ext/getbalance/" + address)
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

	floatValue, err := strconv.ParseFloat(string(body[:]), 64)

	if err == nil {
		return big.NewInt(int64(floatValue * 1.0E8))
	} else {
		log.Print(err)
		return nil
	}
}

func (processor *BitcoinGoldProcessor) GetSumBalance(addresses []string) *big.Int {
	sumBalance := big.NewInt(0)

	for _, walletAddress := range addresses {
		balance := processor.GetBalance(walletAddress)
		if balance != nil {
			sumBalance.Add(sumBalance, balance)
		}
	}

	return sumBalance
}

func (processor *BitcoinGoldProcessor) GetBalanceBunch(addresses []string) []*big.Int {
	balances := make([]*big.Int, len(addresses))

	for i, walletAddress := range addresses {
		balances[i] = processor.GetBalance(walletAddress)
	}

	return balances
}

func (processor *BitcoinGoldProcessor) GetToUsdRate() *big.Float {
	return getCurrencyToUsdRate("bitcoin-gold")
}
