package utils

import "github.com/mishozz/Library/entities"

func Contains(slice []entities.Book, book entities.Book) bool {
	for _, x := range slice {
		if x.Isbn == book.Isbn {
			return true
		}
	}
	return false
}

func Remove(slice []entities.Book, book entities.Book) []entities.Book {
	var s int
	for index, x := range slice {
		if x.Isbn == book.Isbn {
			s = index
			break
		}
	}
	return append(slice[:s], slice[s+1:]...)
}
