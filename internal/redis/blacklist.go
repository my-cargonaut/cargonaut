package redis

import (
	"context"
	"time"

	"github.com/gomodule/redigo/redis"
	uuid "github.com/satori/go.uuid"

	"github.com/my-cargonaut/cargonaut"
)

var _ cargonaut.TokenBlacklist = (*TokenBlacklist)(nil)

// TokenBlacklist is a Redis based token blacklist.
type TokenBlacklist struct {
	conn redis.Conn
}

// NewTokenBlacklist returns a new token blacklist based on the provided redis
// connection.
func NewTokenBlacklist(conn redis.Conn) *TokenBlacklist {
	return &TokenBlacklist{conn}
}

// IsTokenBlacklisted retrieves a token by its unique token ID and returns true
// if it is blacklisted.
func (b *TokenBlacklist) IsTokenBlacklisted(ctx context.Context, id uuid.UUID) (bool, error) {
	// Get the token identified by its ID from the token blacklist. If the
	// response is anything else but nil, the requested token is on the
	// blacklist.
	resp, err := b.conn.Do("GET", id)
	if err != nil {
		return true, err
	}
	return resp != nil, nil
}

// BlacklistToken blacklists one or more tokens by putting them onto the token
// blacklist. The tokens are identified by their unique token ID.
func (b *TokenBlacklist) BlacklistToken(ctx context.Context, tokens ...*cargonaut.Token) error {
	for _, token := range tokens {
		// Skip expired tokens. Adding them would return an error from Redis
		// because the tokens expiration time is below the current time.
		if time.Now().After(token.ExpiresAt) {
			continue
		} else if err := b.blacklistToken(token); err != nil {
			return err
		}
	}
	return nil
}

// blacklistToken blacklists a token by putting it onto the token blacklist. The
// token is identified by its unique token ID.
func (b *TokenBlacklist) blacklistToken(token *cargonaut.Token) error {
	// Set the token onto the blacklist until it expires.
	exp := time.Until(token.ExpiresAt).Seconds()
	if _, err := b.conn.Do("SETEX", token.ID, int(exp), "blacklisted"); err != nil {
		return err
	}
	return nil
}
