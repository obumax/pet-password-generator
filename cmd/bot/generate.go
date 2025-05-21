package main

import (
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"

	"github.com/obumax/pet-password-generator/internal/generator"
	i18nutil "github.com/obumax/pet-password-generator/internal/i18n"
)

const (
	minLen = 4
	maxLen = 35
)

// handleGenerate parses the arguments, calls generator.Generate,
// localizes the error or success and sends the response
func handleGenerate(
	bot *tgbotapi.BotAPI,
	chatID int64,
	args string,
	loc *goi18n.Localizer,
) {
	parts := strings.Fields(args)
	if len(parts) < 2 {
		// Hint on using the command
		usage := mustLocalize(loc, "prompt_length", map[string]interface{}{"Min": minLen, "Max": maxLen}) + "\n" +
			mustLocalize(loc, "prompt_flags", nil) + "\n" +
			"/generate 12 ULDSX"
		bot.Send(tgbotapi.NewMessage(chatID, usage))
		return
	}

	// 1) Parse the length
	length, err := strconv.Atoi(parts[0])
	if err != nil {
		bot.Send(tgbotapi.NewMessage(chatID, mustLocalize(loc, "unknown_command", nil)))
		return
	}

	// 2) flags
	var flags generator.FlagsSet
	for _, c := range parts[1] {
		switch c {
		case 'U', 'u':
			flags.Upper = true
		case 'L', 'l':
			flags.Lower = true
		case 'D', 'd':
			flags.Digits = true
		case 'S', 's':
			flags.SpecSymbols = true
		case 'X', 'x':
			flags.ExcludeSimilar = true
		}
	}

	// 3) Generate a password
	pass, err := generator.Generate(length, flags)
	if err != nil {
		text := i18nutil.LocalizeError(loc, err, map[string]interface{}{
			"Min": minLen,
			"Max": maxLen,
		})
		bot.Send(tgbotapi.NewMessage(chatID, text))
		return
	}

	// 4) Success
	success := mustLocalize(loc, "generation_success", map[string]interface{}{"Password": pass})
	bot.Send(tgbotapi.NewMessage(chatID, success))
}
