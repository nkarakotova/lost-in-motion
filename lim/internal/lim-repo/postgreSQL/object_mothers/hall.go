package postgreSQLObjectMother

import (
	"lim/internal/lim-core/models"
)

func CreateTestHall() *models.Hall {
	return &models.Hall{
		ID:     1,
		Number: 1,
	}
}
