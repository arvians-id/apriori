package auth

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"log"
	"os"
	"strconv"
	"time"
)

type jwtCustomClaim struct {
	IdUser   int `json:"id_user"`
	RoleUser int `json:"role"`
	jwt.StandardClaims
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AtExpires    int64
	RtExpires    int64
}

type JsonWebToken struct {
	accessSecretKey  string
	refreshSecretKey string
	jwtSigningMethod jwt.SigningMethod
}

func NewJsonWebToken() *JsonWebToken {
	return &JsonWebToken{
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

func (auth *JsonWebToken) GenerateToken(id int, role int, expirationTime time.Time) (*TokenDetails, error) {
	tokens := &TokenDetails{}
	tokens.AtExpires = expirationTime.Unix()
	expiredTimeRefresh, err := strconv.Atoi(os.Getenv("JWT_REFRESH_EXPIRED_TIME"))
	if err != nil {
		log.Println("[JWTService][GenerateToken] problem in conversion string to integer, err: ", err.Error())
		return nil, err
	}
	tokens.RtExpires = time.Now().Add(time.Duration(expiredTimeRefresh) * 24 * time.Hour).Unix()

	// Access token
	accessToken := jwtCustomClaim{
		IdUser:   id,
		RoleUser: role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokens.AtExpires,
		},
	}

	tokenAt := jwt.NewWithClaims(auth.jwtSigningMethod, accessToken)
	signedAt, err := tokenAt.SignedString([]byte(auth.accessSecretKey))
	if err != nil {
		log.Println("[JWTService][GenerateToken] problem in first signed string, err: ", err.Error())
		return nil, err
	}
	tokens.AccessToken = signedAt

	// Refresh token
	refreshToken := jwtCustomClaim{
		IdUser:   id,
		RoleUser: role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokens.RtExpires,
		},
	}

	tokenRt := jwt.NewWithClaims(auth.jwtSigningMethod, refreshToken)
	signedRt, err := tokenRt.SignedString([]byte(auth.refreshSecretKey))
	if err != nil {
		log.Println("[JWTService][GenerateToken] problem in second signed string, err: ", err.Error())
		return nil, err
	}
	tokens.RefreshToken = signedRt

	return tokens, nil
}

func (auth *JsonWebToken) RefreshToken(refreshToken string) (*TokenDetails, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		method, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			str := fmt.Sprintf("unexpected signing method %v", token.Header["alg"])
			log.Println("[JWTService][RefreshToken] problem in refresh token, err: ", str)
			return nil, errors.New(str)
		} else if method != auth.jwtSigningMethod {
			str := fmt.Sprintf("unexpected signing method %v", token.Header["alg"])
			log.Println("[JWTService][RefreshToken] problem in refresh token, err: ", str)
			return nil, errors.New(str)
		}
		return []byte(auth.refreshSecretKey), nil
	})
	if err != nil {
		log.Println("[JWTService][RefreshToken] problem in parsing token, err: ", err.Error())
		return nil, err
	}

	// Check if token is expired
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		log.Println("[JWTService][RefreshToken] problem in token expired, err: ", err.Error())
		return nil, errors.New("invalid token")
	}

	// Get id user
	id := int(claims["id_user"].(float64))
	role := int(claims["role"].(float64))

	// Delete the previous Refresh Token
	// --

	// Create new pairs of refresh and access tokens
	expiredTimeAccess, err := strconv.Atoi(os.Getenv("JWT_ACCESS_EXPIRED_TIME"))
	if err != nil {
		log.Println("[JWTService][RefreshToken] problem in conversion string to integer, err: ", err.Error())
		return nil, err
	}

	tokens, err := auth.GenerateToken(id, role, time.Now().Add(time.Duration(expiredTimeAccess)*24*time.Hour))
	if err != nil {
		log.Println("[JWTService][RefreshToken] problem in getting generate token, err: ", err.Error())
		return nil, err
	}

	// Save the token
	// --

	return tokens, nil
}

func (auth *JsonWebToken) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		method, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			str := fmt.Sprintf("unexpected signing method %v", token.Header["alg"])
			log.Println("[JWTService][ValidateToken] problem in validating token, err: ", str)
			return nil, errors.New(str)
		} else if method != auth.jwtSigningMethod {
			str := fmt.Sprintf("unexpected signing method %v", token.Header["alg"])
			log.Println("[JWTService][ValidateToken] problem in validating token, err: ", str)
			return nil, errors.New(str)
		}
		return []byte(auth.accessSecretKey), nil
	})
}
