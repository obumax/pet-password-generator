package session

import "time"

type Session struct {
	ChatID         int64     `json:"chat_id"`         // Telegram chat ID
	Language       string    `json:"language"`        // Selected language
	State          string    `json:"state"`           // Current status
	PasswordLength int       `json:"password_length"` // Selected password length
	Flags          []string  `json:"flags"`           // Selected flags
	LastActive     time.Time `json:"last_active"`     // Time of last activity
}

// Интерфейс Store используется для работы с сессиями
type Store interface {
	// Get возвращает сессию или ошибку
	Get(chatID int64) (*Session, error)
	// Set сохраняет сессию по установленному TTL
	Set(chatID int64, sess *Session) error
	// Delete удаляет сессию
	Delete(chatID int64) error
}

// TTLSession - время жизни сессии в Redis
const TTLSession = 24 * time.Hour
