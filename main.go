package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
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
	db.AutoMigrate(&Book{}, &User{})

	//set up fiber app
	app := fiber.New()

	app.Get("/books", func(c *fiber.Ctx) error {
		// return c.SendString("Hello, World!")
		return c.JSON(getBooks(db))
	})

	app.Get("/book/:id", func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		book := getBook(db, uint(id))
		return c.JSON(book)
	})

	app.Post("/book", func(c *fiber.Ctx) error {
		newBook := new(Book)

		if err := c.BodyParser(newBook); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		createBook(db, newBook)

		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		return c.JSON(fiber.Map{
			"status": "success",
			"data":   newBook,
		})
	})

	app.Put("/book/:id", func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))

		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		book := new(Book)

		if err := c.BodyParser(book); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		book.ID = uint(id)

		err = updateBook(db, book)

		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		return c.JSON(fiber.Map{
			"status": "success",
			"data":   book,
		})
	})

	app.Delete("/book/:id", func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))

		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		err = deleteBook(db, id)

		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		return c.JSON(fiber.Map{
			"status": "success",
		})
	})

	currentBook := seachBookAll(db, "The Alchemist")
	fmt.Println(currentBook)

	app.Post("/register", func(c *fiber.Ctx) error {
		user := new(User)

		if err := c.BodyParser(user); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		err = createUser(db, user)

		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		return c.JSON(fiber.Map{
			"status": "Register Success",
		})
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		user := new(User)

		if err := c.BodyParser(user); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		token, err := loginUser(db, user)

		if err != nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		return c.JSON(fiber.Map{
			"status": "Login Success",
			"token":  token,
		})

		// err = createUser(db, user)

		// if err != nil {
		// 	return c.SendStatus(fiber.StatusBadRequest)
		// }

		// return c.JSON(fiber.Map{
		// 	"status": "Register Success",
		// })
	})

	// for _, book := range currentBook {
	// 	fmt.Println(book)
	// }

	app.Listen(":8083")

	// fmt.Println("Database migration completed!")

	// newBook := &Book{
	// 	Name:        "The Alchemist",
	// 	Author:      "Paulo Coelho",
	// 	Publisher:   "HarperCollins",
	// 	Description: "A novel by Brazilian author Paulo Coelho that was first published in 1988.",
	// }
	// createBook(db, newBook)

	// currentBook := getBook(db, 1)
	// fmt.Println(currentBook)

	// currentBook.Name = "The Alchemist (Updated)"
	// currentBook.Author = "Paulo Coelho (Updated)"
	// updateBook(db, 1, currentBook)

	// deleteBook(db, 1)

	//not use gorm.model can't delete from field - soft delete

	// currentBook := seachBook(db, "The Alchemist")
	// fmt.Println(currentBook)
}
