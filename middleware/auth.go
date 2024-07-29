package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/IsraeliPS/GoMongo/config"
	"github.com/dgrijalva/jwt-go"
)

func JWTAuthentication(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Authorization header is required", http.StatusUnauthorized)
            return
        }

        bearerToken := strings.Split(authHeader, " ")
        fmt.Println("bearerToken", bearerToken)
        if len(bearerToken) != 1 {
            http.Error(w, "Invalid authorization token", http.StatusUnauthorized)
            return
        }

        tokenString := bearerToken[1]
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }
            return []byte(config.LoadEnv().JWT_SECRET), nil
        })

        if err != nil || !token.Valid {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        // If the token is valid, proceed with the request
        ctx := context.WithValue(r.Context(), "user", token.Claims)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
