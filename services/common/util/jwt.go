package util

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var JWTSecret = "smart-kid"

type JWT struct {
	SigningKey []byte
}

var (
	TokenExpired     error = errors.New("Token is expired")
	TokenNotValidYet error = errors.New("Token not active yet")
	TokenMalformed   error = errors.New("That's not even a token")
	TokenInvalid     error = errors.New("Couldn't handle this token:")
)

type CustomClaims struct {
	UserId string `json:"user_id"`
	Email  string `json:"email"`
	jwt.StandardClaims
}

type TokenInfo struct {
	Token     string    `json:"token"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewJWT() *JWT {
	return &JWT{
		[]byte(JWTSecret),
	}
}

func (j *JWT) GenerateToken(UserId, Email string) (TokenInfo, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(30 * 24 * time.Hour)

	claims := CustomClaims{
		UserId,
		Email,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "SmartKid-Education",
		},
	}
	token, err := j.CreateToken(claims)
	return TokenInfo{
		Token:     token,
		ExpiredAt: expireTime,
	}, err
}

func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid
	}
	return nil, TokenInvalid
}

func (j *JWT) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.CreateToken(*claims)
	}
	return "", TokenInvalid
}
