package repository

import (
	"errors"
	"testToDoRestAPI/model"
	"time"

	"gorm.io/gorm"
)

type IUserRepository interface {
	GetUserByEmail(user *model.User, email string) error
	CreateUser(user *model.User) error
	CreatePinRecord(verification *model.VerificationCode) error
	FindVerificationCode(code string) (model.Verificationresponse, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &userRepository{db}
}

func (ur *userRepository) GetUserByEmail(user *model.User, email string) error {
	if err := ur.db.Where("email=?", email).First(user).Error; err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) CreateUser(user *model.User) error {
	if err := ur.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) CreatePinRecord(verification *model.VerificationCode) error {
	if err := ur.db.Create(verification).Error; err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) FindVerificationCode(code string) (model.Verificationresponse, error) {
	var verificationCode model.VerificationCode
	err := ur.db.Where("code = ?", code).First(&verificationCode).Error
	if err != nil {
		return model.Verificationresponse{}, err
	}

	// もし有効期限を確認したいなら、ここで確認し、
	// 期限切れならエラーを返すこともできます。
	if time.Now().After(verificationCode.ExpiresAt) {
		return model.Verificationresponse{}, errors.New("code has expired")
	}

	resUser := model.Verificationresponse{
		UserName: verificationCode.UserName,
		Email:    verificationCode.Email,
		Password: verificationCode.Password,
	}

	return resUser, nil
}
