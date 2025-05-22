package i18n

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"

	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var (
	//go:embed locales/*.json
	localeFS embed.FS

	bundle *goi18n.Bundle
)

// InitBundle creates a Bundle with a fallback in English,
// Registers a JSON parser and loads all locales/*.json files
func InitBundle() error {
	bundle = goi18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	files, err := fs.Glob(localeFS, "locales/*.json")
	if err != nil {
		return fmt.Errorf("i18n: glob locale files: %w", err)
	}
	for _, file := range files {
		data, err := localeFS.ReadFile(file)
		if err != nil {
			return fmt.Errorf("i18n: read %s: %w", file, err)
		}
		if _, err := bundle.ParseMessageFileBytes(data, file); err != nil {
			return fmt.Errorf("i18n: parse %s: %w", file, err)
		}
	}
	return nil
}

// Localizer returns the localizer for the specified code (for example, "en" or "ru")
func Localizer(lang string) *goi18n.Localizer {
	if bundle == nil {
		log.Println("i18n: bundle not initialized, initializing")
		if err := InitBundle(); err != nil {
			log.Fatalf("i18n: InitBundle failed: %v", err)
		}
	}
	return goi18n.NewLocalizer(bundle, lang)
}
