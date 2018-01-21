package database

import (
	"log"
)

const (
	minimalVersion = "0.1"
	latestVersion  = "0.1"
)

type dbUpdater struct {
	version  string
	updateDb func(db *AccountDb)
}

func UpdateVersion(db *AccountDb) {
	currentVersion := db.GetDatabaseVersion()

	if currentVersion != latestVersion {
		updaters := makeUpdaters(currentVersion, latestVersion)

		for _, updater := range updaters {
			updater.updateDb(db)
		}
		log.Printf("Update DB version from %s to %s", currentVersion, latestVersion)
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
