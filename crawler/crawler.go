package crawler

import (
	"fmt"

	"github.com/Samollo/maain/parseutils"
)

var categories = []string{"ingenierie", "voiture", "pilot", "moteur", "auto", "mobile", "constructeur"}

type Crawler struct {
	inputPath      string
	wordDictionary []string
	wpr            *parseutils.WordsPagesRelation
}

//NewCrawler is a constructor for a basic Crawler struct with a path to the xml file to be processed
//and a wordDictionary containing the n most frequent words
func NewCrawler(path string) *Crawler {
	return &Crawler{inputPath: path,
		wordDictionary: make([]string, 0),
		wpr:            nil,
	}
}

//Prepare generates dataset and fill dictionary before CLI
func (c *Crawler) Prepare() error {
	fmt.Println("Prepare..")
	words, titles, err := c.dataset()
	if err != nil {
		return fmt.Errorf("error occured in Prepare: %v", err)
	}

	c.wordDictionary = words
	c.wpr = parseutils.NewWordPagesRelation(c.wordDictionary, titles...)
	return nil
}

func (c *Crawler) dataset() ([]string, []string, error) {
	fmt.Println("Dataset..")
	return parseutils.GenerateDataset(c.inputPath, categories)
}
