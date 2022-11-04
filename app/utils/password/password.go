package password

import (
	"crypto/sha512"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"m/v2/app/utils/password/argon2id"
)

// Generate return a hashed password
func Generate(raw string) string {
	argon2ID := argon2id.New()
	hash, err := argon2ID.GenerateFromPassword(raw)

	if err != nil {
		panic(err)
	}

	return string(hash)
}

// Verify compares a hashed password with plaintext password
func Verify(hash string, raw string) error {
	switch {
	case strings.HasPrefix(hash, "$2a$"),
		strings.HasPrefix(hash, "$2x$"),
		strings.HasPrefix(hash, "$2y$"),
		strings.HasPrefix(hash, "$2b$"):
		// Deprecated: Bcrypt with sha384 hash is deprecated.
		// TODO: remove Bcrypt from verify after users migrated to argon2id hash
		return bcrypt.CompareHashAndPassword([]byte(hash), sha384Sum([]byte(raw)))
	case strings.HasPrefix(hash, "$argon2id$"):
		argon2ID := argon2id.New()
		return argon2ID.ComparePasswordAndHash(hash, raw)
	default:
		return errors.New("password/verfiy: unknown hash type")
	}

}

func sha384Sum(raw []byte) []byte {
	sha384 := sha512.New384()
	sha384.Write(raw)
	return sha384.Sum(nil)
}
