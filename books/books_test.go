package books

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/cucumber/godog"
)

const (
	bookKey              = "filename"
	contentKey           = "content"
	bookTitleNonExistent = "Non Existent"
	content              = `Háblame, Musa, del varón de gran ingenio, que anduvo errante muchísimo tiempo.
y que vio las ciudades de muchos hombres y conoció su manera de pensar,
pero padeció aún en el mar muchos dolores en su ánimo,
procurando conservar su vida y el regreso de sus compañeros.
Mas ni siquiera así terminó de salvar a sus compañeros, aunque lo deseaba vivamente, pues perecerieron por sus propios actos temerarios.`
)

type BookFeature struct {
	client *BookClient
}

//dummy step, just ignore it
func (bf *BookFeature) iGotOneLineAtTime() error {
	return nil
}

func (bf *BookFeature) iOpenTheBook(ctx context.Context) error {
	//get book title from filename. It contains the relative path (booksfile/laodisea.txt) so we have to clean up and extract only laodisea
	//split return in a string slice [0] -> booksfile and [1] -> laodisea.txt
	filenameWithExt := strings.Split(bf.client.book.Name(), "/")
	bookTitleformat := strings.Split(filenameWithExt[2], ".txt")[0]
	err := bf.client.OpenExistentBook(bookTitleformat)
	if err != nil {
		return err
	}
	return nil
}

func (bf *BookFeature) iOpenTheBookNonExistent(ctx context.Context) error {
	//get book title from context
	bookTitle, _ := ctx.Value(bookTitleNonExistent).(string)
	if len(bookTitle) == 0 {
		return fmt.Errorf("Could not recover booktile from Context")
	}

	err := bf.client.OpenNonExistentBook(bookTitle)
	if err != nil {
		return err
	}
	//check is created
	if _, err := bf.client.book.Stat(); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("File %s is not created", bookTitle)
		}
	}

	return nil
}

func (bf *BookFeature) iReadTheBookByLine() error {
	defer bf.client.book.Close()

	mockOutputBuffer := new(bytes.Buffer)
	err := bf.client.ReadBookByLine(mockOutputBuffer)
	if err != nil {
		return err
	}
	contentBuf := make([]byte, len(content))
	nbytes, err := mockOutputBuffer.Read(contentBuf)
	if nbytes != len(content) {
		return err
	}
	return nil
}

func (bf *BookFeature) iReadTheBookFully(ctx context.Context, bookTitle string) (context.Context, error) {
	content, err := bf.client.ReadFullBook(bookTitle)
	if err != nil {
		return ctx, err
	}
	return context.WithValue(ctx, contentKey, content), nil
}

func (bf *BookFeature) iWriteContent() error {
	err := bf.client.WriteContent(content)
	if err != nil {
		return err
	}
	return nil
}

func (bf *BookFeature) iWriteContentToUnexistentBook(ctx context.Context) (context.Context, error) {
	err := bf.client.WriteBook(bookTitleNonExistent, content)
	if err != nil {
		return ctx, err
	}
	return context.WithValue(ctx, bookKey, formatBookFilename(bookTitleNonExistent)), nil
}

func (bf *BookFeature) theBookContentIsReturned(ctx context.Context) error {
	contentBytes, _ := ctx.Value(contentKey).([]byte)
	if len(content) == 0 {
		return fmt.Errorf("Content from context is empty")
	}

	if string(contentBytes) != content {
		return fmt.Errorf("Content from context is %v \n but want %s \n", string(contentBytes), content)
	}

	return nil
}

func (bf *BookFeature) theBookContentIsUpdated() error {

	//read by bytes
	defer bf.client.book.Close()
	buf := make([]byte, len(content))
	//move pointer to initial file and read from that
	_, err := bf.client.book.Seek(0, 0)
	if err != nil {
		return err
	}
	nbytes, err := bf.client.book.Read(buf)
	if nbytes != len(content) {
		return fmt.Errorf("Bytes content readed from %s not equals", bf.client.book.Name())
	}
	return nil
}

func (bf *BookFeature) theBookDoesNotExists(ctx context.Context, bookTitle string) (context.Context, error) {

	if _, err := os.Stat(bookPath + bookTitle + ".txt"); os.IsExist(err) {
		return ctx, err
	}

	return context.WithValue(ctx, bookTitleNonExistent, bookTitle), nil
}

func (bf *BookFeature) theBookIsCreatedAndUpdated(ctx context.Context) error {

	//recover from context
	bookTitle, ok := ctx.Value("Non Existent").(string)
	if !ok {
		return fmt.Errorf("Could not get Book title NonExistent from context")
	}

	//check is created and compare content
	//we can use stat to check is created but we're going to open fully
	content, err := bf.client.ReadFullBook(bookTitle)
	if err != nil {
		if os.IsNotExist(err) {
			return err
		}
	}

	//compare content
	if len(content) == 0 {
		return fmt.Errorf("Content is empty %v", content)
	}

	return nil

}

func (bf *BookFeature) theFileIsCreatedEmpty() error {
	//check book with that filename exists.
	defer bf.client.book.Close()
	fileInfo, err := bf.client.book.Stat()
	if err != nil {
		return err
	}
	//check bookfile size is 0 (empty)
	if fileInfo.Size() != 0 {
		return fmt.Errorf("Book file exists but its not empty")
	}
	return nil
}

func (bf *BookFeature) theFileIsOpenedCorrectly() error {

	defer bf.client.book.Close()
	if _, err := bf.client.book.Stat(); err != nil {
		return err
	}

	return nil
}

func (bf *BookFeature) iCreateTheBook(bookTitle string) error {

	err := bf.client.CreateBook(bookTitle)
	if err != nil {
		return err
	}

	return nil
}

func (bf *BookFeature) theBookDoesNotExist(ctx context.Context, bookTitle string) (context.Context, error) {

	if _, err := os.Stat(formatBookFilename(bookTitle)); err != nil {
		if os.IsExist(err) {
			return ctx, fmt.Errorf("File %s exists", bookTitle)
		}
	}
	return context.WithValue(ctx, bookTitleNonExistent, bookTitle), nil
}

func (bf *BookFeature) theBookExists(bookTitle string) error {
	if _, err := os.Stat(formatBookFilename(bookTitle)); err != nil {
		if os.IsNotExist(err) {
			// we create it
			err = bf.client.CreateBook(bookTitle)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (bf *BookFeature) theFileIsCreatedAndOpenedCorrectly() error {
	if _, err := bf.client.book.Stat(); err != nil {
		return err
	}
	bf.client.book.Close()
	return nil
}

func (bf *BookFeature) theBookIsNotEmpty() error {
	//write content to guarantee it
	_, err := bf.client.book.WriteString(content)
	if err != nil {
		return err
	}
	// reset pos
	bf.client.book.Seek(0, 0)
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {

	//Execute before each scenario. Init client and create dir
	var bookClient BookClient
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		bookClient = BookClient{book: &os.File{}}
		err := os.Mkdir("./"+bookPath, 0755)
		if err != nil {
			fmt.Println(err.Error())
			log.Fatalf("Fail to create directory... aborting")
		}
		return ctx, nil
	})

	bookFeature := BookFeature{client: &bookClient}
	//We want to cleanup our current dir and remove all files created after scenario execution
	ctx.AfterScenario(func(c *godog.Scenario, err error) {
		if info, err := os.Stat(bookPath); err == nil {
			fmt.Println("Cleaning up resources")
			if info.IsDir() {
				err = os.RemoveAll("." + bookPath)
				if err != nil {
					fmt.Printf("Could not cleanup book directory and files")
				}
			}
		}
	})

	ctx.Step(`^I got one line at time$`, bookFeature.iGotOneLineAtTime)
	ctx.Step(`^I open the book$`, bookFeature.iOpenTheBook)
	ctx.Step(`^I read the book by line$`, bookFeature.iReadTheBookByLine)
	ctx.Step(`^i read the book "([^"]*)" fully$`, bookFeature.iReadTheBookFully)
	ctx.Step(`^I write content$`, bookFeature.iWriteContent)
	ctx.Step(`^I write content to unexistent book$`, bookFeature.iWriteContentToUnexistentBook)
	ctx.Step(`^the book content is returned$`, bookFeature.theBookContentIsReturned)
	ctx.Step(`^the book content is updated$`, bookFeature.theBookContentIsUpdated)
	ctx.Step(`^the book "([^"]*)" does not exists$`, bookFeature.theBookDoesNotExists)
	ctx.Step(`^the book is created and updated$`, bookFeature.theBookIsCreatedAndUpdated)
	ctx.Step(`^the file is created empty\.$`, bookFeature.theFileIsCreatedEmpty)
	ctx.Step(`^the file is created$`, bookFeature.theFileIsCreatedEmpty)
	ctx.Step(`^the file is opened correctly$`, bookFeature.theFileIsOpenedCorrectly)
	ctx.Step(`^I create the book "([^"]*)"$`, bookFeature.iCreateTheBook)
	ctx.Step(`^the book "([^"]*)" does not exist$`, bookFeature.theBookDoesNotExist)
	ctx.Step(`^the book "([^"]*)" exists$`, bookFeature.theBookExists)
	ctx.Step(`^the file is created and opened correctly$`, bookFeature.theFileIsCreatedAndOpenedCorrectly)
	ctx.Step(`^I open the book non existent$`, bookFeature.iOpenTheBookNonExistent)
	ctx.Step(`^the book is not empty$`, bookFeature.theBookIsNotEmpty)

}
