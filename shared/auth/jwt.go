package auth

import (
	"encoding/json"
	"errors"
	"github.com/tpp/msf/shared/constant"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/tpp/msf/config"
	"github.com/tpp/msf/model"
)

// UserAuth user auth
type UserAuth = model.User

type frClaims struct {
	jwt.StandardClaims
	Context *UserAuth `json:"context"`
}
type ID struct {
	ID int64 `json:"context"`
}

var (
	errInvalidToken = errors.New("invalid token")
)

// ParseJWTToken to userID
func ParseJWTToken(tokenString string) (Id int64, err error) {
	// get private key.
	key := config.GetConfigByte("auth.secret")
	if len(key) == 0 {
		key = constant.Secret
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errInvalidToken
		}
		return key, nil
	})
	if err != nil {
		return 0, errInvalidToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var claim frClaims
		var id ID
		bs, err := json.Marshal(claims)
		if err != nil {
			return 0, errInvalidToken
		}
		json.Unmarshal(bs, &claim)
		if err != nil {
			return 0, errInvalidToken
		}
		if claim.Context.ID == 0 {
			json.Unmarshal(bs, &id)
			return id.ID, nil
		}

		return int64(claim.Context.ID), nil
	}

	return 0, errInvalidToken

}

// GenerateJWTToken Generate token
func GenerateJWTToken(object any) (string, error) {

	exp := config.GetConfig[int64]("auth.claim.expire_in")
	issuer := config.GetConfig[string]("auth.claim.issuer")
	customClaims := object

	standardClaims := jwt.StandardClaims{
		Issuer:    issuer,
		ExpiresAt: time.Now().Add(time.Duration(exp) * time.Second).Unix(),
		IssuedAt:  time.Now().Unix(),
		NotBefore: time.Now().Unix(),
	}

	pay := payload{
		StandardClaims: standardClaims,
		Context:        customClaims,
	}

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = pay
	key := config.GetConfigByte("auth.secret")

	return token.SignedString(key)
}

type payload struct {
	jwt.StandardClaims
	Context any `json:"context"`
}
