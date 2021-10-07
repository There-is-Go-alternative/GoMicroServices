package domain

import (
	"fmt"
	"log"
	"os"
	"time"

	jwt "github.com/form3tech-oss/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var SecretKey = []byte(os.Getenv("SECRET_KEY"))

type Token struct {
	Token string `json:"token"`
}
type Authorize struct {
	UserID string `json:"user_id"`
}
type Auth struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

//id et clé d'encryption
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Fatal(err)
	}

	return string(bytes), err
}

func VerifyPassword(hashed, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
}

type Claims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

func CreateToken(userID string) (string, error) {
	expirationTime := time.Now().Add(5 * time.Hour)
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SecretKey)
}

func VerifyToken(tokenStr string) (string, error) {
	claims := new(Claims)
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", fmt.Errorf("token is invalid")
	}

	return claims.UserID, nil
}
