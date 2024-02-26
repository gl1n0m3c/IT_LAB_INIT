package jwt

import (
	"github.com/gl1n0m3c/IT_LAB_INIT/pkg/config"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"time"
)

const (
	Public     = "public"
	Manager    = "manager"
	Specialist = "specialist"
)

type JWT interface {
	CreateToken(id int, userType string) string
	Authorize(tokenString string, access string) (userClaim, bool, error)
}

type JWTUtil struct {
	expireTimeOut time.Duration
	secret        string
}

func InitJWTUtil() JWT {
	return JWTUtil{
		expireTimeOut: time.Duration(viper.GetInt(config.JWTExpire)) * time.Minute,
		secret:        viper.GetString(config.JWTSecret),
	}
}

type userClaim struct {
	jwt.RegisteredClaims
	ID       int
	UserType string
}

func (j JWTUtil) CreateToken(id int, userType string) string {

	expiredAt := time.Now().Add(j.expireTimeOut)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{
				Time: expiredAt,
			},
		},
		ID:       id,
		UserType: userType,
	})

	signedString, _ := token.SignedString([]byte(j.secret))

	return signedString
}

func (j JWTUtil) Authorize(tokenString string, access string) (userClaim, bool, error) {
	var claim userClaim

	token, err := jwt.ParseWithClaims(tokenString, &claim, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})
	if err != nil {
		return userClaim{}, false, err
	}

	if !token.Valid {
		return userClaim{}, false, nil
	}

	switch access {
	case Manager:
		return claim, claim.UserType == Manager, nil
	case Specialist:
		return claim, claim.UserType == Specialist, nil
	default:
		panic("you are passing wrong access")

	}
}
