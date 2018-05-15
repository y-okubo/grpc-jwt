package user

import (
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
)

// PrivateKey uses JWT signing
const PrivateKey = "BlPeUYZIibC6FGUjeWo5HT3RYmnkr4C99Q4XTqElCjqMIa1d8ei6QXzhY1m4dS2S" // 512-bit

// User is an user informations.
type User struct {
	Account  string `json:"account"`
	FullName string `json:"full_name"`
	jwt.StandardClaims
}

// Authenticate creates JWT string.
func Authenticate(u, p string) *string {
	// User information can be included in token.
	c := User{
		Account:  "y-okubo",
		FullName: "Yuki Okubo",
	}

	c.Id = uuid.NewV4().String()
	c.IssuedAt = time.Now().UTC().Unix()
	c.NotBefore = time.Now().UTC().Unix()
	c.ExpiresAt = time.Now().UTC().Unix()
	c.Issuer = "Y.O"
	c.Audience = "Y.O"
	c.Subject = "AccessToken"

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), &c)

	// Convert JSON to a string (JWT) using a secret key.
	str, err := token.SignedString([]byte(PrivateKey))
	if err != nil {
		log.Fatalln(err)
	}
	return &str
}
