package books

import (
	"bufio"
	"fmt"
	"os"
)

type BookClient struct {
	book *os.File
}

func (b *BookClient) OpenBook(title string) error {
	file, err := os.OpenFile(title, os.O_RDWR, 0755)
	if err != nil {
		return err
	}

	b.book = file
	return nil
}

func (b *BookClient) ReadFullBook(title string) ([]byte, error) {
	fullContent, err := os.ReadFile(title)
	if err != nil {
		return nil, err
	}

	return fullContent, nil
}

func (b *BookClient) ReadBookByLine() {

	scanner := bufio.NewScanner(b.book)
	for scanner.Scan() {
		fmt.Printf(scanner.Text())
	}
}

// Write content to filename specified. If it does not exist, it creates it.
func (b *BookClient) WriteBook(filename string, content string) error {
	err := os.WriteFile(filename, []byte(content), 0)
	if err != nil {
		return err
	}
	return nil
}
