package dialogFactories

import (
	"github.com/gameraccoon/telegram-bot-skeleton/dialog"
	"github.com/gameraccoon/telegram-bot-skeleton/dialogFactory"
	"github.com/gameraccoon/telegram-bot-skeleton/processing"
	"github.com/nicksnyder/go-i18n/i18n"
	"gitlab.com/gameraccoon/telegram-accountant-bot/database"
	"strconv"
)

type deleteConfirmationVariantPrototype struct {
	id string
	textId string
	process func(int64, *processing.ProcessData) bool
	isActiveFn func(int64, *processing.StaticProccessStructs) bool
	rowId int
}

type deleteConfirmationDialogFactory struct {
	variants []deleteConfirmationVariantPrototype
}

func MakeDeleteConfirmationDialogFactory() dialogFactory.DialogFactory {
	return &(deleteConfirmationDialogFactory{
		variants: []deleteConfirmationVariantPrototype{
			deleteConfirmationVariantPrototype{
				id: "del1",
				textId: "accept_del_watch_only",
				process: deleteWalletFinally,
				isActiveFn: isWatchOnlyWallet,
				rowId:1,
			},
			deleteConfirmationVariantPrototype{
				id: "del2",
				textId: "accept_del_full",
				process: deleteWalletFinally,
				isActiveFn: isFullWallet, // declared in walletDialogFactory.go
				rowId:1,
			},
			deleteConfirmationVariantPrototype{
				id: "back",
				textId: "reject_del",
				process: backToWallet, // implemented in walletSettingsDialogFactory.go
				rowId:2,
			},
		},
	})
}

func isWatchOnlyWallet(walletId int64, staticData *processing.StaticProccessStructs) bool {
	return true
}

func deleteWalletFinally(walletId int64, data *processing.ProcessData) bool {
	database.DeleteWallet(data.Static.Db, walletId)
	data.SubstitudeMessage(data.Trans("deleted_success"))
	data.SendDialog(data.Static.MakeDialogFn("wl", data.UserId, data.Trans, data.Static))
	return true
}

func (factory *deleteConfirmationDialogFactory) createVariants(walletId int64, trans i18n.TranslateFunc, staticData *processing.StaticProccessStructs) (variants []dialog.Variant) {
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

func (factory *deleteConfirmationDialogFactory) MakeDialog(walletId int64, trans i18n.TranslateFunc, staticData *processing.StaticProccessStructs) *dialog.Dialog {
	var text string

	if isWatchOnlyWallet(walletId, staticData) {
		text = trans("title_watch_only_deleting")
	} else {
		text = trans("title_full_deleting")
	}

	return &dialog.Dialog{
		Text:     text,
		Variants: factory.createVariants(walletId, trans, staticData),
	}
}

func (factory *deleteConfirmationDialogFactory) ProcessVariant(variantId string, additionalId string, data *processing.ProcessData) bool {
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
