package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/yokoshima228/sso/internal/domain/models"
	"time"
)

func NewToken(u models.User, app models.App, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = u.Id
	claims["email"] = u.Email
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["appId"] = app.Id

	tokenStr, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}
