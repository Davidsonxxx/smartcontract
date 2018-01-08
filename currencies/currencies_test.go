package currencies

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetCodeAndName(t *testing.T) {
	assert := require.New(t)

	btcCode := GetCurrencyCode(Bitcoin)
	assert.Equal("BTC", btcCode)

	ethCode := GetCurrencyCode(Ether)
	assert.Equal("ETH", ethCode)

	btcName := GetCurrencyFullName(Bitcoin)
	assert.Equal("Bitcoin", btcName)
}
