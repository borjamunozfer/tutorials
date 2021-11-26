package books

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

const bookPath = "/booksfile/"

type BookClient struct {
	book   *os.File
	isOpen bool
}

func (b *BookClient) OpenExistentBook(title string) error {
	//this method is only READ mode
	file, err := b.Open(title)
	if err != nil {
		return err
	}
	b.book = file
	return nil
}

func (b *BookClient) Open(title string) (*os.File, error) {
	file, err := os.Open(bookPath + formatBookFilename(title))

	b.isOpen = err != nil
	return file, err
}

func (b *BookClient) OpenNonExistentBook(title string) error {
	//this method returns an writable and readable file
	file, err := os.OpenFile(bookPath+formatBookFilename(title), os.O_RDWR, 0755)
	if err != nil {
		return err
	}
	b.book = file
	return nil
}

func (b *BookClient) CreateBook(title string) error {

	file, err := os.Create(bookPath + formatBookFilename(title))
	if err != nil {
		return err
	}
	b.book = file
	return nil
}

func (b *BookClient) ReadFullBook(title string) ([]byte, error) {

	fullContent, err := os.ReadFile(bookPath + formatBookFilename(title))
	if err != nil {
		return nil, err
	}
	if len(fullContent) == 0 {
		return nil, fmt.Errorf("Book content is empty")
	}
	return fullContent, nil
}

func (b *BookClient) ReadBookByLine(w io.Writer) error {

	scanner := bufio.NewScanner(b.book)

	for scanner.Scan() {
		fmt.Fprint(w, scanner.Text())
	}

	defer b.book.Close()

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

// Write content to filename specified. If it does not exist, it creates it.
func (b *BookClient) WriteBook(filename string, content string) error {
	err := os.WriteFile("."+bookPath+formatBookFilename(filename), []byte(content), 0755)
	if err != nil {
		return err
	}
	return nil
}

// Write string content
func (b *BookClient) WriteContent(content string) error {

	nbytes, err := b.book.WriteString(content)
	if err != nil {
		return err
	}
	if len(content) != nbytes {
		return fmt.Errorf("Fail to write body %s", string(content))
	}
	return nil
}

func (b *BookClient) DeleteBook(title string) error {

	err := os.Remove(bookPath + formatBookFilename(title))
	if err != nil {
		return err
	}
	return nil
}

func (b *BookClient) Close() error {
	if err := b.book.Close(); err != nil {
		return err
	}
	b.isOpen = false
	return nil
}

func formatBookFilename(title string) string {
	return strings.ToLower(strings.ReplaceAll(title, " ", "") + ".txt")
}
