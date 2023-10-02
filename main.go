package main

import (
	"fmt"
	"testToDoRestAPI/controller"
	"testToDoRestAPI/db"
	"testToDoRestAPI/model"
	"testToDoRestAPI/repository"
	"testToDoRestAPI/router"
	"testToDoRestAPI/usecase"
	"testToDoRestAPI/validator"
	"time"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

func main() {
	db := db.NewDB()
	userValidator := validator.NewUserValidator()
	sentenceValidator := validator.NewSentenceValidator()
	userRepository := repository.NewUserRepository(db)
	sentenceRepository := repository.NewSentenceRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepository, userValidator)
	sentenceUsecase := usecase.NewSentenceUsecase(sentenceRepository, sentenceValidator)
	userController := controller.NewUserController(userUsecase)
	sentenceController := controller.NewSentenceController(sentenceUsecase)
	e := router.NewRouter(userController, sentenceController)

	// Initialize and start the cron scheduler
	c := cron.New(cron.WithSeconds())
	_, err := c.AddFunc("0 */3 * * * *", func() {
		fmt.Println("Running job every 3 minutes")
		DeleteOldRecords(db) // dbはあなたのデータベースのインスタンス
	})
	if err != nil {
		fmt.Printf("Error scheduling job: %s\n", err)
		return
	}
	c.Start()

	e.Logger.Fatal(e.Start(":8080"))
}

func DeleteOldRecords(db *gorm.DB) {
	// gormを使ってクエリを実行する
	err := db.Where("expires_at < ?", time.Now()).Delete(&model.VerificationCode{}).Error
	if err != nil {
		fmt.Printf("Error deleting old records: %s\n", err)
	} else {
		fmt.Println("Successfully deleted old records")
	}
}
