package system

import (
	"context"
	"time"

	"github.com/tgs266/fleet/fleet/auth"
	"github.com/tgs266/fleet/fleet/panicker"
	"github.com/tgs266/fleet/fleet/utils"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/system"
	"github.com/tgs266/rest-gen/runtime/authentication"
	"github.com/tgs266/rest-gen/runtime/errors"
)

type SystemService struct {
}

func (ss *SystemService) ListEnvironments(sort string, token authentication.Token) ([]system.Environment, error) {
	panicker.AndPanic(auth.IsAuth(token)).GetOrPanicFunc(func(err error) error {
		panic(err)
	})
	sortMap := utils.GetSortMap(sort)

	envs := panicker.AndPanic(ListEnvs(context.TODO(), sortMap)).GetOrPanicFunc(errors.NewNotFound)
	return envs, nil
}

func (ss *SystemService) CreateEnvironment(body system.CreateEnvironmentRequest, token authentication.Token) (system.Environment, error) {
	panicker.AndPanic(auth.IsAdmin(token)).GetOrPanicFunc(func(err error) error {
		panic(err)
	})
	now := time.Now()
	env := system.Environment{
		CreatedAt:  now,
		ModifiedAt: now,
		Name:       body.Name,
		Image:      body.Image,
	}

	panicker.CheckAndPanicFunc(CreateEnv(context.TODO(), env), errors.NewInternalError)
	return env, nil
}
