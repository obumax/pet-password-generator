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

// The Store interface is used to work with sessions
type Store interface {
	// Get returns session or error
	Get(chatID int64) (*Session, error)
	// Set saves the session by the set TTL
	Set(chatID int64, sess *Session) error
	// Delete deletes the session
	Delete(chatID int64) error
}

var defaultStore Store

func InitStore(s Store) {
	defaultStore = s
}

func SetLang(chatID int64, lang string) error {
	if defaultStore == nil {
		return nil
	}
	sess, err := defaultStore.Get(chatID)
	if err != nil && err != ErrNotFound {
		return err
	}
	if sess == nil {
		sess = &Session{ChatID: chatID}
	}
	sess.Language = lang
	return defaultStore.Set(chatID, sess)
}

func GetLang(chatID int64) string {
	if defaultStore == nil {
		return ""
	}
	sess, err := defaultStore.Get(chatID)
	if err != nil {
		return ""
	}
	return sess.Language
}

const TTLSession = 24 * time.Hour
