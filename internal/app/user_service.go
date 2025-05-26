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

func (s *UserService) GetProfileByID(userID uint) (*model.User, error) {
	// We need to access the AuthRepository to get user by ID
	// Since we don't have direct access, let's implement this in the UserRepository
	return s.Repo.GetByID(userID)
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

func (s *UserService) UpdateProfileFull(userID uint, newPassword, nim, programStudi, phoneNumber string) error {
	user, err := s.Repo.GetByID(userID)
	if err != nil {
		return err
	}

	// Update password if provided
	if newPassword != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hash)
	}

	// Update optional fields if provided
	if nim != "" {
		user.NIM = nim
	}

	if programStudi != "" {
		user.ProgramStudi = programStudi
	}

	if phoneNumber != "" {
		user.PhoneNumber = phoneNumber
	}

	return s.Repo.Update(user)
}
