package lib

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JWTParser struct {
	hmacSecret []byte
}

func NewJWTParser(secret string) JWTParser {
	sig := hmac.New(sha256.New, []byte(secret))
	hmacSecret := sig.Sum(nil)
	return JWTParser{hmacSecret: hmacSecret}
}

func (p *JWTParser) CreateToken(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(p.hmacSecret)
}

func (p *JWTParser) ParseToken(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return p.hmacSecret, nil
	})

	if token == nil || !token.Valid {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		t, expiredErr := claims.GetExpirationTime()
		if expiredErr != nil {
			return nil, expiredErr
		}
		if t == nil {
			return &claims, fmt.Errorf("no expiration time found")
		}
		if t.Time.Before(time.Now()) {
			return nil, fmt.Errorf("token expired")
		}

		return &claims, nil
	}
	return nil, err
}
