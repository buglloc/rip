package www

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"

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
	channel := genChannelID()
	t := jwt.New()
	if err := t.Set(jwt.SubjectKey, channel); err != nil {
		return "", fmt.Errorf("can't set 'sub' to JWT token: %w", err)
	}

	if err := t.Set(jwt.IssuedAtKey, time.Now()); err != nil {
		return "", fmt.Errorf("can't set 'iat' to JWT token: %w", err)
	}

	token, err := jwt.Sign(t, jwa.HS256, c.signKey)
	if err != nil {
		return "", fmt.Errorf("can't sign token: %w", err)
	}

	return string(token), nil
}

func (c *tokenManager) ParseToken(in string) (string, error) {
	token, err := jwt.Parse([]byte(in), jwt.WithVerify(jwa.HS256, c.signKey))
	if err != nil {
		return "", err
	}

	iat := token.IssuedAt()
	if time.Since(iat) >= c.ttl {
		return "", fmt.Errorf("expired token")
	}

	return token.Subject(), nil
}

func genChannelID() string {
	bs := make([]byte, 6)
	rand.Read(bs[:2])
	binary.BigEndian.PutUint32(bs[2:], uint32(time.Now().Unix()))
	return hex.EncodeToString(bs)
}
