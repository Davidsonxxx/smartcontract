package currencies

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetSymbolAndName(t *testing.T) {
	assert := require.New(t)

	btcSymbol := GetCurrencySymbol(Bitcoin)
	assert.Equal("BTC", btcSymbol)

	ethSymbol := GetCurrencySymbol(Ether)
	assert.Equal("ETH", ethSymbol)

	btcName := GetCurrencyFullName(Bitcoin)
	assert.Equal("Bitcoin", btcName)
}
