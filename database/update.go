package database

import (
	"bytes"
	"fmt"
	"gitlab.com/gameraccoon/telegram-accountant-bot/currencies"
	"log"
)

const (
	minimalVersion = "0.1"
	latestVersion  = "0.3"
)

type dbUpdater struct {
	version  string
	updateDb func(db *AccountDb)
}

func UpdateVersion(db *AccountDb) {
	currentVersion := db.GetDatabaseVersion()

	if currentVersion != latestVersion {
		updaters := makeUpdaters(currentVersion, latestVersion)

		log.Printf("Update DB version from %s to %s in %d iterations", currentVersion, latestVersion, len(updaters))
		for _, updater := range updaters {
			log.Printf("Updating to %s", updater.version)
			updater.updateDb(db)
		}
	}

	db.SetDatabaseVersion(latestVersion)
}

func makeUpdaters(versionFrom string, versionTo string) (updaters []dbUpdater) {
	allUpdaters := makeAllUpdaters()

	isFirstFound := (versionFrom == minimalVersion)
	for _, updater := range allUpdaters {
		if isFirstFound {
			updaters = append(updaters, updater)
			if updater.version == versionTo {
				break
			}
		} else {
			if updater.version == versionFrom {
				isFirstFound = true
			}
		}
	}

	if len(updaters) > 0 {
		lastFoundVersion := updaters[len(updaters) - 1].version
		if lastFoundVersion != versionTo {
			log.Fatalf("Last version updater not found. Expected: %s Found: %s", versionTo, lastFoundVersion)
		}
	}
	return
}

func makeAllUpdaters() (updaters []dbUpdater) {
	updaters = []dbUpdater{
		dbUpdater{
			version: "0.2",
			updateDb: func(db *AccountDb) {
				db.db.Exec("ALTER TABLE wallets ADD COLUMN contract_address TEXT NOT NULL DEFAULT('')")
			},
		},
		dbUpdater{
			version: "0.3",
			updateDb: func(db *AccountDb) {
				db.db.Exec("ALTER TABLE wallets ADD COLUMN price_id TEXT NOT NULL DEFAULT('')")
				availableCurrencies := currencies.GetAllCurrencies()
				var b bytes.Buffer
				for _, currency := range availableCurrencies {
					priceId := currencies.GetCurrencyPriceId(currency)
					b.WriteString(fmt.Sprintf("UPDATE wallets SET price_id='%s' WHERE currency=%d;", priceId, currency))
				}
				db.db.Exec(b.String())
			},
		},
	}
	return
}
