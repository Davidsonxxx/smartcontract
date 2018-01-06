package dialogFactories

import (
	"gitlab.com/gameraccoon/telegram-accountant-bot/processing"
	"gitlab.com/gameraccoon/telegram-accountant-bot/dialog"
	"gitlab.com/gameraccoon/telegram-accountant-bot/dialogFactory"
	"github.com/nicksnyder/go-i18n/i18n"
)

type languageSelectVariantPrototype struct {
	id string
	text string
	process func(*processing.ProcessData) bool
}

type languageSelectDialogFactory struct {
}

func MakeLanguageSelectDialogFactory() dialogFactory.DialogFactory {
	return &(languageSelectDialogFactory{})
}

func applyNewLanguage(data *processing.ProcessData, newLang string) bool {
	data.Static.Db.SetUserLanguage(data.UserId, newLang)
	data.Trans = data.Static.FindTransFunction(data.UserId)
	data.SendDialog(data.Static.MakeDialogFn("mn", data.UserId, data.Trans, data.Static))
	return true
}

func (factory *languageSelectDialogFactory) createVariants(staticData *processing.StaticProccessStructs) (variants []dialog.Variant) {
	variants = make([]dialog.Variant, 0)

	itemId := 0
	itemsInRow := 2

	for _, lang := range staticData.Config.AvailableLanguages {
		variants = append(variants, dialog.Variant{
			Id:   lang.Key,
			Text: lang.Name,
			RowId: itemId / itemsInRow + 1,
		})
		itemId = itemId + 1
	}
	return
}

func (factory *languageSelectDialogFactory) MakeDialog(itemId int64, trans i18n.TranslateFunc, staticData *processing.StaticProccessStructs) *dialog.Dialog {
	return &dialog.Dialog{
		Text:     trans("select_language"),
		Variants: factory.createVariants(staticData),
	}
}

func (factory *languageSelectDialogFactory) ProcessVariant(variantId string, additionalId string, data *processing.ProcessData) bool {
	return applyNewLanguage(data, variantId)
}
