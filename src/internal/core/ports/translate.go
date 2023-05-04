package ports

import (
	repositories "project/guidemysteps/src/internal/adapters/repositories/open_street_map/models"
	"project/guidemysteps/src/internal/core/models"
)

type TranslateRepo interface {
	Translate(repositories.Steps, string, string, models.User) string
}
