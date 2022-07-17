package auth

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tgs266/fleet/fleet/panicker"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/auth"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/errors"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct{}

func (service *AuthService) Login(ctx *gin.Context, body auth.LoginRequest) (auth.LoginResponse, error) {
	user := panicker.Panicker(GetUserByEmail(body.Email)).
		GetOrPanicFunc(func(err error) error { return errors.NewLoginFailed(err, body.Email) })
	pwd := []byte(body.Password)
	if err := bcrypt.CompareHashAndPassword(user.Password, pwd); err != nil {
		return auth.LoginResponse{}, errors.NewLoginFailed(err, body.Email)
	}
	expiration := time.Now().Add(1 * time.Hour)
	jwt := panicker.Panicker(GenerateJwt(user, expiration)).GetOrPanic()
	ctx.SetCookie("F_TOKEN", jwt.String(), int(expiration.Unix()), "/", "", false, true)
	return auth.LoginResponse{
		Token: string(jwt),
	}, nil
}
