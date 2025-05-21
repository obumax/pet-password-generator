package session

import (
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRedisStoreGetSet(t *testing.T) {
	client, mock := redismock.NewClientMock()
	store := &RedisStore{
		client: client,
		prefix: "session",
		ttl:    24 * time.Hour,
	}

	chatID := int64(12345)
	key := store.key(chatID)

	sess := &Session{
		ChatID:     chatID,
		Language:   "ru",
		State:      "init",
		LastActive: time.Now(),
	}

	// Mock Get(redis.Nil - session not found)
	mock.ExpectGet(key).RedisNil()
	_, err := store.Get(chatID)
	assert.Equal(t, ErrNotFound, err)

	// Mock Set and Get
	mock.ExpectSet(key, gomock.Any(), store.ttl).SetVal("OK")
	mock.ExpectGet(key).SetVal(`{"chat_id":12345,"language":"ru","state":"init","last_active":"` + sess.LastActive.Format(time.RFC3339Nano) + `"}`)
	err = store.Set(chatID, sess)
	assert.NoError(t, err)

	got, err := store.Get(chatID)
	assert.NoError(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, "ru", got.Language)
	assert.Equal(t, chatID, got.ChatID)

	// Check that the mock calls were made
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRedisStoreDelete(t *testing.T) {
	client, mock := redismock.NewClientMock()
	store := &RedisStore{
		client: client,
		prefix: "session",
		ttl:    24 * time.Hour,
	}

	chatID := int64(12345)
	key := store.key(chatID)

	mock.ExpectDel(key).SetVal(1)
	err := store.Delete(chatID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
