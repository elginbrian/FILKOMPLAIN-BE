package app

import (
	"errors"
	"os"
	"time"

	"github.com/elginbrian/FILKOMPLAIN-BE/internal/model"
	"github.com/elginbrian/FILKOMPLAIN-BE/internal/repository"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	Repo repository.AuthRepository
}

func NewAuthService(repo repository.AuthRepository) *AuthService {
	return &AuthService{Repo: repo}
}

func (s *AuthService) Register(username, email, password string) error {
	_, err := s.Repo.GetByUsername(username)
	if err == nil {
		return errors.New("username already exists")
	}
	
	_, err = s.Repo.GetByEmail(email)
	if err == nil {
		return errors.New("email already exists")
	}
	
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("could not hash password")
	}
	user := &model.User{
		Username: username,
		Email:    email,
		Password: string(hash),
		Type:     "user",
	}
	return s.Repo.Create(user)
}

func (s *AuthService) RegisterWithType(username, email, password, userType string) error {
	_, err := s.Repo.GetByUsername(username)
	if err == nil {
		return errors.New("username already exists")
	}
	
	_, err = s.Repo.GetByEmail(email)
	if err == nil {
		return errors.New("email already exists")
	}
	
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("could not hash password")
	}
	user := &model.User{
		Username: username,
		Email:    email,
		Password: string(hash),
		Type:     userType,
	}
	return s.Repo.Create(user)
}

func (s *AuthService) GenerateTokenForUser(user *model.User) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", errors.New("JWT_SECRET not configured")
	}
	
	claims := jwt.MapClaims{
		"username": user.Username,
		"type":     user.Type,
		"user_id":  user.ID,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", errors.New("could not generate token")
	}
	
	return signed, nil
}

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.Repo.GetByEmail(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}
	
	return s.GenerateTokenForUser(user)
}

func (s *AuthService) LoginWithType(email, password, userType string) (string, error) {
	user, err := s.Repo.GetByEmail(email)
	if err != nil || user.Type != userType {
		return "", errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}
	
	return s.GenerateTokenForUser(user)
}

func (s *AuthService) GetProfile(username string) (*model.User, error) {
	return s.Repo.GetByUsername(username)
}

func (s *AuthService) GetProfileByID(userID uint) (*model.User, error) {
	return s.Repo.GetByID(userID)
}

func (s *AuthService) UpdateProfile(username, newPassword string) error {
	user, err := s.Repo.GetByUsername(username)
	if err != nil {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	return s.Repo.Update(user)
}

func (s *AuthService) UpdateProfileFull(userID uint, newPassword, username, nim, programStudi, phoneNumber string) error {
	user, err := s.Repo.GetByID(userID)
	if err != nil {
		return err
	}
	
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
