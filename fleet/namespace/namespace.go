package namespace

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tgs266/fleet/fleet/auth"
	"github.com/tgs266/fleet/fleet/panicker"
	"github.com/tgs266/fleet/fleet/persistence"
	"github.com/tgs266/fleet/fleet/utils"
	"github.com/tgs266/fleet/fleet/validation"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/common"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/errors"
	"github.com/tgs266/rest-gen/runtime/authentication"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type NamespaceService struct {
}

func (ss *NamespaceService) CreateNamespace(ctx *gin.Context, body common.CreateNamespaceRequest, token authentication.Token) (common.Namespace, error) {
	panicker.AndPanic(auth.IsAdmin(token)).GetOrPanicFunc(func(err error) error {
		panic(err)
	})
	panicker.OnFalse(validation.ValidateName(body.Namespace), errors.NewInvalidName(nil, body.Namespace))
	now := time.Now()
	entity := common.Namespace{
		Name:       body.Namespace,
		ModifiedAt: now,
		CreatedAt:  now,
	}

	if err := persistence.InsertOneToCollection(ctx, "namespaces", entity); err != nil {
		return common.Namespace{}, err
	}

	return entity, nil
}

func (ss *NamespaceService) ListNamespaces(ctx *gin.Context, sort string, token authentication.Token) ([]common.Namespace, error) {
	panicker.AndPanic(auth.IsAuth(token)).GetOrPanicFunc(func(err error) error {
		panic(err)
	})
	sortMap := utils.GetSortMap(sort)
	namespaces, err := persistence.List[common.Namespace](ctx, "namespaces", options.Find().SetSort(sortMap))
	if err != nil {
		return nil, err
	}

	return namespaces, nil
}
