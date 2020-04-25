package serverData

import (
	"github.com/gameraccoon/telegram-accountant-bot/currencies"
	"math/big"
)

type ServerDataInterface interface {
	GetBalance(address currencies.AddressData) *big.Int
	GetRateToUsd(priceId string) *big.Float
	GetErc20TokenData(contractAddress string) *currencies.Erc20TokenData
}
