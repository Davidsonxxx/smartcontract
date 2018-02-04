package dialogFactories

import (
	"bytes"
	"github.com/gameraccoon/telegram-bot-skeleton/dialog"
	"github.com/gameraccoon/telegram-bot-skeleton/dialogFactory"
	"github.com/gameraccoon/telegram-bot-skeleton/processing"
	"github.com/nicksnyder/go-i18n/i18n"
	"gitlab.com/gameraccoon/telegram-accountant-bot/cryptoFunctions"
	"gitlab.com/gameraccoon/telegram-accountant-bot/serverData"
	"gitlab.com/gameraccoon/telegram-accountant-bot/staticFunctions"
	"strconv"
	"time"
)

type historyVariantPrototype struct {
	id string
	textId string
	process func(int64, *processing.ProcessData) bool
	rowId int
}

type historyDialogFactory struct {
	variants []historyVariantPrototype
}

func MakeHistoryDialogFactory() dialogFactory.DialogFactory {
	return &(historyDialogFactory{
		variants: []historyVariantPrototype{
			historyVariantPrototype{
				id: "back",
				textId: "back_to_wallet",
				process: backToWallet, // declared in walletSettingsDialogFactory.go
				rowId:1,
			},
		},
	})
}

func (factory *historyDialogFactory) createText(walletId int64, trans i18n.TranslateFunc, staticData *processing.StaticProccessStructs) string {
	serverData := serverData.GetServerData(staticData)

	if serverData == nil {
		return "Error"
	}

	var textBuffer bytes.Buffer
	textBuffer.WriteString(trans("history_title"))

	walletAddress := staticFunctions.GetDb(staticData).GetWalletAddress(walletId)

	processor := cryptoFunctions.GetProcessor(walletAddress.Currency)

	if processor != nil {
		history := (*processor).GetTransactionsHistory(walletAddress, 0)

		for _, item := range history {
			textBuffer.WriteString("\n")

			textBuffer.WriteString(trans("transaction_time"))
			textBuffer.WriteString(item.Time.Format(time.UnixDate))

			if item.From != "" {
				textBuffer.WriteString(trans("from_addr"))
				if item.From == walletAddress.Address {
					textBuffer.WriteString(trans("me"))
				} else {
					textBuffer.WriteString(item.From)
				}
			}

			if item.To != "" {
				textBuffer.WriteString(trans("to_addr"))
				if item.To == walletAddress.Address {
					textBuffer.WriteString(trans("me"))
				} else {
					textBuffer.WriteString(item.To)
				}
			}

			textBuffer.WriteString(trans("amount"))

			currencySymbol, currencyDecimals := staticFunctions.GetCurrencySymbolAndDecimals(serverData, walletAddress.Currency, walletAddress.ContractAddress)
			amountText := cryptoFunctions.FormatCurrencyAmount(item.Amount, currencyDecimals)
			textBuffer.WriteString(amountText + " " + currencySymbol)
		}
	}

	return textBuffer.String()
}

func (factory *historyDialogFactory) createVariants(walletId int64, trans i18n.TranslateFunc, staticData *processing.StaticProccessStructs) (variants []dialog.Variant) {
	variants = make([]dialog.Variant, 0)

	for _, variant := range factory.variants {
		variants = append(variants, dialog.Variant{
			Id:   variant.id,
			Text: trans(variant.textId),
			AdditionalId: strconv.FormatInt(walletId, 10),
			RowId: variant.rowId,
		})
	}
	return
}

func (factory *historyDialogFactory) MakeDialog(walletId int64, trans i18n.TranslateFunc, staticData *processing.StaticProccessStructs) *dialog.Dialog {
	return &dialog.Dialog{
		Text:     factory.createText(walletId, trans, staticData),
		Variants: factory.createVariants(walletId, trans, staticData),
	}
}

func (factory *historyDialogFactory) ProcessVariant(variantId string, additionalId string, data *processing.ProcessData) bool {
	walletId, err := strconv.ParseInt(additionalId, 10, 64)

	if err != nil {
		return false
	}

	if !staticFunctions.GetDb(data.Static).IsWalletBelongsToUser(data.UserId, walletId) {
		return false
	}

	for _, variant := range factory.variants {
		if variant.id == variantId {
			return variant.process(walletId, data)
		}
	}
	return false
}
