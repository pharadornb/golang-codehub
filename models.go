package main

import (
	"fmt"

	"gorm.io/gorm"
)

type Book struct {
	//camel case to snell case pass gorm
	gorm.Model
	// ID          int64
	Name        string
	Author      string
	Publisher   string
	Description string
	// price       uint
}

func createBook(db *gorm.DB, book *Book) {
	result := db.Create(book)

	if result.Error != nil {
		panic(fmt.Sprintf("Error creating book: %v", result.Error))
	}

	fmt.Println("Book created successfully!")
}

func getBook(db *gorm.DB, id uint) *Book {
	var book Book
	result := db.First(&book, id)

	if result.Error != nil {
		panic(fmt.Sprintf("Error getting book: %v", result.Error))
	}

	return &book
}


