package jwt_token

import (
	"context"
	"fmt"
	"time"

	"log"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hilmiikhsan/simple-messaging-app/pkg/env"
	"go.elastic.co/apm"
)

type ClaimToken struct {
	Username string `json:"username"`
	Fullname string `json:"full_name"`
	jwt.RegisteredClaims
}

var MapTypeToken = map[string]time.Duration{
	"token":         time.Hour * 3,
	"refresh_token": time.Hour * 72,
}

var jwtSecret = []byte(env.GetEnv("APP_SECRET", ""))

func GenerateToken(ctx context.Context, username, fullname, tokenType string, now time.Time) (string, error) {
	span, _ := apm.StartSpan(ctx, "GenerateToken", "jwt")
	defer span.End()

	claimToken := ClaimToken{
		Username: username,
		Fullname: fullname,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    env.GetEnv("APP_NAME", ""),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(MapTypeToken[tokenType])),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimToken)

	resultToken, err := token.SignedString(jwtSecret)
	if err != nil {
		log.Println("failed to generate token: ", err)
		return resultToken, fmt.Errorf("failed to generate token: %v", err)
	}

	return resultToken, nil
}

func ValidateToken(ctx context.Context, token string) (*ClaimToken, error) {
	span, _ := apm.StartSpan(ctx, "ValidateToken", "jwt")
	defer span.End()

	var (
		claimToken *ClaimToken
		ok         bool
	)

	jwtToken, err := jwt.ParseWithClaims(token, &ClaimToken{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Println("failed to validate method jwt: ", t.Header["alg"])
			return nil, fmt.Errorf("failed to validate method jwt: %v", t.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		log.Println("failed to parse jwt: ", err)
		return nil, fmt.Errorf("failed to parse jwt: %v", err)
	}

	if claimToken, ok = jwtToken.Claims.(*ClaimToken); !ok || !jwtToken.Valid {
		log.Println("token invalid")
		return nil, fmt.Errorf("token invalid")
	}

	return claimToken, nil
}
