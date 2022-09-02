package grpc

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/rs/zerolog"
	"gitlab.com/g6834/team17/task-service/internal/config"
	"gitlab.com/g6834/team17/task-service/internal/domain/models"
	pb "gitlab.com/g6834/team17/task-service/pkg/auth_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"strconv"
)

const (
	authorized = "authorized"
	userId     = "user_id"
	expired    = "expired"
	email      = "email"
	firstName  = "first_name"
	lastName   = "last_name"
	username   = "username"
)

type AuthClient struct {
	pb.AuthServiceClient
	conn      *grpc.ClientConn
	secretKey string
	logger    *zerolog.Logger
}

func NewAuthClient(cfg *config.Config, l *zerolog.Logger) (*AuthClient, error) {
	conn, err := grpc.Dial(net.JoinHostPort(cfg.Auth.Host, strconv.Itoa(cfg.Auth.Port)), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
		//t.Fatalf("grpc.Dial err:%v", err)
	}
	authClient := pb.NewAuthServiceClient(conn)

	return &AuthClient{
		conn:              conn,
		AuthServiceClient: authClient,
		secretKey:         cfg.Auth.SecretKey,
		logger:            l,
	}, nil
}

func (c *AuthClient) Validate(ctx context.Context, token *models.TokenPair) (*models.TokenPair, error) {
	v, err := c.AuthServiceClient.Validate(ctx, &pb.ValidateTokenRequest{
		AccessToken:  token.Access,
		RefreshToken: token.Refresh,
	})
	if err != nil {
		return nil, err
	}

	switch v.GetStatus() {
	case pb.Statuses_valid:
		return &models.TokenPair{
			Access:  v.GetAccessToken(),
			Refresh: v.GetRefreshToken(),
		}, nil
	case pb.Statuses_invalid:
		return nil, err
	}
	return nil, err
}

func (c *AuthClient) ParseToken(ctx context.Context, tokenString string) (*models.User, bool, error) {
	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(c.secretKey), nil
	})
	if err != nil || !token.Valid {
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

func (c *AuthClient) Close() {
	err := c.conn.Close()
	if err != nil {
		c.logger.Error().Err(err)
	}
}
