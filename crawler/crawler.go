package crawler

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"

	"github.com/Samollo/maain/constants"
	"github.com/Samollo/maain/parseutils"
)

var categories = []string{"ingenierie", "voiture", "pilot", "moteur", "auto", "mobile", "constructeur"}

type Crawler struct {
	inputPath      string
	wordDictionary []string
	wpr            *parseutils.WordsPagesRelation
	CLI            *CLI
}

//NewCrawler is a constructor for a basic Crawler struct with a path to the xml file to be processed
//and a wordDictionary containing the n most frequent words
func NewCrawler(path string) *Crawler {
	return &Crawler{
		inputPath:      path,
		wordDictionary: make([]string, 0),
		wpr:            nil,
		CLI:            NewCLI(constants.DumpFactor, constants.Zap),
	}
}

//Prepare generates dataset and fill dictionary before CLI
func (c *Crawler) Prepare() error {
	fmt.Println("Prepare..")
	words, titles, pagesLength, err := c.dataset()
	if err != nil {
		return fmt.Errorf("error occured in Prepare: %v", err)
	}

	c.wordDictionary = words
	c.wpr = parseutils.NewWordPagesRelation(c.wordDictionary, pagesLength, titles...)
	err = c.cliRelation()
	if err != nil {
		fmt.Println("error")
		return err
	}

	pageRanks := c.CLI.PageRank()

	c.wpr.Update(pageRanks)
	c.serialize()

	return err
}

func (c *Crawler) dataset() ([]string, []string, []int, error) {
	fmt.Println("Dataset..")
	return parseutils.GenerateDataset(c.inputPath, categories)
}

func (c *Crawler) cliRelation() error {
	fmt.Println("CLI and Word-Pages relation...")
	stopWords := parseutils.StopWords()
	xmlFile, err := os.Open(constants.Output)
	if err != nil {
		return fmt.Errorf("an error occured. os.Open: %v", err)
	}

	decoder := xml.NewDecoder(xmlFile)
	for {
		t, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("decoder.Token() failed with '%s'\n", err)
			break
		}

		switch v := t.(type) {
		case xml.StartElement:
			if v.Name.Local == "page" {
				title, _ := parseutils.Extract("title", decoder)
				content, _ := parseutils.Extract("text", decoder)
				ids, err := parseutils.InternalLinks(content, c.wpr)
				if err != nil {
					return err
				}
				c.wpr.AddPage(title, content, stopWords)
				err = c.CLI.AddPage(ids)
				if err != nil {
					fmt.Println(err)
					return err
				}
			}
		}
	}
	return nil
}

func (c *Crawler) serialize() error {
	fmt.Println("Serialize...")
	f, err := os.Create(constants.ResultFile)
	if err != nil {
		return err
	}

	relations := c.wpr.Relations()

	for i, v := range relations {
		f.WriteString(c.wpr.WordByID(i))
		f.WriteString(":")
		for i, w := range v {
			f.WriteString(c.wpr.PageByID(w.Id))
			if i < len(v)-1 {
				f.WriteString(",")
			}
		}
		f.WriteString("\n")
	}
	return nil
}

//mot:pages,pages
