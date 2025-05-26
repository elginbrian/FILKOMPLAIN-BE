package app

import (
	"errors"
	"mime/multipart"

	"github.com/elginbrian/FILKOMPLAIN-BE/config"
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

func (s *UserService) UpdateProfileFull(userID uint, newPassword, username, nim, programStudi, phoneNumber string) error {
	user, err := s.Repo.GetByID(userID)
	if err != nil {
		return err
	}

	// Update username if provided and different from current
	if username != "" && username != user.Username {
		// Check if username is already taken
		existingUser, err := s.Repo.GetByUsername(username)
		if err == nil && existingUser.ID != userID {
			return errors.New("username already exists")
		}
		user.Username = username
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

// Add new method for updating profile with image
func (s *UserService) UpdateProfileWithImage(userID uint, newPassword, username, nim, programStudi, phoneNumber string, profileImage *multipart.FileHeader) error {
	user, err := s.Repo.GetByID(userID)
	if err != nil {
		return err
	}

	// Update username if provided and different from current
	if username != "" && username != user.Username {
		// Check if username is already taken
		existingUser, err := s.Repo.GetByUsername(username)
		if err == nil && existingUser.ID != userID {
			return errors.New("username already exists")
		}
		user.Username = username
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

	// Upload profile image if provided
	if profileImage != nil {
		imageURL, err := s.UploadProfileImage(profileImage)
		if err != nil {
			return err
		}
		user.ProfileImageURL = imageURL
	}

	return s.Repo.Update(user)
}

func (s *UserService) UploadProfileImage(file *multipart.FileHeader) (string, error) {
	supabaseStorage := config.NewSupabaseStorage()

	fileURL, err := supabaseStorage.UploadFile(file)
	if err != nil {
		return "", err
	}

	return fileURL, nil
}
