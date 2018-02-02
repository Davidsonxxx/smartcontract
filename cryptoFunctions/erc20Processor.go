package cryptoFunctions

import (
	"encoding/json"
	"gitlab.com/gameraccoon/telegram-accountant-bot/currencies"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
)

type Erc20Processor struct {
}

type Erc20Resp struct {
	Name string `json:"name"`
	Symbol string `json:"symbol"`
	Balance string `json:"balance"`
	Decimals int64 `json:"decimals"`
}

func (processor *Erc20Processor) GetBalance(address currencies.AddressData) *big.Int {
	resp, err := http.Get("https://api.tokenbalance.com/token/" + address.ContractId + "/" + address.Address)
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

	var parsedResp = new(Erc20Resp)
	err = json.Unmarshal(body, &parsedResp)
	if(err != nil){
		log.Print(string(body[:]))
		log.Print(err)
		return nil
	}

	floatValue, _, err := new(big.Float).Parse(parsedResp.Balance, 10)
	decimals := big.NewInt(parsedResp.Decimals)
	decimalsMultiplier := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), decimals, big.NewInt(0)))

	if err == nil {
		intValue, _ := new(big.Float).Mul(floatValue, decimalsMultiplier).Int(nil)
		return intValue
	} else {
		return nil
	}
}

func (processor *Erc20Processor) GetBalanceBunch(addresses []currencies.AddressData) []*big.Int {
	balances := make([]*big.Int, len(addresses))

	for i, walletAddress := range addresses {
		balances[i] = processor.GetBalance(walletAddress)
	}

	return balances
}

func (processor *Erc20Processor) GetTokenData(contractId string) *currencies.Erc20TokenData {
	resp, err := http.Get("https://api.tokenbalance.com/token/" + contractId + "/0x0")
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

	var parsedResp = new(Erc20Resp)
	err = json.Unmarshal(body, &parsedResp)
	if(err != nil){
		log.Print(string(body[:]))
		log.Print(err)
		return nil
	}

	tokenData := currencies.Erc20TokenData{
		Name: parsedResp.Name,
		Symbol: parsedResp.Symbol,
		Decimals: parsedResp.Decimals,
	}

	return &tokenData
}

func (processor *Erc20Processor) GetToUsdRate() *big.Float {
	return nil
}
