package dialogFactories

import (
	"gitlab.com/gameraccoon/telegram-accountant-bot/processing"
	"gitlab.com/gameraccoon/telegram-accountant-bot/dialog"
	"gitlab.com/gameraccoon/telegram-accountant-bot/dialogFactory"
	"github.com/nicksnyder/go-i18n/i18n"
	"fmt"
	"strconv"
)

type walletVariantPrototype struct {
	id string
	textId string
	process func(int64, *processing.ProcessData) bool
	// nil if the variant is always active
	isActiveFn func(int64, *processing.StaticProccessStructs) bool
	rowId int
}

type walletDialogFactory struct {
	variants []walletVariantPrototype
}

func MakeWalletDialogFactory() dialogFactory.DialogFactory {
	return &(walletDialogFactory{
		variants: []walletVariantPrototype{
			walletVariantPrototype{
				id: "send",
				textId: "send",
				process: sendFromWallet,
				isActiveFn: isFullWallet,
				rowId:1,
			},
			walletVariantPrototype{
				id: "get",
				textId: "receive",
				process: receiveToWallet,
				rowId:1,
			},
			walletVariantPrototype{
				id: "hist",
				textId: "history",
				process: showHistory,
				rowId:2,
			},
			walletVariantPrototype{
				id: "set",
				textId: "settings",
				process: walletSettings,
				rowId:2,
			},
			walletVariantPrototype{
				id: "back",
				textId: "back_to_list",
				process: backToList,
				rowId:3,
			},
		},
	})
}

func isFullWallet(walletId int64, staticData *processing.StaticProccessStructs) bool {
	return false
}

func sendFromWallet(walletId int64, data *processing.ProcessData) bool {
	data.SubstitudeMessage(data.Trans("not_supported"))
	data.SendDialog(data.Static.MakeDialogFn("wa", data.UserId, data.Trans, data.Static))
	return true
}

func receiveToWallet(walletId int64, data *processing.ProcessData) bool {
	data.SubstitudeDialog(data.Static.MakeDialogFn("rc", walletId, data.Trans, data.Static))
	return true
}

func showHistory(walletId int64, data *processing.ProcessData) bool {
	data.SubstitudeMessage(data.Trans("not_supported"))
	data.SendDialog(data.Static.MakeDialogFn("wa", data.UserId, data.Trans, data.Static))
	return true
}

func walletSettings(walletId int64, data *processing.ProcessData) bool {
	data.SubstitudeDialog(data.Static.MakeDialogFn("ws", walletId, data.Trans, data.Static))
	return true
}

func backToList(walletId int64, data *processing.ProcessData) bool {
	data.SubstitudeDialog(data.Static.MakeDialogFn("wl", data.UserId, data.Trans, data.Static))
	return true
}

func (factory *walletDialogFactory) getListItemDialogText(walletId int64, trans i18n.TranslateFunc, staticData *processing.StaticProccessStructs) string {
	return fmt.Sprintf("<b>%s</b>", staticData.Db.GetWalletName(walletId))
}

func (factory *walletDialogFactory) createVariants(walletId int64, trans i18n.TranslateFunc, staticData *processing.StaticProccessStructs) (variants []dialog.Variant) {
	variants = make([]dialog.Variant, 0)

	for _, variant := range factory.variants {
		if variant.isActiveFn == nil || variant.isActiveFn(walletId, staticData) {
			variants = append(variants, dialog.Variant{
				Id:   variant.id,
				Text: trans(variant.textId),
				AdditionalId: strconv.FormatInt(walletId, 10),
				RowId: variant.rowId,
			})
		}
	}
	return
}

func (factory *walletDialogFactory) MakeDialog(walletId int64, trans i18n.TranslateFunc, staticData *processing.StaticProccessStructs) *dialog.Dialog {
	return &dialog.Dialog{
		Text:     factory.getListItemDialogText(walletId, trans, staticData),
		Variants: factory.createVariants(walletId, trans, staticData),
	}
}

func (factory *walletDialogFactory) ProcessVariant(variantId string, additionalId string, data *processing.ProcessData) bool {
	walletId, err := strconv.ParseInt(additionalId, 10, 64)

	if err != nil {
		return false
	}

	if !data.Static.Db.IsWalletBelongsToUser(data.UserId, walletId) {
		return false
	}

	for _, variant := range factory.variants {
		if variant.id == variantId {
			return variant.process(walletId, data)
		}
	}
	return false
}
