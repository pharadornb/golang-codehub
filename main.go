package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	host     = "localhost"  // or the Docker service name if running in another container
	port     = 5433         // default PostgreSQL port
	user     = "myuser"     // as defined in docker-compose.yml
	password = "mypassword" // as defined in docker-compose.yml
	dbname   = "mydatabase" // as defined in docker-compose.yml
)

func main() {
	// Configure your PostgreSQL database details here
	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			// IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			// ParameterizedQueries:      true,          // Don't include params in the SQL log
			Colorful: true, // Disable color
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		//panic is kill and show error
		panic("failed to connect to database")
	}

	// db.Migrator().DropColumn(&Book{}, "name")

	// Migrate the schema
	db.AutoMigrate(&Book{})

	fmt.Println("Database migration completed!")

	// newBook := &Book{
	// 	Name:        "The Alchemist",
	// 	Author:      "Paulo Coelho",
	// 	Publisher:   "HarperCollins",
	// 	Description: "A novel by Brazilian author Paulo Coelho that was first published in 1988.",
	// }

	// createBook(db, newBook)

	// currentBook := getBook(db, 1)
	// fmt.Println(currentBook)
}
