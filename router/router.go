package router

import (
	"os"
	"testToDoRestAPI/controller"

	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
)

func NewRouter(uc controller.IUserController, sc controller.ISentenceController) *echo.Echo {
	e := echo.New()
	e.POST("/signup", uc.SignUp)
	e.POST("/validate", uc.Validate)
	e.POST("/login", uc.LogIn)
	e.POST("/logout", uc.LogOut)

	s := e.Group("/:userId/sentences")
	s.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(os.Getenv("SECRET")),
		TokenLookup: "cookie:token",
	}))
	s.GET("", sc.GetAllSentences)
	s.GET("/:sentenceId", sc.GetSentenceById)
	s.POST("", sc.CreateSentence)
	s.DELETE("/:sentenceId", sc.DeleteSentence)

	return e
}
