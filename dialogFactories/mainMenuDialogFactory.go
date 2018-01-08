package dialogFactories

import (
	"github.com/gameraccoon/telegram-bot-skeleton/dialog"
	"github.com/gameraccoon/telegram-bot-skeleton/dialogFactory"
	"github.com/gameraccoon/telegram-bot-skeleton/processing"
	"github.com/nicksnyder/go-i18n/i18n"
)

type mainMenuItemVariantPrototype struct {
	id string
	textId string
	process func(*processing.ProcessData) bool
}

type mainMenuDialogFactory struct {
	variants []mainMenuItemVariantPrototype
}

func MakeMainMenuDialogFactory() dialogFactory.DialogFactory {
	return &(mainMenuDialogFactory{
		variants: []mainMenuItemVariantPrototype{
			mainMenuItemVariantPrototype{
				id: "wl",
				textId: "wallets",
				process: showWallets,
			},
			mainMenuItemVariantPrototype{
				id: "ex",
				textId: "exchange_rates",
				process: showExchangeRates,
			},
		},
	})
}

func showWallets(data *processing.ProcessData) bool {
	data.SubstitudeDialog(data.Static.MakeDialogFn("wl", data.UserId, data.Trans, data.Static))
	return true
}

func showExchangeRates(data *processing.ProcessData) bool {
	data.SubstitudeMessage(data.Trans("not_supported"))
	data.SendDialog(data.Static.MakeDialogFn("wl", data.UserId, data.Trans, data.Static))
	return true
}

func (factory *mainMenuDialogFactory) createVariants(staticData *processing.StaticProccessStructs, trans i18n.TranslateFunc) (variants []dialog.Variant) {
	variants = make([]dialog.Variant, 0)

	for _, variant := range factory.variants {
		variants = append(variants, dialog.Variant{
			Id:   variant.id,
			Text: trans(variant.textId),
		})
	}
	return
}

func (factory *mainMenuDialogFactory) MakeDialog(userId int64, trans i18n.TranslateFunc, staticData *processing.StaticProccessStructs) *dialog.Dialog {
	return &dialog.Dialog{
		Text:     trans("main_menu_title"),
		Variants: factory.createVariants(staticData, trans),
	}
}

func (factory *mainMenuDialogFactory) ProcessVariant(variantId string, additionalId string, data *processing.ProcessData) bool {
	for _, variant := range factory.variants {
		if variant.id == variantId {
			return variant.process(data)
		}
	}
	return false
}
