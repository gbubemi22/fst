package user

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

// func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
// 	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
// 	check := true
// 	msg := ""

// 	if err != nil {
// 		msg = "email of password is incorrect"
// 		check = false
// 	}
// 	return check, msg
// }



func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(providedPassword))
	check := true
	msg := ""

	if err != nil {
		msg = "email or password is incorrect"
		check = false
	}
	return check, msg
}