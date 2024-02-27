package helpers

import (
	//"fmt"
	"log"
	"os"
	"time"

	"github.com/gbubemi22/moni/src/database"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

type SignedDetails struct {
	Email      string
	First_name string
	Last_name  string
	Uid        string
	User_type  string
	jwt.RegisteredClaims
}

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

var SECRET_KEY string = os.Getenv("SECRET_KEY")

func GenerateAllTokens(email string, firstName string, lastName string, uid string, user_type string) (signedToken string, err error) {
	claims := &SignedDetails{
		Email:       email,
		First_name:  firstName,
		Last_name:   lastName,
		Uid:         uid,
		User_type:   user_type,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(time.Hour * time.Duration(24))),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panic(err)
		return
	}

	return token, err
}



// func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {

// 	token, err := jwt.ParseWithClaims(
// 		signedToken,
// 		&SignedDetails{},
// 		func(token *jwt.Token) (interface{}, error) {
// 			return []byte(SECRET_KEY), nil
// 		},
// 	)

// 	//the token is invalid

// 	claims, ok := token.Claims.(*SignedDetails)
// 	if !ok {
// 		msg = fmt.Sprintf("the token is invalid")
// 		msg = err.Error()
// 		return
// 	}

// 	//the token is expired
// 	if claims.ExpiresAt < time.Now().Local().Unix() {
// 		msg = fmt.Sprint("token is expired")
// 		msg = err.Error()
// 		return
// 	}

// 	return claims, msg

// }