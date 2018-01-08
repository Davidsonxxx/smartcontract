package dialogFactories

import (
	"github.com/gameraccoon/telegram-bot-skeleton/dialog"
	"github.com/gameraccoon/telegram-bot-skeleton/dialogFactory"
	"github.com/gameraccoon/telegram-bot-skeleton/processing"
	"github.com/nicksnyder/go-i18n/i18n"
)

type createWalletItemVariantPrototype struct {
	id string
	textId string
	process func(*processing.ProcessData) bool
}

type createWalletDialogFactory struct {
	variants []createWalletItemVariantPrototype
}

func MakeCreateWalletDialogFactory() dialogFactory.DialogFactory {
	return &(createWalletDialogFactory{
		variants: []createWalletItemVariantPrototype{
			createWalletItemVariantPrototype{
				id: "wo",
				textId: "watch_only",
				process: createWatchOnlyWallet,
			},
			createWalletItemVariantPrototype{
				id: "fl",
				textId: "full_wallet",
				process: createFullWallet,
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

func createFullWallet(data *processing.ProcessData) bool {
	data.SubstitudeMessage(data.Trans("not_supported"))
	data.SendDialog(data.Static.MakeDialogFn("wl", data.UserId, data.Trans, data.Static))
	return true
}

func (factory *createWalletDialogFactory) createVariants(staticData *processing.StaticProccessStructs, trans i18n.TranslateFunc) (variants []dialog.Variant) {
	variants = make([]dialog.Variant, 0)

	for _, variant := range factory.variants {
		variants = append(variants, dialog.Variant{
			Id:   variant.id,
			Text: trans(variant.textId),
		})
	}
	return
}

func (factory *createWalletDialogFactory) MakeDialog(userId int64, trans i18n.TranslateFunc, staticData *processing.StaticProccessStructs) *dialog.Dialog {
	return &dialog.Dialog{
		Text:     trans("choose_wallet_type"),
		Variants: factory.createVariants(staticData, trans),
	}
}

func (factory *createWalletDialogFactory) ProcessVariant(variantId string, additionalId string, data *processing.ProcessData) bool {
	for _, variant := range factory.variants {
		if variant.id == variantId {
			return variant.process(data)
		}
	}
	return false
}
