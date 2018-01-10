package dialogFactories

import (
	"github.com/gameraccoon/telegram-bot-skeleton/dialog"
	"github.com/gameraccoon/telegram-bot-skeleton/dialogFactory"
	"github.com/gameraccoon/telegram-bot-skeleton/processing"
	"github.com/nicksnyder/go-i18n/i18n"
	"gitlab.com/gameraccoon/telegram-accountant-bot/database"
	"strconv"
)

type receiveVariantPrototype struct {
	id string
	textId string
	process func(int64, *processing.ProcessData) bool
	rowId int
}

type receiveDialogFactory struct {
	variants []receiveVariantPrototype
}

func MakeReceiveDialogFactory() dialogFactory.DialogFactory {
	return &(receiveDialogFactory{
		variants: []receiveVariantPrototype{
			receiveVariantPrototype{
				id: "back",
				textId: "back_to_wallet",
				process: backToWallet, // declared in walletSettingsDialogFactory.go
				rowId:1,
			},
		},
	})
}

func (factory *receiveDialogFactory) createText(walletId int64, trans i18n.TranslateFunc, staticData *processing.StaticProccessStructs) string {
	walletAddress := database.GetWalletAddress(staticData.Db, walletId)
	
	return trans("receive_title") + "\n" + walletAddress.Address
}

func (factory *receiveDialogFactory) createVariants(walletId int64, trans i18n.TranslateFunc, staticData *processing.StaticProccessStructs) (variants []dialog.Variant) {
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

func (factory *receiveDialogFactory) MakeDialog(walletId int64, trans i18n.TranslateFunc, staticData *processing.StaticProccessStructs) *dialog.Dialog {
	return &dialog.Dialog{
		Text:     factory.createText(walletId, trans, staticData),
		Variants: factory.createVariants(walletId, trans, staticData),
	}
}

func (factory *receiveDialogFactory) ProcessVariant(variantId string, additionalId string, data *processing.ProcessData) bool {
	walletId, err := strconv.ParseInt(additionalId, 10, 64)

	if err != nil {
		return false
	}

	if !database.IsWalletBelongsToUser(data.Static.Db, data.UserId, walletId) {
		return false
	}

	for _, variant := range factory.variants {
		if variant.id == variantId {
			return variant.process(walletId, data)
		}
	}
	return false
}
