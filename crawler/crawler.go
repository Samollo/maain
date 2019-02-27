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
	wordDictionary []string
}

type Word struct {
	value string
	freq  int
}

//NewCrawler is a constructor for a basic Crawler struct with a path to the xml file to be processed
//and a wordDictionary containing the n most frequent words
func NewCrawler(path string) *Crawler {
	return &Crawler{inputPath: path, wordDictionary: make([]string, 0)}
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

	wordIndex := make(map[string]int)
	wordFreq := make([]*Word, 0)
	decoder := xml.NewDecoder(file)
	for {
		title, text, err := extractPage(decoder)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error occured in fillDictionary: %v", err)
		}

		corpus := doCorpus(title, text)
		//Iterate through corpus and generate list of Word with their freq
		for _, word := range corpus {
			if val, ok := wordIndex[word]; ok {
				wordFreq[val].freq++
			} else {
				wordFreq = append(wordFreq, &Word{value: word, freq: 1})
				wordIndex[word] = len(wordFreq) - 1
			}
		}
	}

	//Sorted from biggest freq to lowest
	sort.SliceStable(wordFreq, func(i, j int) bool { return wordFreq[i].freq > wordFreq[j].freq })
	wordFreq = wordFreq[:wordsToKeep]
	sort.SliceStable(wordFreq, func(i, j int) bool { return wordFreq[i].value < wordFreq[j].value })
	//add to dico
	for i := 0; i < len(wordFreq); i++ {
		c.wordDictionary = append(c.wordDictionary, wordFreq[i].value)
	}
	return nil
}

//doCorpus returns a string slice containing words of title and text of an extracted page
func doCorpus(title, text string) []string {
	corpus := strings.Split(title, " ")
	tmp := strings.Replace(text, "\n", " ", -1)
	return append(corpus, strings.Split(tmp, " ")...)
}

//extractPage returns the title of the page and its content.
//Throw an error if it fails to read the token
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
