package crawler

import (
	"encoding/xml"
	"io"
	"os"

	"github.com/Samollo/maain/parseutils"
)

const pagesToExtract = 200000
const output = "dataset.xml"

var categories = []string{"ingenierie", "voiture", "pilot", "moteur", "auto", "mobile", "constructeur"}

type Crawler struct {
	inputPath      string
	wordDictionary map[string]int
}

func NewCrawler(path string) *Crawler {
	return &Crawler{inputPath: path, wordDictionary: make(map[string]int)}
}

func (c *Crawler) Dataset() error {
	return parseutils.GenerateDataset(c.inputPath, output, pagesToExtract, categories)
}

func (c *Crawler) fillDictionary() error {
	file, err := os.Open(output)
	if err != nil {
		return err
	}
	decoder := xml.NewDecoder(file)
	for {
		_, _, err := extractPage(decoder)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

	}
	return nil
}
