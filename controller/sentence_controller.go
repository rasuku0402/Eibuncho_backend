package controller

import (
	"net/http"
	"strconv"
	"testToDoRestAPI/model"
	"testToDoRestAPI/usecase"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type ISentenceController interface {
	GetAllSentences(c echo.Context) error
	GetSentenceById(c echo.Context) error
	CreateSentence(c echo.Context) error
	DeleteSentence(c echo.Context) error
}

type sentenceController struct {
	su usecase.ISentenceUsecase
}

func NewSentenceController(su usecase.ISentenceUsecase) ISentenceController {
	return &sentenceController{su}
}

func (sc *sentenceController) GetAllSentences(c echo.Context) error {
	userIdParam := c.Param("userId")
	userId, _ := strconv.Atoi(userIdParam)

	//テスト段階では、以下をコメントアウトし、JWT認証を無効化する
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	jwtuserId := claims["user_id"]

	// // チェック
	if uint(jwtuserId.(float64)) != uint(userId) {
		return c.JSON(http.StatusUnauthorized, "Unauthorized access")
	}

	sentenceRes, err := sc.su.GetAllSentences(uint(jwtuserId.(float64)))
	// ここまで

	// sentenceRes, err := sc.su.GetAllSentences(uint(userId))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, sentenceRes)
}

func (sc *sentenceController) GetSentenceById(c echo.Context) error {
	userIdParam := c.Param("userId")
	userId, _ := strconv.Atoi(userIdParam)

	//テスト段階では、以下をコメントアウトし、JWT認証を無効化する
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	jwtuserId := claims["user_id"]

	// チェック
	if uint(jwtuserId.(float64)) != uint(userId) {
		return c.JSON(http.StatusUnauthorized, "Unauthorized access")
	}
	id := c.Param("sentenceId")
	sentenceId, _ := strconv.Atoi(id)

	// テスト段階では、以下をコメントアウトし、JWT認証を無効化する
	sentenceRes, err := sc.su.GetSentenceById(uint(jwtuserId.(float64)), uint(sentenceId))
	// sentenceRes, err := sc.su.GetSentenceById(uint(userId), uint(sentenceId))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, sentenceRes)
}

func (sc *sentenceController) CreateSentence(c echo.Context) error {
	userIdParam := c.Param("userId")
	userId, _ := strconv.Atoi(userIdParam)

	//*テスト段階では、以下をコメントアウトし、JWT認証を無効化する
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	jwtuserId := claims["user_id"]
	sentence := model.Sentences{}
	sentence.UserId = uint(jwtuserId.(float64))
	//*ここまでコメントアウト
	if uint(jwtuserId.(float64)) != uint(userId) {
		return c.JSON(http.StatusUnauthorized, "Unauthorized access")
	}

	// sentence := model.Sentences{}
	if err := c.Bind(&sentence); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	// sentence.UserId = uint(userId)
	sentenceRes, err := sc.su.CreateSentence(sentence)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, sentenceRes)
}

func (sc *sentenceController) DeleteSentence(c echo.Context) error {
	userIdParam := c.Param("userId")
	userId, _ := strconv.Atoi(userIdParam)

	//*テスト段階では、以下をコメントアウトし、JWT認証を無効化する
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	jwtuserId := claims["user_id"]
	sentence := model.Sentences{}
	sentence.UserId = uint(jwtuserId.(float64))
	//*ここまでコメントアウト
	if uint(jwtuserId.(float64)) != uint(userId) {
		return c.JSON(http.StatusUnauthorized, "Unauthorized access")
	}

	id := c.Param("sentenceId")
	sentenceId, _ := strconv.Atoi(id)
	err := sc.su.DeleteSentence(uint(userId), uint(sentenceId))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, "Sentence is deleted.")

}
