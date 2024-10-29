package ports

import (
	repositories "github.com/gcarrenho/guidemysteps/src/internal/adapters/repositories/open_street_map/models"
	"github.com/gcarrenho/guidemysteps/src/internal/core/models"
)

type TranslateRepo interface {
	Translate(repositories.Steps, string, string, models.User) string
}
