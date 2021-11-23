package books

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

const bookPath = "/booksfile/"

type Author struct {
	Name  string
	Books map[string]*os.File
}

func NewAuthor(name string) Author {
	books := make(map[string]*os.File)
	return Author{Name: name, Books: books}
}

func (a *Author) CreateBook(title string) error {
	//Check book for this author does not exist
	if _, ok := a.Books[title]; ok {
		return fmt.Errorf("Book %s for author %s already exists", title, a.Name)
	}

	log.Printf("Book %s does not exist, proceeding to create it...", title)

	//Need to create empty file
	createdFile, err := newFile(title)
	createdFile.Close()
	if err != nil {
		return err
	}
	a.Books[title] = createdFile
	return nil
}

func formatBookFilename(title string) string {
	return strings.ToLower(strings.ReplaceAll(title, " ", "") + ".txt")
}

func newFile(title string) (*os.File, error) {
	toValidFilename := formatBookFilename(title)

	//create dir if not exist
	//if any err received creating directory, then exit program
	err := bookDirExists()
	if err != nil {
		log.Fatal(err)
	}
	//create file
	bookFile, err := os.Create(bookPath + toValidFilename)
	if err != nil {
		return nil, err
	}

	return bookFile, nil
}

func bookDirExists() error {
	//check if path exist; create if not
	if _, err := os.Stat(bookPath); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(bookPath, os.ModePerm)
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

func (a *Author) RemoveLine() {

}

func (a *Author) RemoveBook() {

}

func (a *Author) ReadBook() {

}

func (a *Author) WriteBook(book *os.File, body string) error {

	n, err := book.WriteString(body)
	defer book.Close()
	if err != nil {
		return err
	}
	if len(body) != n {
		return fmt.Errorf("Fail to write body")
	}

	return nil
}
