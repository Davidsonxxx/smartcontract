package dialogFactories

import (
	"github.com/gameraccoon/telegram-bot-skeleton/dialog"
	"github.com/gameraccoon/telegram-bot-skeleton/dialogFactory"
	"github.com/gameraccoon/telegram-bot-skeleton/processing"
	"gitlab.com/gameraccoon/telegram-accountant-bot/database"
	"github.com/nicksnyder/go-i18n/i18n"
	"strconv"
)

type walletSettingsVariantPrototype struct {
	id string
	textId string
	process func(int64, *processing.ProcessData) bool
	rowId int
}

type walletSettingsDialogFactory struct {
	variants []walletSettingsVariantPrototype
}

func MakeWalletSettingsDialogFactory() dialogFactory.DialogFactory {
	return &(walletSettingsDialogFactory{
		variants: []walletSettingsVariantPrototype{
			walletSettingsVariantPrototype{
				id: "ren",
				textId: "rename",
				process: renameWallet,
				rowId:1,
			},
			walletSettingsVariantPrototype{
				id: "del",
				textId: "delete",
				process: deleteWallet,
				rowId:1,
			},
			walletSettingsVariantPrototype{
				id: "back",
				textId: "back_to_wallet",
				process: backToWallet,
				rowId:2,
			},
		},
	})
}

func renameWallet(walletId int64, data *processing.ProcessData) bool {
	data.Static.SetUserStateTextProcessor(data.UserId, &processing.AwaitingTextProcessorData{
		ProcessorId: "renamingWallet",
		AdditionalId: walletId,
	})
	data.SubstitudeMessage(data.Trans("rename_wallet_request"))
	return true
}

func deleteWallet(walletId int64, data *processing.ProcessData) bool {
	data.SubstitudeDialog(data.Static.MakeDialogFn("de", walletId, data.Trans, data.Static))
	return true
}

func backToWallet(walletId int64, data *processing.ProcessData) bool {
	data.SubstitudeDialog(data.Static.MakeDialogFn("wa", walletId, data.Trans, data.Static))
	return true
}

func (factory *walletSettingsDialogFactory) createVariants(walletId int64, trans i18n.TranslateFunc, staticData *processing.StaticProccessStructs) (variants []dialog.Variant) {
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

func (factory *walletSettingsDialogFactory) MakeDialog(walletId int64, trans i18n.TranslateFunc, staticData *processing.StaticProccessStructs) *dialog.Dialog {
	return &dialog.Dialog{
		Text:     trans("settings_title"),
		Variants: factory.createVariants(walletId, trans, staticData),
	}
}

func (factory *walletSettingsDialogFactory) ProcessVariant(variantId string, additionalId string, data *processing.ProcessData) bool {
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
