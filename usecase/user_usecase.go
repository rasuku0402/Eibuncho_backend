package usecase

import (
	"os"
	"testToDoRestAPI/model"
	"testToDoRestAPI/repository"
	"testToDoRestAPI/validator"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type IUserUsecase interface {
	// SignUp(user model.User) (model.UserResponse, error)
	SignUp(verification model.VerificationCode) error
	Login(user model.User) (model.UserResponse, string, error)
	FindVerificationCode(req model.VerificationRequest) error
}

type userUsecase struct {
	ur repository.IUserRepository
	uv validator.IUserValidator
}

func NewUserUsecase(ur repository.IUserRepository, uv validator.IUserValidator) IUserUsecase {
	return &userUsecase{ur, uv}
}

// func (uu *userUsecase) SignUp(user model.User) (model.UserResponse, error) {
// 	if err := uu.uv.UserValidate(user); err != nil {
// 		return model.UserResponse{}, err
// 	}
// 	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
// 	if err != nil {
// 		return model.UserResponse{}, err
// 	}
// 	newUser := model.User{UserName: user.UserName, Email: user.Email, Password: string(hash)}

// 	if err := uu.ur.CreateUser(&newUser); err != nil {
// 		return model.UserResponse{}, err
// 	}

// 	resUser := model.UserResponse{
// 		UserID:   newUser.UserID,
// 		UserName: newUser.UserName,
// 		Email:    newUser.Email,
// 	}
// 	return resUser, nil
// }

func (uu *userUsecase) SignUp(verification model.VerificationCode) error {
	if err := uu.uv.RecordValidate(verification); err != nil {
		return err
	}
	// passwordをハッシュ化
	pshash, err := bcrypt.GenerateFromPassword([]byte(verification.Password), 10)
	if err != nil {
		return err
	}

	// Codeをハッシュ化
	// codehash, err := bcrypt.GenerateFromPassword([]byte(verification.Code), 10)
	// if err != nil {
	// 	return err
	// }
	newVerifyrecord := model.VerificationCode{UserName: verification.UserName, Email: verification.Email, Code: verification.Code, Password: string(pshash), ExpiresAt: verification.ExpiresAt}

	if err := uu.ur.CreatePinRecord(&newVerifyrecord); err != nil {
		return err
	}
	return nil
}

func (uu *userUsecase) Login(user model.User) (model.UserResponse, string, error) {
	if err := uu.uv.UserValidate(user); err != nil {
		return model.UserResponse{}, "", err
	}
	storedUser := model.User{}
	if err := uu.ur.GetUserByEmail(&storedUser, user.Email); err != nil {
		return model.UserResponse{}, "", err
	}
	err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password))
	if err != nil {
		return model.UserResponse{}, "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": storedUser.UserID,
		"exp":     time.Now().Add(time.Hour * 12).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return model.UserResponse{}, "", err
	}

	//フロントの処理のために、user_idとかをresponseさせる

	resUser := model.UserResponse{
		UserID:   storedUser.UserID,
		UserName: storedUser.UserName,
		Email:    storedUser.Email,
	}

	return resUser, tokenString, nil
}

func (uu *userUsecase) FindVerificationCode(req model.VerificationRequest) error {
	// var verificationCode model.VerificationCode
	// req.codeをハッシュ化
	// reqcodehash, err := bcrypt.GenerateFromPassword([]byte(req.Code), 10)
	// if err != nil {
	// 	return err
	// }
	code := string(req.Code)

	resUser, err := uu.ur.FindVerificationCode(code)
	if err != nil {
		return err // "Invalid or expired code"
	}

	// ユーザーレコードを作成
	user := model.User{
		UserName: resUser.UserName,
		Email:    resUser.Email,
		Password: resUser.Password,
	}

	// ユーザーをデータベースに保存
	if err := uu.ur.CreateUser(&user); err != nil {
		return err
	}

	return nil

}
