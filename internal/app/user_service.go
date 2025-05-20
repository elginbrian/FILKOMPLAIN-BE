package app

import (
	"errors"

	"github.com/elginbrian/FILKOMPLAIN-BE/internal/model"
	"github.com/elginbrian/FILKOMPLAIN-BE/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

func (s *UserService) GetProfile(username string) (*model.User, error) {
	return s.Repo.GetByUsername(username)
}

func (s *UserService) UpdateProfile(username, newPassword string) error {
	user, err := s.Repo.GetByUsername(username)
	if err != nil {
		return err
	}
	if newPassword == "" {
		return errors.New("password required")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	return s.Repo.Update(user)
}
