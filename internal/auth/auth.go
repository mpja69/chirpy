package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var ErrNoAuthHeaderIncluded = errors.New("no auth header included in request")
var ErrAuthHeaderMalformed = errors.New("auth header is malformed")

func GetBearerToken(request *http.Request, secret string) (string, error) {
	header := request.Header.Get("Authorization")
	if header == "" {
		return "", ErrNoAuthHeaderIncluded
	}
	splitHeader := strings.Split(header, " ")
	if len(splitHeader) != 2 || splitHeader[0] != "Bearer" {
		return "", ErrAuthHeaderMalformed
	}
	return splitHeader[1], nil
}

func ValidateJWT(tokenString string, tokenSecret []byte) (string, error) {

	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return tokenSecret, nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}
	userIdString, err := token.Claims.GetSubject()
	if err != nil {
		return "", err
	}
	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return "", err
	}

	if issuer == "" {
		return "", fmt.Errorf("invalid user")
	}
	return userIdString, err
}
