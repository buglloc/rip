package www

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwt"

	"github.com/buglloc/rip/v2/pkg/cfg"
)

type tokenManager struct {
	signKey []byte
	ttl     time.Duration
}

func NewTokenManager() *tokenManager {
	return &tokenManager{
		signKey: []byte(cfg.HubSign),
		ttl:     cfg.HubSignTTL,
	}
}

func (c *tokenManager) NewToken() (string, error) {
	channelID, err := newChannelID()
	if err != nil {
		return "", fmt.Errorf("can't create channel id: %w", err)
	}

	t := jwt.New()
	if err := t.Set(jwt.SubjectKey, channelID); err != nil {
		return "", fmt.Errorf("can't set 'sub' to JWT token: %w", err)
	}

	if err := t.Set(jwt.IssuedAtKey, time.Now()); err != nil {
		return "", fmt.Errorf("can't set 'iat' to JWT token: %w", err)
	}

	token, err := jwt.Sign(t, jwt.WithKey(jwa.HS256(), c.signKey))
	if err != nil {
		return "", fmt.Errorf("can't sign token: %w", err)
	}

	return string(token), nil
}

func (c *tokenManager) ParseToken(in string) (string, error) {
	token, err := jwt.Parse([]byte(in), jwt.WithKey(jwa.HS256(), c.signKey))
	if err != nil {
		return "", err
	}

	iat, ok := token.IssuedAt()
	if !ok {
		return "", errors.New("missing token issue time")
	}

	if time.Since(iat) >= c.ttl {
		return "", errors.New("expired token")
	}

	subject, ok := token.Subject()
	if !ok {
		return "", errors.New("missing token subject")
	}

	return subject, nil
}

func newChannelID() (string, error) {
	bs := make([]byte, 6)
	if _, err := rand.Read(bs[:2]); err != nil {
		return "", fmt.Errorf("rand read: %w", err)
	}

	binary.BigEndian.PutUint32(bs[2:], uint32(time.Now().Unix()))
	return hex.EncodeToString(bs), nil
}
