package db

import (
	"fmt"
	"log"
	"sync"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	dbInstance *gorm.DB
	once       sync.Once
)

func ConnectDB(path string) (*gorm.DB, error) {
	once.Do(func() {
		var err error
		dbInstance, err = gorm.Open(sqlite.Open(path), &gorm.Config{})
		if err != nil {
			log.Fatalf("Failed to connect to the dbInstance: %v", err)
		}

		tableExists, err := checkIfTableExists(dbInstance, "messages")
		if err != nil {
			log.Fatalf("Failed to check if table exists: %v", err)
		}

		if !tableExists {
			fmt.Println("Messages table does not exist, creating...")

			err := dbInstance.AutoMigrate(&Message{})
			if err != nil {
				log.Fatalf("Failed to create table: %v", err)
			}

			err = insertInitialMessages(dbInstance)
			if err != nil {
				log.Fatalf("Failed to insert initial messages: %v", err)
			}

			fmt.Println("Table created and 100 rows inserted.")
		}
	})

	return dbInstance, nil
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
