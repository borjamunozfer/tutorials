package books

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/cucumber/godog"
)

type BookFeature struct {
	author *Author
	file   *os.File
	body   string
}

const filenameKey = "bookFilename"

func (bk *BookFeature) aNewEmptyBookIsCreated(ctx context.Context) error {
	//recover the filename of the book created from Context
	bookFilename := ctx.Value(filenameKey).(string)

	//check book with that filename exists.
	fileInfo, err := os.Stat(bookPath + bookFilename)
	if err != nil {
		return err
	}
	//check bookfile size is 0 (empty)
	if fileInfo.Size() != 0 {
		return fmt.Errorf("Book file exists but its not empty")
	}

	return nil
}

func (bk *BookFeature) iCreateTheBook(ctx context.Context, bookTitle string) (context.Context, error) {

	if err := bk.author.CreateBook(bookTitle); err != nil {
		return ctx, err
	}

	//save filename to context
	return context.WithValue(ctx, filenameKey, formatBookFilename(bookTitle)), nil
}

func (bf *BookFeature) iDeleteLine(pos string) error {

	//a lot of different approaches to this
	//declare buffer
	numPos, _ := strconv.Atoi(pos)
	scanner := bufio.NewScanner(bf.file)
	start := numPos
	//read by pos
	scanLines := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		advance, token, err = bufio.ScanLines(data, atEOF)
		start += int(advance)
		return
	}
	scanner.Split(scanLines)
	for scanner.Scan() {
		fmt.Printf("Pos: %d, Scanned: %s\n", start, scanner.Text())
	}
	return nil
}

func (bf *BookFeature) iOpenTheBook(bookTitle string) error {

	file, err := os.OpenFile(bookPath+formatBookFilename(bookTitle), os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	bf.file = file
	return nil
}

func (bf *BookFeature) iWriteLines() error {

	content := `Háblame, Musa, del varón de gran ingenio, que anduvo errante muchísimo tiempo.
	y que vio las ciudades de muchos hombres y conoció su manera de pensar,
	pero padeció aún en el mar muchos dolores en su ánimo, 
	procurando conservar su vida y el regreso de sus compañeros. 
	Mas ni siquiera así terminó de salvar a sus compañeros, aunque lo deseaba vivamente, pues perecerieron por sus propios actos temerarios.`
	err := bf.author.WriteBook(bf.file, content)
	if err != nil {
		return err
	}
	bf.body = content
	return nil
}

func (bf *BookFeature) theBookDoesNotExist(bookTitle string) error {

	// two different ways to check book does not exist
	// 1. find by key in author books map
	// 2. parse current dir, search filename.
	if _, ok := bf.author.Books[bookTitle]; ok {
		return fmt.Errorf("Book %s already exists", bookTitle)
	}

	if _, err := os.Stat(bookPath + bookTitle + ".txt"); os.IsExist(err) {
		return err
	}

	return nil
}

func theBookExists(arg1 string) error {
	return nil
}

func (bf *BookFeature) theBookIsUpdated(bookFile string) error {

	// Read buffered
	/*defer bf.file.Close()
	b1 := make([]byte, 5)
	n1, err := bf.file.Read(b1)
	if err != nil {
		return err
	}
	fmt.Printf("%d bytes readed: %s\n", n1, string(b1[:n1]))*/
	bcontent, err := os.ReadFile(bookPath + formatBookFilename(bookFile))
	if err != nil {
		return err
	}

	if string(bcontent) != bf.body {
		return fmt.Errorf("Got %s content, expected %s", string(bcontent), bf.body)
	}
	return nil
}

func theBookLineIsRemoved() error {
	return godog.ErrPending
}

func InitializeScenario(ctx *godog.ScenarioContext) {

	var author Author
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		fmt.Println("Initialize author struct")
		author = NewAuthor("Homero")
		return ctx, nil
	})

	//We want to cleanup our current dir and remove all files created after scenario execution
	ctx.AfterScenario(func(c *godog.Scenario, err error) {
		if info, err := os.Stat(bookPath); err == nil {
			fmt.Println("Cleaning up resources")
			if info.IsDir() {
				time.Sleep(time.Second + 10)
				err = os.RemoveAll("." + bookPath)
				if err != nil {
					fmt.Printf("Could not cleanup book directory and files")
				}
			}
		}
	})

	bookFeature := &BookFeature{author: &author}
	ctx.Step(`^a new empty book is created$`, bookFeature.aNewEmptyBookIsCreated)
	ctx.Step(`^I create the book "([^"]*)"$`, bookFeature.iCreateTheBook)
	ctx.Step(`^I delete "([^"]*)" line$`, bookFeature.iDeleteLine)
	ctx.Step(`^I open the book "([^"]*)"$`, bookFeature.iOpenTheBook)
	ctx.Step(`^I write lines$`, bookFeature.iWriteLines)
	ctx.Step(`^the book "([^"]*)" does not exist$`, bookFeature.theBookDoesNotExist)
	ctx.Step(`^the book "([^"]*)" exists$`, theBookExists)
	ctx.Step(`^the book "([^"]*)" is updated$`, bookFeature.theBookIsUpdated)
	ctx.Step(`^the book line is removed$`, theBookLineIsRemoved)
}
