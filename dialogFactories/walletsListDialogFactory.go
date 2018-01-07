package dialogFactories

import (
	"gitlab.com/gameraccoon/telegram-accountant-bot/processing"
	"gitlab.com/gameraccoon/telegram-accountant-bot/dialog"
	"gitlab.com/gameraccoon/telegram-accountant-bot/dialogFactory"
	"github.com/nicksnyder/go-i18n/i18n"
	"strconv"
)

type walletsListDialogVariantPrototype struct {
	id string
	additionalIdFn func(*walletsListDialogCache) string
	textId string
	textFn func(*walletsListDialogCache) string
	// nil if the variant is always active
	isActiveFn func(*walletsListDialogCache) bool
	process func(string, *processing.ProcessData) bool
}

type cachedItem struct {
	id int64
	text string
}

type walletsListDialogCache struct {
	cachedItems []cachedItem
	currentPage int
	pagesCount int
	countOnPage int
}

type walletsListDialogFactory struct {
	textId string
	variants []walletsListDialogVariantPrototype
}

func MakeWalletsListDialogFactory() dialogFactory.DialogFactory {
	return &(walletsListDialogFactory{
		textId: "choose_wallet",
		variants: []walletsListDialogVariantPrototype{
			walletsListDialogVariantPrototype{
				id: "add",
				textId: "add_wallet_btn",
				isActiveFn: isTheFirstPage,
				process: addWallet,
			},
			walletsListDialogVariantPrototype{
				id: "it1",
				additionalIdFn: getFirstItemId,
				textFn: getFirstItemText,
				isActiveFn: isFirstElementVisible,
				process: openWallet,
			},
			walletsListDialogVariantPrototype{
				id: "it2",
				additionalIdFn: getSecondItemId,
				textFn: getSecondItemText,
				isActiveFn: isSecondElementVisible,
				process: openWallet,
			},
			walletsListDialogVariantPrototype{
				id: "it3",
				additionalIdFn: getThirdItemId,
				textFn: getThirdItemText,
				isActiveFn: isThirdElementVisible,
				process: openWallet,
			},
			walletsListDialogVariantPrototype{
				id: "it4",
				additionalIdFn: getFourthItemId,
				textFn: getFourthItemText,
				isActiveFn: isFourthElementVisible,
				process: openWallet,
			},
			walletsListDialogVariantPrototype{
				id: "it5",
				additionalIdFn: getFifthItemId,
				textFn: getFifthItemText,
				isActiveFn: isFifthElementVisible,
				process: openWallet,
			},
			walletsListDialogVariantPrototype{
				id: "back",
				textId: "back_btn",
				isActiveFn: isNotTheFirstPage,
				process: moveBack,
			},
			walletsListDialogVariantPrototype{
				id: "fwd",
				textId: "fwd_btn",
				isActiveFn: isNotTheLastPage,
				process: moveForward,
			},
		},
	})
}

func isTheFirstPage(cahce *walletsListDialogCache) bool {
	return cahce.currentPage == 0
}

func isNotTheFirstPage(cahce *walletsListDialogCache) bool {
	return cahce.currentPage > 0
}

func isNotTheLastPage(cahce *walletsListDialogCache) bool {
	return cahce.currentPage + 1 < cahce.pagesCount
}

func isFirstElementVisible(cahce *walletsListDialogCache) bool {
	return cahce.countOnPage > 0
}

func isSecondElementVisible(cahce *walletsListDialogCache) bool {
	return cahce.countOnPage > 1
}

func isThirdElementVisible(cahce *walletsListDialogCache) bool {
	return cahce.countOnPage > 2
}

func isFourthElementVisible(cahce *walletsListDialogCache) bool {
	return cahce.countOnPage > 3
}

func isFifthElementVisible(cahce *walletsListDialogCache) bool {
	return cahce.countOnPage > 4
}

func getFirstItemText(cahce *walletsListDialogCache) string {
	index := cahce.currentPage * 4
	return cahce.cachedItems[int64(index)].text
}

func getSecondItemText(cahce *walletsListDialogCache) string {
	index := cahce.currentPage * 4 + 1
	return cahce.cachedItems[int64(index)].text
}

func getThirdItemText(cahce *walletsListDialogCache) string {
	index := cahce.currentPage * 4 + 2
	return cahce.cachedItems[int64(index)].text
}

func getFourthItemText(cahce *walletsListDialogCache) string {
	index := cahce.currentPage * 4 + 3
	return cahce.cachedItems[int64(index)].text
}

func getFifthItemText(cahce *walletsListDialogCache) string {
	index := cahce.currentPage * 4 + 4
	return cahce.cachedItems[int64(index)].text
}

func getFirstItemId(cahce *walletsListDialogCache) string {
	index := cahce.currentPage * 4
	return strconv.FormatInt(cahce.cachedItems[int64(index)].id, 10)
}

func getSecondItemId(cahce *walletsListDialogCache) string {
	index := cahce.currentPage * 4 + 1
	return strconv.FormatInt(cahce.cachedItems[int64(index)].id, 10)
}

func getThirdItemId(cahce *walletsListDialogCache) string {
	index := cahce.currentPage * 4 + 2
	return strconv.FormatInt(cahce.cachedItems[int64(index)].id, 10)
}

func getFourthItemId(cahce *walletsListDialogCache) string {
	index := cahce.currentPage * 4 + 3
	return strconv.FormatInt(cahce.cachedItems[int64(index)].id, 10)
}

func getFifthItemId(cahce *walletsListDialogCache) string {
	index := cahce.currentPage * 4 + 4
	return strconv.FormatInt(cahce.cachedItems[int64(index)].id, 10)
}

func addWallet(additionalId string, data *processing.ProcessData) bool {
	data.SubstitudeDialog(data.Static.MakeDialogFn("cw", data.UserId, data.Trans, data.Static))
	return true
}

func moveForward(additionalId string, data *processing.ProcessData) bool {
	ids, _ := data.Static.Db.GetUserWallets(data.UserId)
	itemsCount := len(ids)
	var pagesCount int
	if itemsCount > 2 {
		pagesCount = (itemsCount - 2) / 4 + 1
	} else {
		pagesCount = 1
	}

	currentPage := data.Static.GetUserStateCurrentPage(data.UserId)

	if currentPage + 1 < pagesCount {
		data.Static.SetUserStateCurrentPage(data.UserId, currentPage + 1)
	}
	data.SubstitudeDialog(data.Static.MakeDialogFn("wl", data.UserId, data.Trans, data.Static))
	return true
}

func moveBack(additionalId string, data *processing.ProcessData) bool {
	currentPage := data.Static.GetUserStateCurrentPage(data.UserId)
	if currentPage > 0 {
		data.Static.SetUserStateCurrentPage(data.UserId, currentPage - 1)
	}
	data.SubstitudeDialog(data.Static.MakeDialogFn("wl", data.UserId, data.Trans, data.Static))
	return true
}

func openWallet(additionalId string, data *processing.ProcessData) bool {
	id, err := strconv.ParseInt(additionalId, 10, 64)

	if err != nil {
		return false
	}

	if data.Static.Db.IsWalletBelongsToUser(data.UserId, id) {
		data.SubstitudeMessage("test open")
		//data.SubstitudeDialog(data.ChatId, data.Static.MakeDialogFn("li", id, data.Trans, data.Static))
		return true
	} else {
		return false
	}
}

func (factory *walletsListDialogFactory) createVariants(userId int64, trans i18n.TranslateFunc, staticData *processing.StaticProccessStructs) (variants []dialog.Variant) {
	variants = make([]dialog.Variant, 0)
	cache := getListDialogCache(userId, staticData)

	row := 1
	col := 0

	for _, variant := range factory.variants {
		if variant.isActiveFn == nil || variant.isActiveFn(cache) {
			var text string

			if variant.textFn != nil {
				text = variant.textFn(cache)
			} else {
				text = trans(variant.textId)
			}

			var additionalId string

			if variant.additionalIdFn != nil {
				additionalId = variant.additionalIdFn(cache)
			}

			variants = append(variants, dialog.Variant{
				Id:   variant.id,
				Text: text,
				AdditionalId: additionalId,
				RowId: row,
			})

			col = col + 1
			if col > 1 {
				row = row + 1
				col = 0
			}
		}
	}
	return
}

func getListDialogCache(userId int64, staticData *processing.StaticProccessStructs) (cache *walletsListDialogCache) {

	cache = &walletsListDialogCache{}

	cache.cachedItems = make([]cachedItem, 0)

	ids, names := staticData.Db.GetUserWallets(userId)
	if len(ids) == len(names) {
		for index, id := range ids {
			cache.cachedItems = append(cache.cachedItems, cachedItem{
				id: id,
				text: names[index],
			})
		}
	}

	cache.currentPage = staticData.GetUserStateCurrentPage(userId)
	count := len(cache.cachedItems)
	if count > 2 {
		cache.pagesCount = (count - 2) / 4 + 1
	} else {
		cache.pagesCount = 1
	}

	cache.countOnPage = count - cache.currentPage * 4
	if cache.countOnPage > 5 {
		cache.countOnPage = 4
	}

	return
}

func (factory *walletsListDialogFactory) MakeDialog(userId int64, trans i18n.TranslateFunc, staticData *processing.StaticProccessStructs) *dialog.Dialog {
	return &dialog.Dialog{
		Text:     trans(factory.textId),
		Variants: factory.createVariants(userId, trans, staticData),
	}
}

func (factory *walletsListDialogFactory) ProcessVariant(variantId string, additionalId string, data *processing.ProcessData) bool {
	for _, variant := range factory.variants {
		if variant.id == variantId {
			return variant.process(additionalId, data)
		}
	}
	return false
}
