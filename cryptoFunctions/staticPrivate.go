package cryptoFunctions

import (
	"bytes"
	"github.com/gameraccoon/telegram-accountant-bot/currencies"
)

func joinAddresses(addresses []currencies.AddressData) string {
	var b bytes.Buffer

	for i, addressData := range addresses {
		if i > 0 {
			b.WriteString(",")
		}
		b.WriteString(addressData.Address)
	}

	return b.String()
}
