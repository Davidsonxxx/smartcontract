package cryptoFunctions

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type BitcoinGoldProcessor struct {
}

func (processor *BitcoinGoldProcessor) GetBalance(address string) int64 {
	resp, err := http.Get("http://btgexp.com/ext/getbalance/" + address)
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

	floatValue, err := strconv.ParseFloat(string(body[:]), 64)

	if err == nil {
		return int64(floatValue * 1.0E8)
	} else {
		log.Print(err)
		return 0
	}
}

func (processor *BitcoinGoldProcessor) GetSumBalance(addresses []string) int64 {
	var sumBalance int64 = 0

	for _, walletAddress := range addresses {
		sumBalance = sumBalance + processor.GetBalance(walletAddress)
	}

	return sumBalance
}
