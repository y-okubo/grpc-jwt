package user

import (
	"context"
	"log"

	"github.com/dgrijalva/jwt-go"
	"github.com/y-okubo/grpc-jwt/auth"
	"google.golang.org/grpc"
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
	return authenticateRuby()

	// // User information can be included in token.
	// c := User{
	// 	Account:  "y-okubo",
	// 	FullName: "Yuki Okubo",
	// }

	// c.Id = uuid.NewV4().String()
	// c.IssuedAt = time.Now().UTC().Unix()
	// c.NotBefore = time.Now().UTC().Unix()
	// c.ExpiresAt = time.Now().UTC().Unix()
	// c.Issuer = "Y.O"
	// c.Audience = "Y.O"
	// c.Subject = "AccessToken"

	// token := jwt.NewWithClaims(jwt.SigningMethodRS512, &c)

	// bin, err := ioutil.ReadFile("rsa")
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// key, err := jwt.ParseRSAPrivateKeyFromPEM(bin)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// // Convert JSON to a string (JWT) using a secret key.
	// str, err := token.SignedString(key)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// return &str
}

func authenticateRuby() *string {
	conn, err := grpc.Dial(":7831", grpc.WithInsecure())
	if err != nil {
		log.Println(err)
		return nil
	}
	defer conn.Close()

	c := auth.NewAuthenticatorClient(conn)

	resp, err := c.DoAuth(context.Background(), &auth.AuthRequest{})
	if err != nil {
		log.Println(err)
		return nil
	}

	return &resp.Token
}
