package util

import (
	"crypto/rand"
	"errors"
	"math/big"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	cfg "github.com/dwarvesf/yggdrasil/organization/cmd/config"
	"golang.org/x/crypto/bcrypt"
)

//HashPassword hash password to a hash string with cost is 14
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

//VerifyPassword check password is valid
func VerifyPassword(password, cipherText string) error {
	return bcrypt.CompareHashAndPassword([]byte(cipherText), []byte(password))
}

//AuthToken for user authenticaion
type AuthToken struct {
	Token string `json:"access_token"`
}

//AuthTokenClaim ...
type AuthTokenClaim struct {
	*jwt.StandardClaims
	AuthPayload
}

//AuthPayload ...
type AuthPayload struct {
	LoginType    string `json:"login_type"`
	organization string `json:"organization"`
}

//CreateAccessToken create access token for user
func CreateAccessToken(organization, loginType string) (string, error) {
	//Expires after 100 days
	expiresAt := time.Now().Add(time.Minute * 144000).Unix()

	token := jwt.New(jwt.SigningMethodHS256)

	token.Claims = &AuthTokenClaim{
		&jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
		AuthPayload{
			LoginType:    loginType,
			organization: organization,
		},
	}

	tokenString, err := token.SignedString([]byte(cfg.JwtSecret))

	if err != nil {
		return "", ErrLogin
	}

	return tokenString, nil
}

//VerifyToken ...
func VerifyToken(token string) error {
	token, err := getTokenString(token)
	if err == nil {
		token, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrorUnauthorize
			}
			return []byte(cfg.JwtSecret), nil
		})

		if err != nil {
			return ErrorUnauthorize
		}

		if token.Valid {
			return nil
		}
	}

	return ErrorUnauthorize
}

//GenerateToken for verify user
func GenerateToken() string {
	token := ""
	for i := 0; i < 6; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(10))
		token = token + strconv.Itoa(int(n.Int64()))
	}
	return token
}

func getTokenString(tokenString string) (string, error) {
	bearerToken := strings.Split(tokenString, " ")
	if len(bearerToken) != 2 {
		return "", errors.New("Token string invalid")
	}
	return bearerToken[1], nil
}
