package users

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tgs266/fleet/fleet/auth"
	"github.com/tgs266/fleet/fleet/panicker"
	"github.com/tgs266/fleet/fleet/utils"
	"github.com/tgs266/fleet/frn"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/common"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/entities"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/users"
	"github.com/tgs266/rest-gen/runtime/authentication"
	"github.com/tgs266/rest-gen/runtime/errors"
	"go.mongodb.org/mongo-driver/bson"
)

type UserService struct {
}

func (ss *UserService) ListUsers(ctx *gin.Context, offset *int64, pageSize *int64, sort string, token authentication.Token) (users.ListUsersResponse, error) {
	panicker.AndPanic(auth.CanI(token, "user", "*", auth.ACTION_VIEW)).GetOrPanicFunc(func(err error) error {
		panic(err)
	})
	total, err := Count(ctx)
	if err != nil {
		return users.ListUsersResponse{}, err
	}
	sortMap := bson.D{}
	if sort != "" {
		sortDir := -1 // descending
		if sort[0] == '+' {
			sortDir = 1
			sort = strings.TrimPrefix(sort, "+")
		} else if sort[0] == '-' {
			sort = strings.TrimPrefix(sort, "-")
		}
		sortMap = append(sortMap, bson.E{Key: sort, Value: sortDir})
	}

	res, err := ListUsers(ctx, utils.OrDefault(offset, 0), utils.OrDefault(pageSize, 0), sortMap)
	if err != nil {
		return users.ListUsersResponse{}, err
	}

	results := utils.ConvertList(res, func(input entities.UserEntity) users.User {
		return users.NewUserBuilder().
			SetFrn(input.Frn).
			SetName(input.Name).
			SetEmail(input.Email).
			SetCreatedAt(input.CreatedAt).
			SetModifiedAt(input.ModifiedAt).
			SetGroups(input.Groups).
			Build()
	})

	return users.ListUsersResponse{
		Items: results,
		Count: len(results),
		Total: int(total),
	}, nil
}

func (ss *UserService) CreateUser(ctx *gin.Context, body users.CreateUserRequest, token authentication.Token) (users.CreateUserResponse, error) {
	panicker.AndPanic(auth.CanI(token, "user", "*", auth.ACTION_CREATE)).GetOrPanicFunc(func(err error) error {
		panic(err)
	})
	now := time.Now()
	hashedPw := panicker.AndPanic(auth.GeneratePassword(body.Password)).GetOrPanicFunc(errors.NewInternalError)
	user := entities.NewUserEntityBuilder().
		SetFrn(common.UserFrn(frn.Generate("user", "default").String())).
		SetEmail(body.Email).
		SetName(body.Name).
		SetPassword(hashedPw).
		SetGroups([]string{}).
		SetCreatedAt(now).
		SetModifiedAt(now).Build()

	if err := CreateUser(ctx, user); err != nil {
		return users.CreateUserResponse{}, errors.NewInternalError(err)
	}
	return users.CreateUserResponse{Frn: user.Frn}, nil
}
