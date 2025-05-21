package i18n

import (
	"embed"
	"encoding/json"
	"log"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed ../../configs/locales/*.json
var localeFS embed.FS

var bundle *i18n.Bundle

// InitBundle загружает все JSON-файлы локалей из configs/locales
// InitBundle loads all locale JSON files from configs/locales
func InitBundle() {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	files, err := localeFS.ReadDir("configs/locales")
	if err != nil {
		log.Fatalf("i18n: cannot read locales dir: %v", err)
	}
	for _, f := range files {
		data, err := localeFS.ReadFile("configs/locales/" + f.Name())
		if err != nil {
			log.Fatalf("i18n: cannot read %s: %v", f.Name(), err)
		}
		if _, err := bundle.ParseMessageFileBytes(data, f.Name()); err != nil {
			log.Fatalf("i18n: parse %s failed: %v", f.Name(), err)
		}
	}
}

// Localizer возвращает локализатор для данного языка (например, "ru" или "en")
// Localizer returns the localizer for the given language (for example, "ru" or "en")
func Localizer(lang string) *i18n.Localizer {
	if bundle == nil {
		InitBundle()
	}
	return i18n.NewLocalizer(bundle, lang)
}
