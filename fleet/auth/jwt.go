package auth

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/common"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/entities"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/errors"
	"github.com/tgs266/rest-gen/runtime/authentication"
)

var jwtKey = []byte("RANDOM GEN")

type JWTClaim struct {
	Email   string         `json:"email"`
	UserFrn common.UserFrn `json:"userFrn"`
	jwt.StandardClaims
}

func GenerateJwt(user entities.UserEntity, expirationTime time.Time) (authentication.Token, error) {
	claims := &JWTClaim{
		UserFrn: user.Frn,
		Email:   user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	return authentication.Token(tokenString), err
}

func DecodeToken(signedToken authentication.Token) (*JWTClaim, error) {
	token, err := jwt.ParseWithClaims(
		signedToken.String(),
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)
	if err != nil {
		return nil, errors.NewInvalidToken(nil, signedToken.String())
	}
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		return nil, errors.NewInvalidToken(nil, signedToken.String())
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, errors.NewExpiredToken(nil, signedToken.String())
	}
	return claims, nil
}
