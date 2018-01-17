package database

import (
	"github.com/gameraccoon/telegram-bot-skeleton/database"
	"log"
)

const (
	minimalVersion = "0.1"
	latestVersion  = "0.1"
)

type dbUpdater struct {
	version  string
	updateDb func(db *database.Database)
}

func UpdateVersion(db *database.Database) {
	currentVersion := GetDatabaseVersion(db)

	if currentVersion != latestVersion {
		updaters := makeUpdaters(currentVersion, latestVersion)

		for _, updater := range updaters {
			updater.updateDb(db)
		}
		log.Printf("Update DB version from %s to %s", currentVersion, latestVersion)
	}

	SetDatabaseVersion(db, latestVersion)
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
				updaters = append(updaters, updater)
			}
		}
	}
	return
}

func makeAllUpdaters() (updaters []dbUpdater) {
	updaters = []dbUpdater{
		dbUpdater{
			// version: "1.1",
			// updateDb: func(db *Database) {
			// 	db.execQuery("ALTER TABLE prohibited_words ADD COLUMN removed INTEGER")
			// },
		},
	}
	return
}
