package crawler

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/Samollo/maain/parseutils"
)

const pagesToExtract = 200000
const output = "dataset.xml"
const wordsToKeep = 10000

var categories = []string{"ingenierie", "voiture", "pilot", "moteur", "auto", "mobile", "constructeur"}

type Crawler struct {
	inputPath      string
	wordDictionary map[string]int
}

type Word struct {
	value string
	freq  int
}

//NewCrawler is a constructor for a basic Crawler struct with a path to the xml file to be processed
//and a wordDictionary containing the n most frequent words
func NewCrawler(path string) *Crawler {
	return &Crawler{inputPath: path, wordDictionary: make(map[string]int)}
}

func (c *Crawler) Prepare() error {
	err := c.dataset()
	if err != nil {
		return fmt.Errorf("error occured in Prepare: %v", err)
	}

	err = c.fillDictionary()
	if err != nil {
		return fmt.Errorf("error occured in Prepare: %v", err)
	}
	return nil
}

func (c *Crawler) dataset() error {
	return parseutils.GenerateDataset(c.inputPath, output, pagesToExtract, categories)
}

func (c *Crawler) fillDictionary() error {
	file, err := os.Open(output)
	if err != nil {
		return fmt.Errorf("error occured in fillDictionary: %v", err)
	}
	decoder := xml.NewDecoder(file)
	for {
		title, text, err := extractPage(decoder)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

	}
	return nil
}

func extractPage(decoder *xml.Decoder) (string, string, error) {
	title, err := parseutils.Extract("title", decoder)
	if err != nil {
		return "", "", fmt.Errorf("error occured in extractPage: %v", err)
	}
	text, err := parseutils.Extract("text", decoder)
	if err != nil {
		return "", "", fmt.Errorf("error occured in extractPage: %v", err)
	}
	return title, text, nil
}
