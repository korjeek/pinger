package token

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/korjeek/pinger/backend/internal/dto"
	"github.com/lestrrat-go/jwx/v4/jwa"
	"github.com/lestrrat-go/jwx/v4/jwt"
)

var (
	ErrEmptySecret           = errors.New("access token secret cannot be empty")
	ErrFailedToGenerateBytes = errors.New("failed to generate random bytes for refresh token")
	ErrFailedToBuildToken    = errors.New("failed to build access token")
	ErrFailedToSignToken     = errors.New("failed to sign access token")

	ErrInvalidToken   = errors.New("invalid access token")
	ErrMissingSubject = errors.New("access token has no subject claim")
	ErrInvalidSubject = errors.New("access token subject is not a valid UUID")
)

type GeneratorConfig struct {
	ATC AccessTokenConfig
	RTC RefreshTokenConfig
}

type AccessTokenConfig struct {
	Secret []byte
	TTL    time.Duration
}

type RefreshTokenConfig struct {
	TTL time.Duration
}

type Clock interface {
	Now() time.Time
}

type JWTTokenManager struct {
	clock Clock

	accessTokenSecret []byte
	accessTokenTTL    time.Duration

	refreshTokenTTL time.Duration
}

func NewGenerator(cfg GeneratorConfig, clock Clock) (*JWTTokenManager, error) {
	if len(cfg.ATC.Secret) == 0 {
		return nil, ErrEmptySecret
	}

	return &JWTTokenManager{
		accessTokenSecret: cfg.ATC.Secret,
		accessTokenTTL:    cfg.ATC.TTL,
		refreshTokenTTL:   cfg.RTC.TTL,
		clock:             clock,
	}, nil
}

func (g *JWTTokenManager) GenerateAccessToken(userId uuid.UUID) (dto.Token, error) {
	now := g.clock.Now()
	expiresAt := now.Add(g.accessTokenTTL)

	tok, err := jwt.NewBuilder().
		Subject(userId.String()).
		IssuedAt(now).
		Expiration(expiresAt).
		Build()
	if err != nil {
		return dto.Token{}, fmt.Errorf("%w: %v", ErrFailedToBuildToken, err)
	}

	signed, err := jwt.Sign(tok, jwt.WithKey(jwa.HS256(), g.accessTokenSecret))
	if err != nil {
		return dto.Token{}, fmt.Errorf("%w: %v", ErrFailedToSignToken, err)
	}

	return dto.Token{
		Payload:   string(signed),
		ExpiresAt: expiresAt,
	}, err
}

func (g *JWTTokenManager) GenerateRefreshToken() (refToken dto.Token, error error) {
	now := g.clock.Now()
	buf := make([]byte, 32)

	if _, err := rand.Read(buf); err != nil {
		return dto.Token{}, fmt.Errorf("%w: %v", ErrFailedToGenerateBytes, err)
	}

	token := base64.RawURLEncoding.EncodeToString(buf)
	return dto.Token{
		Payload:   token,
		ExpiresAt: now.Add(g.refreshTokenTTL),
	}, nil
}

func (g *JWTTokenManager) GenerateTokenPair(userId uuid.UUID) (pair dto.TokenPair, err error) {
	accessToken, err := g.GenerateAccessToken(userId)
	if err != nil {
		return pair, err
	}

	refreshToken, err := g.GenerateRefreshToken()
	if err != nil {
		return pair, err
	}

	return dto.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (m *JWTTokenManager) Validate(tokenString string) (uuid.UUID, error) {
	token, err := jwt.Parse(
		[]byte(tokenString),
		jwt.WithKey(jwa.HS256(), m.accessTokenSecret),
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	sub, ok := token.Subject()
	if !ok || sub == "" {
		return uuid.Nil, ErrMissingSubject
	}

	userID, err := uuid.Parse(sub)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: %v", ErrInvalidSubject, err)
	}

	return userID, nil
}
