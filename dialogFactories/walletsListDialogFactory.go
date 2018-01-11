package dialogFactories

import (
	"github.com/gameraccoon/telegram-bot-skeleton/dialog"
	"github.com/gameraccoon/telegram-bot-skeleton/dialogFactory"
	"github.com/gameraccoon/telegram-bot-skeleton/processing"
	"gitlab.com/gameraccoon/telegram-accountant-bot/database"
	"gitlab.com/gameraccoon/telegram-accountant-bot/cryptoFunctions"
	"gitlab.com/gameraccoon/telegram-accountant-bot/currencies"
	"github.com/nicksnyder/go-i18n/i18n"
	"log"
	"math/big"
	"strconv"
)

const maxItemsOnPage int = 10
const maxItemsInRow int = 2

type walletsListDialogVariantPrototype struct {
	isListItem bool
	id string
	additionalIdFn func(*walletsListDialogCache, int) string
	textId string
	textFn func(*walletsListDialogCache, int) string
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
	variants []walletsListDialogVariantPrototype
}

func MakeWalletsListDialogFactory() dialogFactory.DialogFactory {
	return &(walletsListDialogFactory{
		variants: []walletsListDialogVariantPrototype{
			walletsListDialogVariantPrototype{
				id: "add",
				textId: "add_wallet_btn",
				isActiveFn: isTheFirstPage,
				process: addWallet,
			},
			walletsListDialogVariantPrototype{
				isListItem: true,
				id: "it",
				additionalIdFn: getItemId,
				textFn: getItemText,
				process: openWallet,
			},
			walletsListDialogVariantPrototype{
				id: "bck",
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

func getItemText(cahce *walletsListDialogCache, itemIndex int) string {
	index := cahce.currentPage * maxItemsOnPage + itemIndex
	return cahce.cachedItems[int64(index)].text
}

func getItemId(cahce *walletsListDialogCache, itemIndex int) string {
	index := cahce.currentPage * maxItemsOnPage + itemIndex
	return strconv.FormatInt(cahce.cachedItems[int64(index)].id, 10)
}

func addWallet(additionalId string, data *processing.ProcessData) bool {
	data.SubstitudeDialog(data.Static.MakeDialogFn("cc", data.UserId, data.Trans, data.Static))
	return true
}

func moveForward(additionalId string, data *processing.ProcessData) bool {
	ids, _ := database.GetUserWallets(data.Static.Db, data.UserId)
	itemsCount := len(ids)
	var pagesCount int
	if itemsCount > 2 {
		pagesCount = (itemsCount - 2) / maxItemsOnPage + 1
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

	if database.IsWalletBelongsToUser(data.Static.Db, data.UserId, id) {
		data.SubstitudeDialog(data.Static.MakeDialogFn("wa", id, data.Trans, data.Static))
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

		if variant.isListItem {
			for i := 0; i < cache.countOnPage; i++ {
				if variant.textFn == nil || variant.additionalIdFn == nil {
					log.Printf("List element doesn't have a valid functions")
					continue
				}

				variants = append(variants, dialog.Variant{
					Id:   variant.id + strconv.Itoa(i),
					Text: variant.textFn(cache, i),
					RowId: row,
					AdditionalId: variant.additionalIdFn(cache, i),
				})

				col = col + 1
				if col >= maxItemsInRow {
					row = row + 1
					col = 0
				}
			}
		} else {
			if variant.isActiveFn == nil || variant.isActiveFn(cache) {
				variants = append(variants, dialog.Variant{
					Id:   variant.id,
					Text: trans(variant.textId),
					RowId: row,
				})

				col = col + 1
				if col >= maxItemsInRow {
					row = row + 1
					col = 0
				}
			}
		}
	}
	return
}

func getListDialogCache(userId int64, staticData *processing.StaticProccessStructs) (cache *walletsListDialogCache) {

	cache = &walletsListDialogCache{}

	cache.cachedItems = make([]cachedItem, 0)

	ids, names := database.GetUserWallets(staticData.Db, userId)
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
		cache.pagesCount = (count - 2) / maxItemsOnPage + 1
	} else {
		cache.pagesCount = 1
	}

	cache.countOnPage = count - cache.currentPage * maxItemsOnPage
	if cache.countOnPage > maxItemsOnPage + 1 {
		cache.countOnPage = maxItemsOnPage
	}

	return
}

func (factory *walletsListDialogFactory) GetDialogCaption(userId int64, trans i18n.TranslateFunc, staticData *processing.StaticProccessStructs) string {
	walletAddresses := database.GetUserWalletAddresses(staticData.Db, userId)

	if len(walletAddresses) == 0 {
		return ""
	}

	groupedWallets := make(map[currencies.Currency] []string)

	for _, walletAddress := range walletAddresses {
		walletsSlice, ok := groupedWallets[walletAddress.Currency]
		if ok {
			groupedWallets[walletAddress.Currency] = append(walletsSlice, walletAddress.Address)
		} else {
			groupedWallets[walletAddress.Currency] = []string{ walletAddress.Address }
		}
	}

	text := "Balance:\n"
	
	usdSum := new(big.Float)

	for currency, wallets := range groupedWallets {
		processor := cryptoFunctions.GetProcessor(currency)

		if processor == nil {
			continue
		}

		var balance *big.Int
		if processor != nil {
			balance = (*processor).GetSumBalance(wallets)
		}
		
		if balance == nil {
			return "Error"
		}
		
		currencyCode := currencies.GetCurrencyCode(currency)
		currencyDigits := currencies.GetCurrencyDigits(currency)
		
		floatBalance := cryptoFunctions.GetFloatBalance(balance, currencyDigits)
		
		if floatBalance == nil {
			return "Error"
		}
		
		if processor != nil {
			toUsdRate := (*processor).GetToUsdRate()

			if toUsdRate != nil {
				usdSum.Add(usdSum, new(big.Float).Mul(floatBalance, toUsdRate))
			}
		}

		text = text + floatBalance.Text('f', currencyDigits) + " " + currencyCode + "\n"
	}
	
	if usdSum != nil {
		text = text + "Sum: " + usdSum.Text('f', 2) + " USD\n"
	}

	return text
}

func (factory *walletsListDialogFactory) MakeDialog(userId int64, trans i18n.TranslateFunc, staticData *processing.StaticProccessStructs) *dialog.Dialog {
	return &dialog.Dialog{
		Text:     factory.GetDialogCaption(userId, trans, staticData) + trans("choose_wallet"),
		Variants: factory.createVariants(userId, trans, staticData),
	}
}

func (factory *walletsListDialogFactory) ProcessVariant(variantId string, additionalId string, data *processing.ProcessData) bool {
	for _, variant := range factory.variants {
		if variant.isListItem {
			if variant.id == variantId[0:2] { // "id"
				return variant.process(additionalId, data)
			}
		} else if variant.id == variantId {
			return variant.process(additionalId, data)
		}
	}
	return false
}
