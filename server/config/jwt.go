package config

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// custom jwt secret
var jwtSecret = [] byte("G8FzL1dP3Y9Xg7b2c9VqKfQmZyJ7XHb6wTdf8M/UGZk=")

func GenerateSecret(userId string,email string) 	(string,error){
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{
		"userId":userId,
		"email":email,
		"exp":time.Now().Add(time.Hour*24).Unix(),
	})
	return token.SignedString(jwtSecret)

}

func ValidateToken(tokenString string)(*jwt.Token ,error){
	token,err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token")
		}
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("unauthorized")
	}

	return token, nil
}