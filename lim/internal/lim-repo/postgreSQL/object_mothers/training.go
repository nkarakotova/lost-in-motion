package postgreSQLObjectMother

import (
	"lim/internal/lim-core/models"
	"time"
)

func CreateTestTraining() *models.Training {
	return &models.Training{
		ID:        1,
		CoachID:   1,
		HallID:    1,
		Name:      "Name",
		DateTime:  time.Date(2024, 7, 7, 12, 0, 0, 0, time.UTC),
		PlacesNum: 10,
	}
}
