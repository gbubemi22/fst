package helpers

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type SignedDetails struct {
	Email      string
	First_name string
	Last_name  string
	Uid        string
	User_type  string
	Phone      string
	jwt.RegisteredClaims
}

var SECRET_KEY string = os.Getenv("SECRET_KEY")

func GenerateAllTokens(email string, first_name string, last_name string, uid string, user_type string) (signedToken string, err error) {
	// Check for mandatory arguments
	if email == "" || first_name == "" || last_name == "" || uid == "" || user_type == "" {
	    return "", errors.New("missing required information to generate tokens")
	}
  
	// Ensure SECRET_KEY is set
	if SECRET_KEY == "" {
	    return "", errors.New("environment variable SECRET_KEY is not set")
	}
  
	claims := &SignedDetails{
	    Email:      email,
	    First_name: first_name,
	    Last_name:  last_name,
	    Uid:        uid,
	    User_type:  user_type,
	    RegisteredClaims: jwt.RegisteredClaims{
		  ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(time.Hour * 24)),
	    },
	}
  
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
	    log.Printf("Error generating token: %v", err)
	    return "", err
	}
  
	return token, nil
  }
func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {

	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	// Handle invalid token
	if !token.Valid {
		msg = fmt.Sprintf("the token is invalid: %v", err)
		return
	}

	// Convert claims to SignedDetails type
	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = fmt.Sprintf("failed to cast claims to SignedDetails")
		return
	}

	// Check for expired token using Unix timestamps
	expirationUnix := claims.ExpiresAt.Unix()
	if err != nil {
		msg = fmt.Sprint("failed to extract expiration timestamp:", err)
		return
	}

	if expirationUnix < time.Now().Local().Unix() {
		msg = fmt.Sprint("token is expired")
		return
	}

	return claims, msg
}
