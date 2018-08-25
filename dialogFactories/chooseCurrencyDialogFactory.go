package dialogFactories

import (
	"github.com/gameraccoon/telegram-bot-skeleton/dialog"
	"github.com/gameraccoon/telegram-bot-skeleton/dialogFactory"
	"github.com/gameraccoon/telegram-bot-skeleton/processing"
	"gitlab.com/gameraccoon/telegram-accountant-bot/currencies"
	"github.com/nicksnyder/go-i18n/i18n"
)

type chooseCurrencyItemVariantPrototype struct {
	id string
	currencyId currencies.Currency
	rowId int
}

type chooseCurrencyDialogFactory struct {
	variants []chooseCurrencyItemVariantPrototype
}

func MakeChooseCurrencyDialogFactory() dialogFactory.DialogFactory {
	return &(chooseCurrencyDialogFactory{
		variants: []chooseCurrencyItemVariantPrototype{
			chooseCurrencyItemVariantPrototype{
				id: "btc",
				currencyId: currencies.Bitcoin,
				rowId: 1,
			},
			chooseCurrencyItemVariantPrototype{
				id: "eth",
				currencyId: currencies.Ether,
				rowId: 1,
			},
			chooseCurrencyItemVariantPrototype{
				id: "bch",
				currencyId: currencies.BitcoinCash,
				rowId: 2,
			},
			chooseCurrencyItemVariantPrototype{
				id: "btg",
				currencyId: currencies.BitcoinGold,
				rowId: 2,
			},
			chooseCurrencyItemVariantPrototype{
				id: "xrp",
				currencyId: currencies.RippleXrp,
				rowId: 3,
			},
			chooseCurrencyItemVariantPrototype{
				id: "erc20",
				currencyId: currencies.Erc20Token,
				rowId: 3,
			},
		},
	})
}

func processWalletType(data *processing.ProcessData, variantPrototype *chooseCurrencyItemVariantPrototype) bool {
	data.Static.CleanUserStateValues(data.UserId)
	data.Static.SetUserStateValue(data.UserId, "walletCurrency", variantPrototype.currencyId)
	data.SubstitudeMessage(data.Trans("send_wallet_name"))
	data.Static.SetUserStateTextProcessor(data.UserId, &processing.AwaitingTextProcessorData{
		ProcessorId: "newWalletName",
	})
	return true
}

func (factory *chooseCurrencyDialogFactory) createVariants(staticData *processing.StaticProccessStructs, trans i18n.TranslateFunc) (variants []dialog.Variant) {
	variants = make([]dialog.Variant, 0)

	for _, variant := range factory.variants {
		variants = append(variants, dialog.Variant{
			Id:    variant.id,
			Text:  currencies.GetCurrencyFullName(variant.currencyId),
			RowId: variant.rowId,
		})
	}
	return
}

func (factory *chooseCurrencyDialogFactory) MakeDialog(userId int64, trans i18n.TranslateFunc, staticData *processing.StaticProccessStructs) *dialog.Dialog {
	return &dialog.Dialog{
		Text:     trans("choose_wallet_type"),
		Variants: factory.createVariants(staticData, trans),
	}
}

func (factory *chooseCurrencyDialogFactory) ProcessVariant(variantId string, additionalId string, data *processing.ProcessData) bool {
	for _, variant := range factory.variants {
		if variant.id == variantId {
			return processWalletType(data, &variant)
		}
	}
	return false
}
