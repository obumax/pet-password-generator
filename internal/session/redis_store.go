package session

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	// ErrNotFound is returned if there is no state in Redis for the given chatID
	ErrNotFound = fmt.Errorf("session not found")
)

// RedisStore stores sessions in Redis
type RedisStore struct {
	client *redis.Client
	prefix string
	ttl    time.Duration
}

// NewRedisStore creates a store instance
// redisAddr - e.g. "redis:6379", db - database number, pwd - password ("" if none)
func NewRedisStore(redisAddr string, db int, pwd string) *RedisStore {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		DB:       db,
		Password: pwd,
	})
	return &RedisStore{
		client: rdb,
		prefix: "session:",
		ttl:    TTLSession,
	}
}

// key generates a key for Redis from chatID
func (r *RedisStore) key(chatID int64) string {
	return fmt.Sprintf("%s%d", r.prefix, chatID)
}

// Get loads the session from Redis
func (r *RedisStore) Get(chatID int64) (*Session, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	data, err := r.client.Get(ctx, r.key(chatID)).Result()
	if err == redis.Nil {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	var sess Session
	if err := json.Unmarshal([]byte(data), &sess); err != nil {
		return nil, err
	}
	return &sess, nil
}

// Set serializes and stores a session with TTL
func (r *RedisStore) Set(chatID int64, sess *Session) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	sess.LastActive = time.Now()
	buf, err := json.Marshal(sess)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, r.key(chatID), buf, r.ttl).Err()
}

// Delete deletes the session key
func (r *RedisStore) Delete(chatID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := r.client.Del(ctx, r.key(chatID)).Result()
	return err
}
