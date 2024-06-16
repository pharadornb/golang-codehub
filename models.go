package main

import (
	"fmt"

	"gorm.io/gorm"
)

type Book struct {
	//camel case to snell case pass gorm
	gorm.Model
	// ID          int64
	Name        string `json:"name"`
	Author      string `json:"author"`
	Publisher   string `json:"publisher"`
	Description string `json:"description"`
	Price       uint   `json:"price"`
}

func createBook(db *gorm.DB, book *Book) error {
	result := db.Create(book)

	if result.Error != nil {
		panic(fmt.Sprintf("Error creating book: %v", result.Error))
	}

	// fmt.Println("Book created successfully!")
	return nil
}

func getBook(db *gorm.DB, id uint) *Book {
	var book Book
	result := db.First(&book, id)

	if result.Error != nil {
		panic(fmt.Sprintf("Error getting book: %v", result.Error))
	}

	return &book
}

func updateBook(db *gorm.DB, book *Book) error {
	// result := db.Model(&book).Where("id = ?", id).Updates(book)
	result := db.Model(&book).Updates(book)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func deleteBook(db *gorm.DB, id int) error {
	result := db.Delete(&Book{}, id)

	if result.Error != nil {
		// panic(fmt.Sprintf("Error deleting book: %v", result.Error))
		return result.Error
	}

	return nil

	// fmt.Println("Book deleted successfully!")
}

// func seachBook(db *gorm.DB, bookName string) {
// 	var book Book
// 	result := db.Where("name = ?", bookName).First(&book)

// 	if result.Error != nil {
// 		panic(fmt.Sprintf("Error searching book: %v", result.Error))
// 	}

// 	fmt.Println("Book found successfully!")
// }

func seachBookAll(db *gorm.DB, bookName string) []Book {
	var books []Book
	result := db.Where("name = ?", bookName).Order("price desc").Find(&books)

	if result.Error != nil {
		panic(fmt.Sprintf("Error searching book: %v", result.Error))
	}

	fmt.Println("Book found successfully!")
	return books
}

func getBooks(db *gorm.DB) []Book {
	var books []Book
	result := db.Find(&books)

	if result.Error != nil {
		panic(fmt.Sprintf("Error searching book: %v", result.Error))
	}

	fmt.Println("Book found successfully!")
	return books
}
