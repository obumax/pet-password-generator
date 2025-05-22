package main

import (
	"log"
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

func handleGenerate(
	bot *tgbotapi.BotAPI,
	chatID int64,
	args string,
	loc *goi18n.Localizer,
) {
	parts := strings.Fields(args)
	if len(parts) < 2 {
		sendLocalized(bot, chatID, loc, "prompt_length", map[string]interface{}{
			"Min": minLen, "Max": maxLen,
		})
		sendLocalized(bot, chatID, loc, "prompt_flags", nil)
		return
	}

	length, err := strconv.Atoi(parts[0])
	if err != nil || length < minLen || length > maxLen {
		sendLocalized(bot, chatID, loc, "prompt_length", map[string]interface{}{
			"Min": minLen, "Max": maxLen,
		})
		return
	}

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
	if !flags.HasAny() {
		sendLocalized(bot, chatID, loc, "prompt_flags", nil)
		return
	}

	pass, err := generator.Generate(length, flags)
	if err != nil {
		text := i18nutil.MapError(loc, err, map[string]interface{}{
			"Min": minLen, "Max": maxLen,
		})
		if _, e := bot.Send(tgbotapi.NewMessage(chatID, text)); e != nil {
			log.Printf("send err: %v", e)
		}
		return
	}

	sendLocalized(bot, chatID, loc, "generation_success", map[string]interface{}{
		"Password": pass,
	})
}
