package services

import (
	"github.com/gcarrenho/guidemysteps/src/internal/core/models"
	"github.com/gcarrenho/guidemysteps/src/internal/core/ports"
)

type userSvc struct {
	userRepo ports.UserRepository
}

func NewUserSvc(userRepo ports.UserRepository) *userSvc {
	return &userSvc{
		userRepo: userRepo,
	}
}
func (u *userSvc) Get(email string) (*models.User, error) {
	return u.userRepo.Get(email)
}

func (u *userSvc) Create(user models.User) error {
	return u.userRepo.Create(user)
}

func (u *userSvc) Update(user models.User) error {
	return u.userRepo.Update(user)
}
