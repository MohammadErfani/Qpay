package utils

import (
	"Qpay/config"
	"Qpay/models"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"gorm.io/gorm"
)

type Credential struct {
	UserID         int       `json:"id"`
	ExpirationTime time.Time `json:"expired_at"`
	jwt.RegisteredClaims
}

func newCredential(user *models.User, duration time.Duration) *Credential {
	cred := &Credential{
		UserID:         user.ID,
		ExpirationTime: time.Now().Add(duration),
	}

	return cred
}

func CreateToken(config *config.JWT, user *models.User) (string, error) {
	cred := newCredential(user, config.ExpirationTime)

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

var secretKey = []byte("mysecretKey")

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
		return nil, ErrInvalidToken
	}

	cred, ok := jwtToken.Claims.(*Credential)

	if !ok {
		return nil, ErrInvalidToken
	}

	return cred, nil
}
