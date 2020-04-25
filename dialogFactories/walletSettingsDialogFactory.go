package dialogFactories

import (
	"github.com/gameraccoon/telegram-bot-skeleton/dialog"
	"github.com/gameraccoon/telegram-bot-skeleton/dialogFactory"
	"github.com/gameraccoon/telegram-bot-skeleton/processing"
	"github.com/nicksnyder/go-i18n/i18n"
	"github.com/gameraccoon/telegram-accountant-bot/staticFunctions"
	"github.com/gameraccoon/telegram-accountant-bot/currencies"
	"strconv"
)

type walletSettingsData struct {
	walletId int64
	staticData *processing.StaticProccessStructs
	isNotificationsEnabled bool
}

type walletSettingsVariantPrototype struct {
	id string
	textId string
	process func(int64, *processing.ProcessData) bool
	rowId int
	isActiveFn func(*walletSettingsData) bool
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
				id: "prc",
				textId: "change_price_id",
				process: changePriceId,
				rowId:2,
				isActiveFn: isErc20TokenWallet,
			},
			walletSettingsVariantPrototype{
				id: "onntfy",
				textId: "enable_notify",
				process: enableBalanceNotifications,
				rowId:3,
				isActiveFn: isNotificationsDisabled,
			},
			walletSettingsVariantPrototype{
				id: "offntfy",
				textId: "disable_notify",
				process: disableBalanceNotifications,
				rowId:3,
				isActiveFn: isNotificationsEnabled,
			},
			walletSettingsVariantPrototype{
				id: "back",
				textId: "back_to_wallet",
				process: backToWallet,
				rowId:4,
			},
		},
	})
}

func isErc20TokenWallet(settingsData *walletSettingsData) bool {
	walletAddress := staticFunctions.GetDb(settingsData.staticData).GetWalletAddress(settingsData.walletId)
	return walletAddress.Currency == currencies.Erc20Token
}

func isNotificationsEnabled(settingsData *walletSettingsData) bool {
	return settingsData.isNotificationsEnabled
}

func isNotificationsDisabled(settingsData *walletSettingsData) bool {
	return !settingsData.isNotificationsEnabled
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

func enableBalanceNotifications(walletId int64, data *processing.ProcessData) bool {
	staticFunctions.GetDb(data.Static).EnableBalanceNotifies(walletId)

	data.SubstitudeDialog(data.Static.MakeDialogFn("ws", walletId, data.Trans, data.Static))
	return true
}

func disableBalanceNotifications(walletId int64, data *processing.ProcessData) bool {
	staticFunctions.GetDb(data.Static).DisableBalanceNotifies(walletId)

	data.SubstitudeDialog(data.Static.MakeDialogFn("ws", walletId, data.Trans, data.Static))
	return true
}

func changePriceId(walletId int64, data *processing.ProcessData) bool {
	data.Static.SetUserStateTextProcessor(data.UserId, &processing.AwaitingTextProcessorData{
		ProcessorId: "setWalletPriceId",
		AdditionalId: walletId,
	})
	data.SubstitudeMessage(data.Trans("send_price_id"))
	return true
}

func backToWallet(walletId int64, data *processing.ProcessData) bool {
	data.SubstitudeDialog(data.Static.MakeDialogFn("wa", walletId, data.Trans, data.Static))
	return true
}

func (factory *walletSettingsDialogFactory) createVariants(settingsData *walletSettingsData, trans i18n.TranslateFunc) (variants []dialog.Variant) {
	variants = make([]dialog.Variant, 0)

	for _, variant := range factory.variants {
		if variant.isActiveFn == nil || variant.isActiveFn(settingsData) {
			variants = append(variants, dialog.Variant{
				Id:   variant.id,
				Text: trans(variant.textId),
				AdditionalId: strconv.FormatInt(settingsData.walletId, 10),
				RowId: variant.rowId,
			})
		}
	}
	return
}

func (factory *walletSettingsDialogFactory) MakeDialog(walletId int64, trans i18n.TranslateFunc, staticData *processing.StaticProccessStructs) *dialog.Dialog {
	isNotificationsEnabled := staticFunctions.GetDb(staticData).IsBalanceNotifiesEnabled(walletId)

	settingsData := walletSettingsData {
		walletId: walletId,
		staticData: staticData,
		isNotificationsEnabled: isNotificationsEnabled,
	}

	var notificationsText string
	if isNotificationsEnabled {
		notificationsText = trans("balance_notify_enabled")
	} else {
		notificationsText = trans("balance_notify_disabled")
	}

	return &dialog.Dialog{
		Text:     trans("settings_title") + notificationsText,
		Variants: factory.createVariants(&settingsData, trans),
	}
}

func (factory *walletSettingsDialogFactory) ProcessVariant(variantId string, additionalId string, data *processing.ProcessData) bool {
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
