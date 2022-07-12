package health

import (
	"github.com/tgs266/fleet/fleet/persistence"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/health"
)

type HealthService struct{}

func (ss *HealthService) Health() (health.Health, error) {
	if !persistence.Ping() {
		return health.Health{
			Status: "UNAVAILABLE",
			Reason: "db connection failed",
		}, nil
	}
	return health.Health{
		Status: "AVAILABLE",
	}, nil
}
