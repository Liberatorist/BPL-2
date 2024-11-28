package auth

import (
	"os"
	"time"

	"bpl/repository"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	UserID      int      `json:"user_id"`
	Permissions []string `json:"permissions"`
	Exp         int64    `json:"exp"`
}

func (claims *Claims) FromJWTClaims(jwtClaims jwt.Claims) {
	permissions := []string{}
	if jwtClaims.(jwt.MapClaims)["permissions"] != nil {
		for _, perm := range jwtClaims.(jwt.MapClaims)["permissions"].([]interface{}) {
			permissions = append(permissions, perm.(string))
		}
	}
	claims.Permissions = permissions
	claims.UserID = int(jwtClaims.(jwt.MapClaims)["user_id"].(float64))
	claims.Exp = int64(jwtClaims.(jwt.MapClaims)["exp"].(float64))
}

func (claims *Claims) Valid() error {
	if time.Now().Unix() > claims.Exp {
		return jwt.ErrTokenExpired
	}
	return nil
}

func CreateToken(user *repository.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id":     user.ID,
			"permissions": user.Permissions,
			"exp":         time.Now().Add(time.Hour * 24 * 7).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}
	return token, nil
}
