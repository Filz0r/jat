package auth

import (
	"crypto/rand"

	"encoding/hex"
)

func MakeRefreshToken() string {
	key := make([]byte, 32)
	rand.Read(key)
	encoded := hex.EncodeToString(key)
	return encoded
}

//func ValidateRefreshToken(tokenString string, db *database.Queries, ctx context.Context) (database.RefreshToken, bool) {
//	token, err := db.GetRefreshToken(ctx, tokenString)
//	if err != nil {
//		fmt.Printf("Error validating refresh token: %v", err)
//		return database.RefreshToken{}, false
//	}
//
//	if token.RevokedAt.Valid {
//		return database.RefreshToken{}, false
//	}
//	if time.Now().After(token.ExpiresAt) {
//		return database.RefreshToken{}, false
//	}
//
//	return token, true
//}
