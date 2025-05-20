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

func (s *AuthService) Register(username, password string) error {
	_, err := s.Repo.GetByUsername(username)
	if err == nil {
		return errors.New("username already exists")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("could not hash password")
	}
	user := &model.User{
		Username: username,
		Password: string(hash),
		Type:     "user", // default type
	}
	return s.Repo.Create(user)
}

func (s *AuthService) RegisterWithType(username, password, userType string) error {
	_, err := s.Repo.GetByUsername(username)
	if err == nil {
		return errors.New("username already exists")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("could not hash password")
	}
	user := &model.User{
		Username: username,
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

func (s *AuthService) Login(username, password string) (string, error) {
	user, err := s.Repo.GetByUsername(username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}
	
	return s.GenerateTokenForUser(user)
}

func (s *AuthService) LoginWithType(username, password, userType string) (string, error) {
	user, err := s.Repo.GetByUsername(username)
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
