package repository

import (
	"github.com/elginbrian/FILKOMPLAIN-BE/internal/model"
	"gorm.io/gorm"
)

type AuthRepository interface {
	GetByUsername(username string) (*model.User, error)
	Create(user *model.User) error
	Update(user *model.User) error
}

type authRepository struct {
	DB *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{DB: db}
}

func (r *authRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	if err := r.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) Create(user *model.User) error {
	return r.DB.Create(user).Error
}

func (r *authRepository) Update(user *model.User) error {
	return r.DB.Save(user).Error
}
