package repository

import (
	"github.com/elginbrian/FILKOMPLAIN-BE/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetByUsername(username string) (*model.User, error)
	GetByID(userID uint) (*model.User, error) // Add method to get user by ID
	Update(user *model.User) error
}

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{DB: db}
}

func (r *userRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	if err := r.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByID(userID uint) (*model.User, error) {
	var user model.User
	if err := r.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *model.User) error {
	return r.DB.Save(user).Error
}
