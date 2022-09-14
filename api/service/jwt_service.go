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
	IdUser int `json:"id_user"`
	jwt.StandardClaims
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AtExpires    int64
	RtExpires    int64
}

type JwtServiceImpl struct {
	accessSecretKey  string
	refreshSecretKey string
	jwtSigningMethod jwt.SigningMethod
}

func NewJwtService() JwtService {
	return &JwtServiceImpl{
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

func (service *JwtServiceImpl) GenerateToken(id int, expirationTime time.Time) (*TokenDetails, error) {
	tokens := &TokenDetails{}
	tokens.AtExpires = expirationTime.Unix()
	expiredTimeRefresh, _ := strconv.Atoi(os.Getenv("JWT_REFRESH_EXPIRED_TIME"))
	tokens.RtExpires = time.Now().Add(time.Duration(expiredTimeRefresh) * 24 * time.Hour).Unix()

	// Access token
	accessToken := jwtCustomClaim{
		IdUser: id,
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
		IdUser: id,
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

func (service *JwtServiceImpl) RefreshToken(refreshToken string) (*TokenDetails, error) {
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

	// Get id user
	id := int(claims["id_user"].(float64))

	// Delete the previous Refresh Token
	// --

	// Create new pairs of refresh and access tokens
	expiredTimeAccess, _ := strconv.Atoi(os.Getenv("JWT_ACCESS_EXPIRED_TIME"))
	tokens, err := service.GenerateToken(id, time.Now().Add(time.Duration(expiredTimeAccess)*24*time.Hour))
	if err != nil {
		return nil, err
	}

	// Save the token
	// --

	return tokens, nil
}

func (service *JwtServiceImpl) ValidateToken(token string) (*jwt.Token, error) {
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
