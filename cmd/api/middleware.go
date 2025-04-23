package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func (app *Application) AuthTokenAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.unauthorizedError(w, r, fmt.Errorf("authorization header is missing"))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			app.unauthorizedError(w, r, fmt.Errorf("authorization header is malformed"))
			return
		}

		token := parts[1]

		jwtToken, err := app.authenticator.ValidateToken(token)
		if err != nil {
			app.unauthorizedError(w, r, err)
			return
		}

		claims := jwtToken.Claims.(jwt.MapClaims)

		userId, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)
		if err != nil {
			app.unauthorizedError(w, r, err)
			return
		}

		ctx := r.Context()
		user, err := app.store.Users.GetById(ctx, userId)
		if err != nil {
			app.unauthorizedError(w, r, err)
			return
		}
		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
