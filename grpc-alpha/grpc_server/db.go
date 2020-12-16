package main

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	dbDSN = "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable"
)


// Task to be assigned
type Task struct {
	ClientID  string `gorm:"primary_key" json:"client_id"`
	TaskID    int32  `gorm:"primary_key;auto_increment:false" json:"task_id"`
	CommandID int32  `json:"command_id"`
	Args      string `json:"command_args"`
	Completed bool   `gorm:"default:false"`
	Status    string `gorm:"default:Not Run"`  
}

// TempURL to be assigned
type TempURL struct {
	URLPost  string `gorm:"primary_key"`
	FilePath string `gorm:"not null" json:"file_path"`
	Valid    bool   `gorm:"default:true"`
}

// OpenDB function used to connect to dataase and returns db object
func OpenDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	// db, err := gorm.Open(postgres.Open(dbDSN), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	CreateDB(db)
	return db
}

// CreateDB function that auto migrates and creates database if it does not exist
func CreateDB(db *gorm.DB) {
	db.AutoMigrate(&Task{}, &TempURL{})
}

func test() {

	db := OpenDB()
	db.Create(&Task{ClientID: "Test", TaskID: 2, CommandID: 8, Args: "Bob", Completed: true})
	tasks := make([]Task, 0)
	db.Find(&tasks)
	fmt.Println(tasks)
}
