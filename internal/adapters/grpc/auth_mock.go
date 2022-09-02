package grpc

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/rs/zerolog"
	"gitlab.com/g6834/team17/task-service/internal/config"
	"gitlab.com/g6834/team17/task-service/internal/domain/models"
	"strconv"
)

type AuthClientMock struct {
	secretKey string
	logger    *zerolog.Logger
}

func NewAuthClientMock(cfg *config.Config, l *zerolog.Logger) *AuthClientMock {
	return &AuthClientMock{
		logger:    l,
		secretKey: cfg.Auth.SecretKey,
	}
}

func (c *AuthClientMock) Validate(ctx context.Context, token *models.TokenPair) (*models.TokenPair, error) {
	return token, nil
}

func (c *AuthClientMock) ParseToken(ctx context.Context, tokenString string) (*models.User, bool, error) {
	claims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(c.secretKey), nil
	})
	if err != nil {
		return nil, false, err
	}

	stringId := fmt.Sprintf("%v", claims[userId])
	userID, err := strconv.ParseUint(stringId, 10, 64)
	if err != nil {
		return nil, false, err
	}

	user := &models.User{
		ID:       userID,
		Username: fmt.Sprintf("%v", claims[username]),
		Email:    fmt.Sprintf("%v", claims[email]),
	}

	return user, true, nil
}
