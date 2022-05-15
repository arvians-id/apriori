package service

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

type jwtCustomClaim struct {
	IdUser uint64 `json:"id_user"`
	jwt.StandardClaims
}

type JwtService interface {
	GenerateToken(IdUser uint64) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

type jwtService struct {
	secretKey        string
	jwtSigningMethod jwt.SigningMethod
	issuer           string
}

func NewJwtService() JwtService {
	return &jwtService{
		secretKey:        getSecretKey(),
		jwtSigningMethod: jwt.SigningMethodHS256,
		issuer:           "wids",
	}
}

func getSecretKey() string {
	return "asu"
}

func (service *jwtService) GenerateToken(IdUser uint64) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)

	claims := jwtCustomClaim{
		IdUser: IdUser,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(service.jwtSigningMethod, claims)
	tokenString, err := token.SignedString([]byte(service.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (service *jwtService) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		method, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			str := fmt.Sprintf("unexpected signing method %v", token.Header["alg"])
			return nil, errors.New(str)
		} else if method != service.jwtSigningMethod {
			str := fmt.Sprintf("unexpected signing method %v", token.Header["alg"])
			return nil, errors.New(str)
		}

		return []byte(service.secretKey), nil
	})
}
