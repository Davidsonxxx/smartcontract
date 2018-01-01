package dialogFactories

import (
	"gitlab.com/gameraccoon/telegram-accountant-bot/processing"
	"gitlab.com/gameraccoon/telegram-accountant-bot/dialog"
	"gitlab.com/gameraccoon/telegram-accountant-bot/dialogFactory"
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
				id: "wallets",
				textId: "wallets",
				process: showWallets,
			},
			mainMenuItemVariantPrototype{
				id: "exch",
				textId: "exchange_rates",
				process: showExchangeRates,
			},
		},
	})
}

func showWallets(data *processing.ProcessData) bool {
	data.Static.Chat.SendMessage(data.ChatId, "test message 1")
	return true
}

func showExchangeRates(data *processing.ProcessData) bool {
	data.Static.Chat.SendMessage(data.ChatId, "test message 2")
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

func (factory *mainMenuDialogFactory) MakeDialog(itemId int64, trans i18n.TranslateFunc, staticData *processing.StaticProccessStructs) *dialog.Dialog {
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
