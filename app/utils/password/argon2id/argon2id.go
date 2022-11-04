package argon2id

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

var (
	// ErrInvalidHash in returned by ComparePasswordAndHash if the provided
	// hash isn't in the expected format.
	ErrInvalidHash = errors.New("argon2id: hash is not in the correct format")

	// ErrIncompatibleVariant is returned by ComparePasswordAndHash if the
	// provided hash was created using a unsupported variant of Argon2.
	// Currently only argon2id is supported by this package.
	ErrIncompatibleVariant = errors.New("argon2id: incompatible variant of argon2")

	// ErrIncompatibleVersion is returned by ComparePasswordAndHash if the
	// provided hash was created using a different version of Argon2.
	ErrIncompatibleVersion = errors.New("argon2id: incompatible version of argon2")

	// ErrHashDoesNotMatch is returned by ComparePasswordAndHash if the
	// provided hash does not match the plain text password.
	ErrHashDoesNotMatch = errors.New("argon2id: hashedPassword is not the hash of the given password")
)

type Argon2ID struct {
	format  string
	version int
	time    uint32
	memory  uint32
	keyLen  uint32
	saltLen uint32
	threads uint8
}

// return a new Argon2ID struct with pre-selected paramaters
//
//	Paramaters: time: 4, memory: 64MiB, parallelism: 1, saltLen: 128bit, keyLen: 256bit
func New() Argon2ID {
	// TODO: maintain these params per:
	// https://datatracker.ietf.org/doc/html/draft-irtf-cfrg-argon2-03#section-4
	// Summary: Argon2id parameters should change over time
	// as resources allow for reasonable time execution, target
	// between 500ms and 1000ms under normal load on server.
	//
	// FYI: Current deployment taget is B2 Azure App Serivce Instance
	return Argon2ID{
		format:  "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		version: argon2.Version,
		time:    3,
		memory:  64 * 1024,
		keyLen:  32,
		saltLen: 16,
		threads: 2,
	}
}

// GenerateFromPassword returns the Argon2ID hash of the password.
// Use CompareHashAndPassword, as defined in this package,
// to compare the returned hashed password with its cleartext version.
func (a Argon2ID) GenerateFromPassword(plain string) (string, error) {
	salt := make([]byte, a.saltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(plain), salt, a.time, a.memory, a.threads, a.keyLen)

	return fmt.Sprintf(
			a.format,
			a.version,
			a.memory,
			a.time,
			a.threads,
			base64.RawStdEncoding.EncodeToString(salt),
			base64.RawStdEncoding.EncodeToString(hash),
		),
		nil
}

// CompareHashAndPassword compares a Argon2ID hashed password with its possible
// plaintext equivalent. Returns nil on success, or an error on failure.
func (a Argon2ID) ComparePasswordAndHash(hash, plain string) error {
	hashParts := strings.Split(hash, "$")
	if len(hashParts) != 6 {
		return ErrInvalidHash
	}

	if hashParts[1] != "argon2id" {
		return ErrIncompatibleVariant
	}

	var version int
	_, err := fmt.Sscanf(hashParts[2], "v=%d", &version)
	if err != nil || version != argon2.Version {
		return ErrIncompatibleVersion
	}

	_, err = fmt.Sscanf(hashParts[3], "m=%d,t=%d,p=%d", &a.memory, &a.time, &a.threads)
	if err != nil {
		return err
	}

	salt, err := base64.RawStdEncoding.DecodeString(hashParts[4])
	if err != nil {
		return err
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(hashParts[5])
	if err != nil {
		return err
	}

	hashToCompare := argon2.IDKey([]byte(plain), salt, a.time, a.memory, a.threads, uint32(len(decodedHash)))

	if subtle.ConstantTimeCompare(decodedHash, hashToCompare) == 1 {
		return nil
	} else {
		return ErrHashDoesNotMatch
	}
}
