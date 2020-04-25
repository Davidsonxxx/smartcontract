package dialogFactories

import (
	"bytes"
	"fmt"
	"github.com/gameraccoon/telegram-bot-skeleton/dialog"
	"github.com/gameraccoon/telegram-bot-skeleton/dialogFactory"
	"github.com/gameraccoon/telegram-bot-skeleton/processing"
	"github.com/nicksnyder/go-i18n/i18n"
	"github.com/gameraccoon/telegram-accountant-bot/cryptoFunctions"
	"github.com/gameraccoon/telegram-accountant-bot/currencies"
	"github.com/gameraccoon/telegram-accountant-bot/serverData"
	"github.com/gameraccoon/telegram-accountant-bot/staticFunctions"
	"strconv"
	"strings"
)

const (
	maxHistoryRecords = 25
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

	db := staticFunctions.GetDb(staticData)

	userTimezone := db.GetUserTimezone(db.GetWalletOwner(walletId))

	var textBuffer bytes.Buffer

	walletAddress := db.GetWalletAddress(walletId)

	processor := cryptoFunctions.GetProcessor(walletAddress.Currency)

	if processor != nil {
		if currencies.IsHistoryEnabled(walletAddress.Currency) {
			history := (*processor).GetTransactionsHistory(walletAddress, 25)

			if (len(history) == maxHistoryRecords) {
				textBuffer.WriteString(fmt.Sprintf(trans("history_cut_title"), len(history)))
			} else {
				textBuffer.WriteString(fmt.Sprintf(trans("history_title")))
			}
			
			for i := len(history)-1; i >= 0; i-- {
				item := history[i]
				textBuffer.WriteString("\n\n")

				textBuffer.WriteString(staticFunctions.FormatTimestamp(item.Time, userTimezone))

				currencySymbol, currencyDecimals := staticFunctions.GetCurrencySymbolAndDecimals(serverData, walletAddress.Currency, walletAddress.ContractAddress)
				amountText := cryptoFunctions.FormatCurrencyAmount(item.Amount, currencyDecimals)

				if strings.EqualFold(item.From, walletAddress.Address) {
					textBuffer.WriteString(fmt.Sprintf(trans("sent_format"), amountText, currencySymbol, item.To))
				} else if strings.EqualFold(item.To, walletAddress.Address) {
					textBuffer.WriteString(fmt.Sprintf(trans("recieved_format"), amountText, currencySymbol, item.From))
				}
			}
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
