package ports

import "github.com/gcarrenho/guidemysteps/src/internals/core/models"

type UserSvc interface {
	Get(email string) (*models.User, error)
	Create(user models.User) error
	Update(user models.User) error
}

type UserRepository interface {
	Get(email string) (*models.User, error)
	Create(user models.User) error
	Update(user models.User) error
}
