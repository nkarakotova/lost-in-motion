package postgreSQLObjectMother

import (
	"lim/internal/lim-core/models"
)

func CreateTestCoach() *models.Coach {
	return &models.Coach{
		ID:   1,
		Name: "Name",
	}
}
