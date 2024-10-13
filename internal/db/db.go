package db

import (
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	database, err := gorm.Open(sqlite.Open("myapp.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
		return nil, err
	}

	tableExists, err := checkIfTableExists(database, "messages")
	if err != nil {
		return nil, err
	}

	if !tableExists {
		fmt.Println("Messages table does not exist, creating...")

		err := database.AutoMigrate(&Message{})
		if err != nil {
			return nil, err
		}

		err = insertInitialMessages(database)
		if err != nil {
			return nil, err
		}

		fmt.Println("Table created and 100 rows inserted.")
	}

	return database, nil
}

func checkIfTableExists(db *gorm.DB, tableName string) (bool, error) {
	exists := db.Migrator().HasTable(tableName)
	return exists, nil
}

func insertInitialMessages(db *gorm.DB) error {
	for i := 1; i <= 100; i++ {
		content := fmt.Sprintf("Message #%d", i)
		receiver := fmt.Sprintf("+90555111%04d", i)
		status := "unsent"

		message := Message{
			Content:  content,
			Receiver: receiver,
			Status:   status,
		}

		if err := db.Create(&message).Error; err != nil {
			return err
		}
	}

	return nil
}
