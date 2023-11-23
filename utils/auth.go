package utils

import (
	"Qpay/config"
	"Qpay/models"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"gorm.io/gorm"
)

type Credential struct {
	ID             int       `json:"id"`
	ExpirationTime time.Time `json:"expired_at"`
}

func newCredential(userID int, duration time.Duration) *Credential {
	fmt.Printf("TIME IN CONFIG: %v\n", duration*600000000)
	cred := &Credential{
		ID:             userID,
		ExpirationTime: time.Now().Add(duration * 600000000),
	}
	return cred
}

func (credential *Credential) Valid() error {
	if time.Now().After(credential.ExpirationTime) {
		return errors.New("token expired")
	}
	return nil
}

func CreateToken(config *config.JWT, userID int) (string, error) {
	cred := newCredential(userID, config.ExpirationTime)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cred)

	return token.SignedString([]byte(config.SecretKey))
}

func GetUser(db *gorm.DB, email, phoneNumber, password string) (*models.User, error) {
	var dbUser models.User

	if email != "" {
		result := db.First(&dbUser, "Email= ?", email)
		if result.RowsAffected == 0 {
			return nil, errors.New("User not found!")
		}
	} else if phoneNumber != "" && email == "" {
		result := db.First(&dbUser, "phone_number= ?", email)
		if result.RowsAffected == 0 {
			return nil, errors.New("User not found!")
		}
	}

	passChecker := CheckPassword(password, dbUser.Password)

	if passChecker {
		return &dbUser, nil
	}

	return nil, errors.New("Password not correct")
}

var ErrInvalidToken = errors.New("invalid token")

func VerifyToken(cfg *config.JWT, token string) (*Credential, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New(" ")
		}
		return []byte(cfg.SecretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Credential{}, keyFunc)

	if err != nil {
		return nil, err
	}

	cred, ok := jwtToken.Claims.(*Credential)

	if !ok {
		return nil, ErrInvalidToken
	}

	return cred, nil
}
