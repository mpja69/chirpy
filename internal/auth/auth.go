package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrNoAuthHeaderIncluded = errors.New("no auth header included in request")
var ErrAuthHeaderMalformed = errors.New("auth header is malformed")

func GetAuthorizationBearer(request *http.Request) (string, error) {
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

func GetAuthorizationAPIKey(request *http.Request) (string, error) {
	header := request.Header.Get("Authorization")
	if header == "" {
		return "", ErrNoAuthHeaderIncluded
	}
	splitHeader := strings.Split(header, " ")
	if len(splitHeader) != 2 || splitHeader[0] != "ApiKey" {
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

func MakeJWT(userId int, tokenSecret []byte, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   strconv.Itoa(userId),
	})

	return token.SignedString(tokenSecret)
}

func MakeRefreshToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	refreshToken := hex.EncodeToString(b)
	return refreshToken, err
}
