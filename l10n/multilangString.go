package l10n

import (
	"github.com/nicksnyder/go-i18n/i18n"
	"log"
)

type MultilangString struct {
	content map[string]string
}

func MakeMultilangString(key string, trans map[string]i18n.TranslateFunc) (result MultilangString) {
	for lang, translator := range trans {
		result.content[lang] = translator(key)
	}
	return
}

func (multilangString *MultilangString) Get(language string) (result string) {
	result, ok := multilangString.content[language]
	if !ok {
		log.Printf("Key not found for language: %s", language)
	}
	return
}
