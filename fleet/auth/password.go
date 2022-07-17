package auth

import (
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/common"
	"golang.org/x/crypto/bcrypt"
)

func GeneratePassword(password string) (common.Password, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}
