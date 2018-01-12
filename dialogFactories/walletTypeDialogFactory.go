package dialogFactories

import (
	"github.com/gameraccoon/telegram-bot-skeleton/dialog"
	"github.com/gameraccoon/telegram-bot-skeleton/dialogFactory"
	"github.com/gameraccoon/telegram-bot-skeleton/processing"
	"github.com/nicksnyder/go-i18n/i18n"
)

type walletTypeItemVariantPrototype struct {
	id string
	textId string
	process func(*processing.ProcessData) bool
}

type walletTypeDialogFactory struct {
	variants []walletTypeItemVariantPrototype
}

func MakeWalletTypeDialogFactory() dialogFactory.DialogFactory {
	return &(walletTypeDialogFactory{
		variants: []walletTypeItemVariantPrototype{
			walletTypeItemVariantPrototype{
				id: "wo",
				textId: "watch_only",
				process: createWatchOnlyWallet,
			},
		},
	})
}

func createWatchOnlyWallet(data *processing.ProcessData) bool {
	data.SubstitudeMessage(data.Trans("send_wallet_name"))
	data.Static.SetUserStateTextProcessor(data.UserId, &processing.AwaitingTextProcessorData{
		ProcessorId: "newWatchOnlyWalletName",
	})
	return true
}

func (factory *walletTypeDialogFactory) createVariants(staticData *processing.StaticProccessStructs, trans i18n.TranslateFunc) (variants []dialog.Variant) {
	variants = make([]dialog.Variant, 0)

	for _, variant := range factory.variants {
		variants = append(variants, dialog.Variant{
			Id:   variant.id,
			Text: trans(variant.textId),
		})
	}
	return
}

func (factory *walletTypeDialogFactory) MakeDialog(userId int64, trans i18n.TranslateFunc, staticData *processing.StaticProccessStructs) *dialog.Dialog {
	return &dialog.Dialog{
		Text:     trans("choose_wallet_type"),
		Variants: factory.createVariants(staticData, trans),
	}
}

func (factory *walletTypeDialogFactory) ProcessVariant(variantId string, additionalId string, data *processing.ProcessData) bool {
	for _, variant := range factory.variants {
		if variant.id == variantId {
			return variant.process(data)
		}
	}
	return false
}
