package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "jat-http-access",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   userID.String(),
	})
	signed, err := tok.SignedString([]byte(tokenSecret))
	return signed, err
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	keyFunc := func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(tokenSecret), nil
	}

	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, keyFunc)

	if err != nil {
		return uuid.UUID{}, fmt.Errorf("error parsing jwt token: %w", err)
	}

	id, err := token.Claims.GetSubject()

	if err != nil {
		return uuid.UUID{}, fmt.Errorf("error getting subject from the token: %w", err)
	}

	parsed, err := uuid.Parse(id)

	if err != nil {
		return uuid.UUID{}, fmt.Errorf("error parsing uuid: %w", err)
	}

	return parsed, err
}
