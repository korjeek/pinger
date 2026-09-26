package usecase

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/korjeek/pinger/backend/internal/dto"
)

type TokenGeneratorConfig struct {
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

type TokenGenerator struct {
	clock Clock

	accessTokenSecret []byte
	accessTokenTTL    time.Duration

	refreshTokenTTL time.Duration
}

type AccessClaims struct {
	jwt.RegisteredClaims
}

func NewTokenGenerator(cfg TokenGeneratorConfig, clock Clock) *TokenGenerator {
	return &TokenGenerator{
		accessTokenSecret: cfg.ATC.Secret,
		accessTokenTTL:    cfg.ATC.TTL,
		refreshTokenTTL:   cfg.RTC.TTL,
		clock:             clock,
	}
}

func (g *TokenGenerator) GenerateAccessToken(userId uuid.UUID) (dto.Token, error) {
	now := g.clock.Now()
	expiresAt := now.Add(g.accessTokenTTL)

	claims := AccessClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userId.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedJWT, err := token.SignedString(g.accessTokenSecret)

	return dto.Token{
		Payload:   signedJWT,
		ExpiresAt: expiresAt,
	}, err
}

func (g *TokenGenerator) GenerateRefreshToken() (refToken dto.Token, error error) {
	now := g.clock.Now()
	buf := make([]byte, 32)

	if _, err := rand.Read(buf); err != nil {
		return refToken, err
	}

	token := base64.RawURLEncoding.EncodeToString(buf)
	return dto.Token{
		Payload:   token,
		ExpiresAt: now.Add(g.refreshTokenTTL),
	}, nil
}

func (g *TokenGenerator) GenerateTokenPair(userId uuid.UUID) (pair dto.TokenPair, err error) {
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
