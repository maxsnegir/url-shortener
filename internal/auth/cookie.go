package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
)

const AuthorizationCookieName = "AuthToken"

type CookieAuthentication struct {
	key    [sha256.Size]byte
	aesgcm cipher.AEAD
	nonce  []byte
}

func (a *CookieAuthentication) CreateToken() (string, string) {
	userToken := uuid.New().String()
	encodedUserToken := fmt.Sprintf("%x", a.aesgcm.Seal(nil, a.nonce, []byte(userToken), nil))
	return userToken, encodedUserToken
}

func (a *CookieAuthentication) ParseToken(token string) (string, error) {
	decodedToken, err := hex.DecodeString(token)
	if err != nil {
		return "", nil
	}
	userToken, err := a.aesgcm.Open(nil, a.nonce, decodedToken, nil)
	if err != nil {
		return "", err
	}
	return string(userToken), nil
}

func NewCookieAuthentication(secretKey string) (CookieAuthentication, error) {
	var cockieAuth CookieAuthentication
	key := sha256.Sum256([]byte(secretKey))
	aesblock, err := aes.NewCipher(key[:])
	if err != nil {
		return cockieAuth, err
	}
	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return cockieAuth, err
	}
	nonce := key[len(key)-aesgcm.NonceSize():]
	cockieAuth.key = key
	cockieAuth.aesgcm = aesgcm
	cockieAuth.nonce = nonce
	return cockieAuth, nil
}
