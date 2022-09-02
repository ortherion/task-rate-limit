package http

import (
	"context"
	"gitlab.com/g6834/team17/task-service/internal/constants"
	"gitlab.com/g6834/team17/task-service/internal/domain/models"
	"net/http"
)

type ctxKey int

const ridKey = ctxKey(0)

const (
	errInvalidToken = "invalid token"
)

func (s *Server) ValidateAuth() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			accessToken, err := r.Cookie(constants.ACCESS_TOKEN)
			if err != nil {
				http.Error(w, errInvalidToken, http.StatusUnauthorized)
				s.logger.Error().Msg(err.Error())
				return
			}
			refreshToken, err := r.Cookie(constants.REFRESH_TOKEN)
			if err != nil {
				http.Error(w, errInvalidToken, http.StatusUnauthorized)
				s.logger.Error().Msg(err.Error())
				return
			}

			tokens, err := s.auth.Validate(r.Context(), &models.TokenPair{
				Access:  accessToken.Value,
				Refresh: refreshToken.Value,
			})
			if err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				s.logger.Error().Msg(err.Error())
				return
			}

			user, _, _ := s.auth.ParseToken(r.Context(), tokens.Access)

			ctx := context.WithValue(r.Context(), constants.CTX_USER, user)

			atCookie := http.Cookie{
				Name:  constants.ACCESS_TOKEN,
				Value: tokens.Access,
				Path:  "/",
			}
			rtCookie := http.Cookie{
				Name:     constants.REFRESH_TOKEN,
				Value:    tokens.Refresh,
				Path:     "/",
				HttpOnly: true,
			}

			http.SetCookie(w, &atCookie)
			http.SetCookie(w, &rtCookie)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetReqID(ctx context.Context) string {
	return ctx.Value(ridKey).(string)
}
