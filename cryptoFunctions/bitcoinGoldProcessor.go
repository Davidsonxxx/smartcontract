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
		return big.NewInt(-1)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
		return big.NewInt(-1)
	}

	floatValue, err := strconv.ParseFloat(string(body[:]), 64)

	if err == nil {
		return big.NewInt(int64(floatValue * 1.0E8))
	} else {
		log.Print(err)
		return big.NewInt(-1)
	}
}

func (processor *BitcoinGoldProcessor) GetSumBalance(addresses []string) *big.Int {
	sumBalance := big.NewInt(0)

	for _, walletAddress := range addresses {
		sumBalance.Add(sumBalance, processor.GetBalance(walletAddress))
	}

	return sumBalance
}
