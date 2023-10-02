package security

import (
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

const (
	// TokenExpireDuration token expire
	JWT_SECRET     = "secret"
	JWT_EXP_HOUR   = 1
	JWT_EXP_MINUTE = 0
	JWT_EXP_SECOND = 0
)


func NewToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"iat":   time.Now().Unix(),
		"nbf":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Hour*time.Duration(JWT_EXP_HOUR) + time.Minute*time.Duration(JWT_EXP_MINUTE) + time.Second*time.Duration(JWT_EXP_SECOND)).Unix(),
	})

	signedToken, err := token.SignedString([]byte(JWT_SECRET))
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return signedToken, nil
}

func ParseToken(tokenString string) (*jwt.Token, error) {
	var token *jwt.Token
	var err error

	token, err = parseHS256Token(tokenString, token)
	if err != nil && err.Error() != "Token is expired" {
		return nil, err
	}

	return token, nil
}

func parseHS256Token(tokenString string, token *jwt.Token) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok { // SigningMethodHMAC
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(JWT_SECRET), nil
	})
	return token, err
}

func GetClaims(token *jwt.Token) (jwt.MapClaims, error) {
	if !token.valid {
		return nil, fmt.Errorf("Unauthorized")
	}
	err := token.Claims.(jwt.MapClaims).Valid()
	if err != nil {
		return nil, err
	}

	return token.Claims.(jwt.MapClaims), nil
}
