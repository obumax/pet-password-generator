package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/redis/go-redis/v9"

	"github.com/obumax/pet-password-generator/internal/i18n"
	"github.com/obumax/pet-password-generator/internal/session"
)

func init() {
	_ = godotenv.Load() // .env
}

func main() {

	// Localization
	if err := i18n.InitBundle(); err != nil {
		log.Fatalf("i18n init failed: %v", err)
	}

	// Redis-store
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Fatal("REDIS_URL is not set")
	}
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("invalid REDIS_URL: %v", err)
	}
	store := session.NewRedisStore(opt.Addr, opt.DB, opt.Password)
	session.InitStore(store)

	// Telegram-bot
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatal("BOT_TOKEN is not set")
	}
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("NewBotAPI error: %v", err)
	}
	log.Printf("Authorized on %s", bot.Self.UserName)

	updCfg := tgbotapi.NewUpdate(0)
	updCfg.Timeout = 30
	updates := bot.GetUpdatesChan(updCfg)

	for upd := range updates {
		if upd.CallbackQuery != nil {
			handleCallback(bot, upd.CallbackQuery)
			continue
		}
		if upd.Message != nil {
			handleMessage(bot, upd.Message)
		}
	}
}

func handleCallback(bot *tgbotapi.BotAPI, cq *tgbotapi.CallbackQuery) {
	data := cq.Data
	if !strings.HasPrefix(data, "lang:") {
		return
	}
	parts := strings.SplitN(data, ":", 2)
	if len(parts) != 2 || parts[1] == "" {
		bot.Request(tgbotapi.NewCallback(cq.ID, "Error!"))
		return
	}
	lang := parts[1]
	chatID := cq.Message.Chat.ID
	if err := session.SetLang(chatID, lang); err != nil {
		log.Printf("SetLang error: %v", err)
	}
	bot.Request(tgbotapi.NewCallback(cq.ID, ""))

	loc := i18n.Localizer(lang)
	sendLocalized(bot, chatID, loc, "start_greeting", nil)
	sendLocalized(bot, chatID, loc, "start_commands", nil)
}

func handleMessage(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	lang := session.GetLang(chatID)
	if lang == "" {
		lang = "en"
		session.SetLang(chatID, lang)
	}
	loc := i18n.Localizer(lang)

	if !msg.IsCommand() {
		return
	}
	switch msg.Command() {
	case "start":
		sendLocalized(bot, chatID, loc, "start_choose_language", nil)
	case "generate":
		handleGenerate(bot, chatID, msg.CommandArguments(), loc)
	default:
		sendLocalized(bot, chatID, loc, "unknown_command", nil)
	}
}

func sendLocalized(bot *tgbotapi.BotAPI, chatID int64, loc *goi18n.Localizer,
	messageID string, data map[string]interface{},
) {
	text, err := loc.Localize(&goi18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: data,
	})
	if err != nil {
		text = fmt.Sprintf("[%s]", messageID)
		log.Printf("loc err %q: %v", messageID, err)
	}
	if _, err := bot.Send(tgbotapi.NewMessage(chatID, text)); err != nil {
		log.Printf("send err: %v", err)
	}
}
