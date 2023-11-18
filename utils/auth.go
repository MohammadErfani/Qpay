package utils

import (
	"Qpay/config"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

// strconv.FormatInt(time.Now().UTC().UnixNano(), 10) `json:"expired_at"`
type Credential struct {
	email          string    `json:"email"`
	ExpirationTime time.Time `json:"expired_at"`
	jwt.RegisteredClaims
}

func newCredential(email string, duration time.Duration) *Credential {
	cred := &Credential{
		email:          email,
		ExpirationTime: time.Now().Add(duration),
	}

	return cred
}

func CreateToken(config *config.JWT, email string) (string, error) {
	cred := newCredential(email, config.ExpirationTime)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cred)

	return token.SignedString([]byte(config.SecretKey))

}
