package main

import (
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"

	"github.com/obumax/pet-password-generator/internal/i18n"
	"github.com/obumax/pet-password-generator/internal/session"
)

func main() {
	// Initialize the translation bundle
	i18n.InitBundle()

	// Create a Telegram Bot
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatal("BOT_TOKEN is not set")
	}
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)

	updates := bot.GetUpdatesChan(tgbotapi.NewUpdate(0))

	for upd := range updates {
		// 1) Regular text messages
		if upd.Message != nil {
			chatID := upd.Message.Chat.ID
			lang := session.GetLang(chatID) // "ru" or "en"
			loc := i18n.Localizer(lang)     // *goi18n.Localizer

			if upd.Message.IsCommand() {
				switch upd.Message.Command() {
				case "start":
					// Suggest language choice
					msg := tgbotapi.NewMessage(chatID,
						mustLocalize(loc, "start_choose_language", nil))
					kb := tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("🇷🇺 Русский", "lang:ru"),
							tgbotapi.NewInlineKeyboardButtonData("🇬🇧 English", "lang:en"),
						),
					)
					msg.ReplyMarkup = kb
					bot.Send(msg)

				case "generate":
					args := upd.Message.CommandArguments()
					handleGenerate(bot, chatID, args, loc)

				default:
					bot.Send(tgbotapi.NewMessage(chatID,
						mustLocalize(loc, "unknown_command", nil)))
				}
			}
		}

		// 2) Inline-callback (language selection)
		if upd.CallbackQuery != nil {
			data := upd.CallbackQuery.Data
			if strings.HasPrefix(data, "lang:") {
				parts := strings.SplitN(data, ":", 2)
				lang := parts[1]
				session.SetLang(upd.CallbackQuery.Message.Chat.ID, lang)

				// // Confirm the callback to remove the "loading"
				bot.Request(tgbotapi.NewCallback(upd.CallbackQuery.ID, ""))

				loc := i18n.Localizer(lang)
				bot.Send(tgbotapi.NewMessage(upd.CallbackQuery.Message.Chat.ID,
					mustLocalize(loc, "start_greeting", nil)))
				bot.Send(tgbotapi.NewMessage(upd.CallbackQuery.Message.Chat.ID,
					mustLocalize(loc, "start_commands", nil)))
			}
		}
	}
}

// mustLocalize is halper to avoid handling the error every time
func mustLocalize(loc *goi18n.Localizer, messageID string, data map[string]interface{}) string {
	s, err := loc.Localize(&goi18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: data,
	})
	if err != nil {
		return messageID // fallback: just ID
	}
	return s
}
