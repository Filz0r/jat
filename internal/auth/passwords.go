package auth

import "github.com/alexedwards/argon2id"

func HashPassword(password string) (string, error) {
	res, err := argon2id.CreateHash(password, argon2id.DefaultParams)

	return res, err
}

func CheckPasswordHash(password, hash string) (bool, error) {
	res, err := argon2id.ComparePasswordAndHash(password, hash)
	return res, err
}
