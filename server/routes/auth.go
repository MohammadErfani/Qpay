package routes

import (
	"Qpay/utils"
	"encoding/json"
	// "fmt"
	"net/http"
	"time"

	echo "github.com/labstack/echo/v4"
)

type User struct {
	email       string `json:"email"`
	phoneNumber string `json:"phoneNumber"`
	password    string `json:"password"`
}

// Create a struct that will be encoded to a JWT.
// We add jwt.RegisteredClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	token string `json:"token"`
}

func GenerateToken(claims jwt.Claims, key interface{}, method jwt.SigningMethod) (string, error) {
	// Create a new token with the specified claims
	token := jwt.NewWithClaims(method, claims)

	// Sign and encode the token using the key
	return token.SignedStringkey
}

func AuthGroup(authG *echo.Group) {
	json_map := make(map[string]interface{})

	// user := new(User)
	authG.POST("/login", func(ctx echo.Context) error {

		err := json.NewDecoder(ctx.Request().Body).Decode(&json_map)

		email := json_map["email"].(string)
		emailError := utils.IsValidEmail(email)
		if emailError != nil {
			return ctx.JSON(http.StatusBadRequest, emailError)
		}

		phoneNumber := json_map["phoneNumber"].(string)
		phoneNumberError := utils.IsValidPhoneNumber(phoneNumber)
		if phoneNumberError != nil {
			return ctx.JSON(http.StatusBadRequest, phoneNumberError)
		}

		password := json_map["password"].(string)
		if len(password) == 0 {
			return ctx.JSON(http.StatusBadRequest, "Input your password")
		}

		if err != nil {
			return ctx.String(http.StatusBadRequest, "Bad request")
		}

		// create JWT token
		expirationTime := time.Now().Add(5 * time.Minute)
		// Create a new token object, specifying signing method and the claims
		// you would like it to contain.
		token := echoJwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"email":          email,
			"expirationTime": expirationTime,
		})

		return ctx.JSON(http.StatusOK, "askdjaskjdASDJK!J@kjASJDK")
	})
}
