package service

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"os"
	"strconv"
	"time"
)

type jwtCustomClaim struct {
	IdUser uint64 `json:"id_user"`
	jwt.StandardClaims
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AtExpires    int64
	RtExpires    int64
}
type JwtService interface {
	GenerateToken(IdUser uint64, expirationTime time.Time) (*TokenDetails, error)
	RefreshToken(refreshToken string) (*TokenDetails, error)
	ValidateToken(token string) (*jwt.Token, error)
}

type jwtService struct {
	accessSecretKey  string
	refreshSecretKey string
	jwtSigningMethod jwt.SigningMethod
}

func NewJwtService() JwtService {
	return &jwtService{
		accessSecretKey:  getAccessSecretKey(),
		refreshSecretKey: getRefreshSecretKey(),
		jwtSigningMethod: jwt.SigningMethodHS256,
	}
}

func getAccessSecretKey() string {
	return os.Getenv("JWT_SECRET_ACCESS_KEY")
}
func getRefreshSecretKey() string {
	return os.Getenv("JWT_SECRET_REFRESH_KEY")
}

func (service *jwtService) GenerateToken(IdUser uint64, expirationTime time.Time) (*TokenDetails, error) {
	tokens := &TokenDetails{}
	tokens.AtExpires = expirationTime.Unix()
	tokens.RtExpires = time.Now().Add(7 * 24 * time.Hour).Unix()

	// Access token
	accessToken := jwtCustomClaim{
		IdUser: IdUser,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokens.AtExpires,
		},
	}

	tokenAt := jwt.NewWithClaims(service.jwtSigningMethod, accessToken)
	signedAt, err := tokenAt.SignedString([]byte(service.accessSecretKey))
	if err != nil {
		return nil, err
	}
	tokens.AccessToken = signedAt

	// Refresh token
	refreshToken := jwtCustomClaim{
		IdUser: IdUser,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokens.RtExpires,
		},
	}

	tokenRt := jwt.NewWithClaims(service.jwtSigningMethod, refreshToken)
	signedRt, err := tokenRt.SignedString([]byte(service.refreshSecretKey))
	if err != nil {
		return nil, err
	}
	tokens.RefreshToken = signedRt

	return tokens, nil
}

func (service *jwtService) RefreshToken(refreshToken string) (*TokenDetails, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		method, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			str := fmt.Sprintf("unexpected signing method %v", token.Header["alg"])
			return nil, errors.New(str)
		} else if method != service.jwtSigningMethod {
			str := fmt.Sprintf("unexpected signing method %v", token.Header["alg"])
			return nil, errors.New(str)
		}

		return []byte(service.refreshSecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	// Check if token is expired
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Ge user id
	userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["id_user"]), 10, 64)
	if err != nil {
		return nil, err
	}

	// Delete the previous Refresh Token
	// --

	// Create new pairs of refresh and access tokens
	tokens, err := service.GenerateToken(userId, time.Now().Add(5*time.Minute))
	if err != nil {
		return nil, err
	}

	// Save the token
	// --

	return tokens, nil
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

		return []byte(service.accessSecretKey), nil
	})
}
