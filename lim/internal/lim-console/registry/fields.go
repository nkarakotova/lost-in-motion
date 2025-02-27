package registry

import (
	"lim/internal/lim-core/services"
)

type AppServiceFields struct {
	ClientService   services.ClientService
	CoachService    services.CoachService
	HallService     services.HallService
	TrainingService services.TrainingService
}
