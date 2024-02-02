package auth

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

type Service interface {
	GenerateToken(userId int) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

type jwtservice struct {
}

func Newjwtservice() *jwtservice {
	return &jwtservice{}
}

var SECRET_KEY = []byte("makan_malam")

func (j *jwtservice) GenerateToken(userId int) (string, error) {

	// var hmacSampleSecret []byte

	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	// 	"foo": "bar",
	// 	"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	// })

	// Sign and get the complete encoded token as a string using the secret
	// tokenString, err := token.SignedString(hmacSampleSecret)

	// fmt.Println(token, "fff", tokenString, "ff", err)
	// SECRET_KEY = []byte("makan_malam")
	claim := jwt.MapClaims{}

	claim["user_id"] = userId
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	signedToken, err := token.SignedString(SECRET_KEY)
	fmt.Println(token, "fff", signedToken, "ff", err)
	if err != nil {
		return signedToken, err
	}
	return signedToken, nil

}

func (s *jwtservice) ValidateToken(etoken string) (*jwt.Token, error) {
	var ccc []byte // just test this line
	fmt.Println(SECRET_KEY, "scret key atas validate")
	validatetoken, err := jwt.Parse(etoken, func(token *jwt.Token) (interface{}, error) {
		fmt.Println(token, "token inside atas", etoken)
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		fmt.Println(token, "token inside", etoken)
		if !ok {
			return nil, errors.New("invalid token")
		}
		fmt.Println(SECRET_KEY, "makan_malam", []byte(SECRET_KEY), ccc)
		return SECRET_KEY, nil

	})
	fmt.Println(validatetoken, "validation tooo")
	if err != nil {
		return validatetoken, err
	}

	return validatetoken, nil
}
