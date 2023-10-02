package controller

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"testToDoRestAPI/model"
	"testToDoRestAPI/usecase"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/labstack/echo/v4"
)

var sess = session.Must(session.NewSession(&aws.Config{
	Region: aws.String("ap-northeast-1"), // 適切なリージョンに変更してください
}))
var svc = ses.New(sess)

// メール送信関数を実装
func sendEmail(recipient string, subject string, bodyText string) error {
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{aws.String(recipient)},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Data: aws.String(bodyText),
				},
			},
			Subject: &ses.Content{
				Data: aws.String(subject),
			},
		},
		Source: aws.String("taitoo0402@gmail.com"), // ensure this is a verified email in SES
	}
	_, err := svc.SendEmail(input)
	return err
}

// Email validation function
func isValidEmail(email string) bool {
	// Simple regex check for valid email format
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

type IUserController interface {
	SignUp(c echo.Context) error
	LogIn(c echo.Context) error
	LogOut(c echo.Context) error
	Validate(c echo.Context) error
}

type userController struct {
	uu usecase.IUserUsecase
}

func NewUserController(uu usecase.IUserUsecase) IUserController {
	return &userController{uu}
}

// func (uc *userController) SignUp(c echo.Context) error {
// 	user := model.User{}
// 	if err := c.Bind(&user); err != nil {
// 		return c.JSON(http.StatusBadRequest, err.Error())
// 	}

// 	// ユーザーのemailアドレスを取得する
// 	email := user.Email
// 	fmt.Println("Received email:", email) // Log output for debugging
// 	// Email validation check here
// 	if !isValidEmail(email) {
// 		return c.JSON(http.StatusBadRequest, "Invalid email address provided")
// 	}

// 	// ランダムなPINコードを生成する (例として6桁の数字)
// 	pinCode := fmt.Sprintf("%06d", rand.Intn(1000000))

// 	// ... データベースなどにPINコードを保存 ...(まだ未実装)

// 	// メールを送信する
// 	err := sendEmail(email, "Your verification code", "Your verification code is: "+pinCode)
// 	if err != nil {
// 		fmt.Printf(email)
// 		return c.JSON(http.StatusInternalServerError, err.Error())
// 	}

// 	userRes, err := uc.uu.SignUp(user)
// 	if err != nil {

// 		return c.JSON(http.StatusInternalServerError, err.Error())
// 	}

// 	return c.JSON(http.StatusOK, userRes)
// }

func (uc *userController) SignUp(c echo.Context) error {
	verification := model.VerificationCode{}
	if err := c.Bind(&verification); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// ユーザのemailを取得する
	email := verification.Email
	fmt.Println("Received email:", email) // デバック用にemailを標準出力

	// Email validation check here
	if !isValidEmail(email) {
		return c.JSON(http.StatusBadRequest, "Invalid email address provided")
	}

	// ランダムなPINコードを生成する (例として6桁の数字)
	Code := fmt.Sprintf("%06d", rand.Intn(1000000))
	verification.Code = Code

	// 有効期限を1分に設定
	verification.ExpiresAt = time.Now().Add(1 * time.Minute)

	// メールを送信する
	err := sendEmail(email, "Your verification code", "Your verification code is: "+Code)
	if err != nil {
		fmt.Printf(email)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	//データベースなどにPINコードを保存
	if err := uc.uu.SignUp(verification); err != nil {
		return err
	}
	return nil
}

func (uc *userController) LogIn(c echo.Context) error {
	user := model.User{}
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	userRes, tokenString, err := uc.uu.Login(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = tokenString
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.Path = "/"
	cookie.Domain = os.Getenv("API_DOMAIN")
	// cookie.Secure = true
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteNoneMode
	c.SetCookie(cookie)
	return c.JSON(http.StatusOK, userRes)

}

func (uc *userController) LogOut(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = ""
	cookie.Expires = time.Now()
	cookie.Path = "/"
	cookie.Domain = os.Getenv("API_DOMAIN")
	// cookie.Secure = true
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteNoneMode
	c.SetCookie(cookie)
	return c.NoContent(http.StatusOK)
}

// フロントから、PINを受け取って、認証OKならusersテーブルにレコードを追加する
func (uc *userController) Validate(c echo.Context) error {
	// リクエストボディをバインド
	req := model.VerificationRequest{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "Who?")
	}
	// FindVerificationCodeは適切に実装する必要があります
	err := uc.uu.FindVerificationCode(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid or expired code")
	}

	return c.JSON(http.StatusOK, "User successfully created")
}
