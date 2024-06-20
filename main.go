package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

//Book 1-1 Publisher
//Books M-N Authors

type Book struct {
	gorm.Model
	Name        string `json:"name"`
	Author      string `json:"author"`
	Description string `json:"description"`
	PublisherID uint
	Publisher   Publisher
	Authors     []Author `gorm:"many2many:author_books;"`
}

type Publisher struct {
	gorm.Model
	Details string
	Name    string
}

type Author struct {
	gorm.Model
	Name  string
	Books []Book `gorm:"many2many:author_books;"`
}

type AuthorBook struct {
	AuthorID uint
	Author   Author
	BookID   uint
	Book     Book
}

// func createPublisher(db *gorm.DB, publisher *Publisher) error {
// 	result := db.Create(publisher)
// 	if result.Error != nil {
// 		return result.Error
// 	}
// 	return nil
// }

// func createAuthor(db *gorm.DB, author *Author) error {
// 	result := db.Create(author)
// 	if result.Error != nil {
// 		return result.Error
// 	}
// 	return nil
// }

// // , _ []uint
// func createBookWithAuthor(db *gorm.DB, book *Book) error {
// 	// First, create the book
// 	if err := db.Create(book).Error; err != nil {
// 		return err
// 	}

// 	return nil
// }

func getBookWithPublisher(db *gorm.DB, bookID uint) (*Book, error) {
	var book Book
	result := db.Preload("Publisher").First(&book, bookID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &book, nil
}

func getBookWithAuthors(db *gorm.DB, bookID uint) (*Book, error) {
	var book Book
	result := db.Preload("Authors").First(&book, bookID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &book, nil
}

func listBooksOfAuthor(db *gorm.DB, authorID uint) ([]Book, error) {
	var books []Book
	result := db.Joins("JOIN author_books on author_books.book_id = books.id").
		Where("author_books.author_id = ?", authorID).
		Find(&books)
	if result.Error != nil {
		return nil, result.Error
	}
	return books, nil
}

const (
	host     = "localhost"  // or the Docker service name if running in another container
	port     = 5433         // default PostgreSQL port
	user     = "myuser"     // as defined in docker-compose.yml
	password = "mypassword" // as defined in docker-compose.yml
	dbname   = "mydatabase" // as defined in docker-compose.yml
)

// func authRequired(c *fiber.Ctx) error {
// 	cookie := c.Cookies("jwt")
// 	jwtSecretKey := "TestSecret"

// 	// StandardClaims{} is a struct that implements the Claims interface
// 	//MapClaims{} is a struct that implements the Claims interface
// 	//RegisteredClaims{} is a struct that implements the Claims interface
// 	token, err := jwt.ParseWithClaims(cookie, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
// 		return []byte(jwtSecretKey), nil
// 	})

// 	if err != nil || !token.Valid {
// 		return c.SendStatus(fiber.StatusUnauthorized)
// 	}

// 	claim := token.Claims.(*jwt.MapClaims)
// 	fmt.Println(claim["name"].(string))

//		return c.Next()
//	}
func authRequired(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	jwtSecretKey := "TestSecret"

	// StandardClaims{} is a struct that implements the Claims interface
	//MapClaims{} is a struct that implements the Claims interface
	//RegisteredClaims{} is a struct that implements the Claims interface
	token, err := jwt.ParseWithClaims(cookie, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecretKey), nil
	})

	if err != nil || !token.Valid {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// claim := token.Claims.(*jwt.MapClaims)
	// fmt.Println((*claim)["name"].(string))

	return c.Next()
}

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

	// db.Migrator().DropColumn(Booker[])
	// db.Migrator().DropTable(&Booker{})
	// db.Migrator().DropTable("users")

	// Migrate the schema
	db.AutoMigrate(&Booker{}, &User{})
	db.AutoMigrate(&Book{}, &Publisher{}, &Author{}, &AuthorBook{})

	// ขาสร้าง
	// publisher := Publisher{
	// 	Details: "Bomb",
	// 	Name:    "B",
	// }
	// _ = createPublisher(db, &publisher)

	// // Example data for a new author
	// author1 := Author{
	// 	Name: "B",
	// }
	// _ = createAuthor(db, &author1)

	// author2 := Author{
	// 	Name: "B",
	// }
	// _ = createAuthor(db, &author2)

	// // // Example data for a new book with an author
	// book := Book{
	// 	Name:        "GG",
	// 	Author:      "FFFFEDXDD",
	// 	Description: "Book Description",
	// 	PublisherID: publisher.ID,               // Use the ID of the publisher created above
	// 	Authors:     []Author{author1, author2}, // Add the created author
	// }
	// _ = createBookWithAuthor(db, &book)

	// ขาเรียก

	// Example: Get a book with its publisher
	// bookWithPublisher, err := getBookWithPublisher(db, 1) // assuming a book with ID 1
	// // if err != nil {
	// // 	// Handle error
	// // }

	// // // Example: Get a book with its authors
	// bookWithAuthors, err := getBookWithAuthors(db, 1) // assuming a book with ID 1
	// // if err != nil {
	// // 	// Handle error
	// // }

	// // // Example: List books of a specific author
	// authorBooks, err := listBooksOfAuthor(db, 1) // assuming an author with ID 1
	// // if err != nil {
	// // 	// Handle error
	// // }

	// fmt.Println("----------------------------")
	// fmt.Println(bookWithPublisher)
	// fmt.Println("----------------------------")
	// fmt.Println(bookWithAuthors)
	// fmt.Println("----------------------------")
	// fmt.Println(authorBooks)

	// Example: Get a book with its publisher
	bookWithPublisher, err := getBookWithPublisher(db, 1) // assuming a book with ID 1
	if err != nil {
		log.Printf("Error getting book with publisher: %v", err)
	} else {
		fmt.Println("----------------------------")
		fmt.Println(bookWithPublisher)
	}

	// Example: Get a book with its authors
	bookWithAuthors, err := getBookWithAuthors(db, 1) // assuming a book with ID 1
	if err != nil {
		log.Printf("Error getting book with authors: %v", err)
	} else {
		fmt.Println("----------------------------")
		fmt.Println(bookWithAuthors)
	}

	// Example: List books of a specific author
	authorBooks, err := listBooksOfAuthor(db, 1) // assuming an author with ID 1
	if err != nil {
		log.Printf("Error listing books of author: %v", err)
	} else {
		fmt.Println("----------------------------")
		fmt.Println(authorBooks)
	}

	//set up fiber app
	app := fiber.New()
	app.Use("/books", authRequired)

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
		newBook := new(Booker)

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

		book := new(Booker)

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

		c.Cookie(&fiber.Cookie{
			Name:     "jwt",
			Value:    token,
			Expires:  time.Now().Add(time.Hour * 72),
			HTTPOnly: true,
		})

		return c.JSON(fiber.Map{
			"status": "Login Success",
			// "token":  token,
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
