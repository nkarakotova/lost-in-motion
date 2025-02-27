package postgreSQLObjectMother

import (
	"lim/internal/lim-core/models"
)

func CreateTestClient() *models.Client {
	return &models.Client{
		ID:        1,
		Name:      "Name",
		Telephone: "1234567890",
		Mail:      "mail@mail.ru",
		Password:  "123",
	}
}
