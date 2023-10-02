package main

import (
	"fmt"
	"testToDoRestAPI/db"
	"testToDoRestAPI/model"
)

func main() {
	dbConn := db.NewDB()
	defer fmt.Println("Successfully Migrated")
	defer db.CloseDB(dbConn)
	dbConn.AutoMigrate(&model.User{}, &model.Sentences{}, &model.VerificationCode{})
}
