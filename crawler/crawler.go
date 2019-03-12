package crawler

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"
	"unicode"

	"github.com/Samollo/maain/parseutils"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

const pagesToExtract = 200000
const output = "dataset.xml"
const stopwords = "stopwords-fr.txt"
const wordsToKeep = 10000

var categories = []string{"ingenierie", "voiture", "pilot", "moteur", "auto", "mobile", "constructeur"}

type Crawler struct {
	inputPath      string
	wordDictionary []string
	stopWords      map[string]int
}

type Word struct {
	value string
	freq  int
}

//NewCrawler is a constructor for a basic Crawler struct with a path to the xml file to be processed
//and a wordDictionary containing the n most frequent words
func NewCrawler(path string) *Crawler {
	return &Crawler{inputPath: path,
		wordDictionary: make([]string, 0),
		stopWords:      stopWords(),
	}
}

func stopWords() map[string]int {
	hmap := make(map[string]int)
	sw, err := os.Open(stopwords)
	if err != nil {
		fmt.Println("Failed to open stopwords file.")
		return hmap
	}
	s, err := ioutil.ReadAll(sw)
	if err != nil {
		fmt.Println("Failed top read stopwords file.")
		return hmap
	}

	list := strings.Split(bytes.NewBuffer(s).String(), "\n")
	for _, w := range list {
		if _, ok := hmap[w]; !ok {
			hmap[w] = 1
		}
	}
	return hmap
}

//Prepare generates dataset and fill dictionary before CLI
func (c *Crawler) Prepare() error {
	fmt.Println("Prepare..")
	err := c.dataset()
	if err != nil {
		return fmt.Errorf("error occured in Prepare: %v", err)
	}

	err = c.fillDictionary()
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("error occured in Prepare: %v", err)
	}
	return nil
}

func (c *Crawler) dataset() error {
	fmt.Println("Dataset..")
	return parseutils.GenerateDataset(c.inputPath, output, pagesToExtract, categories)
}

func (c *Crawler) fillDictionary() error {
	fmt.Println("fill..")
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
			w, err := formatWord(word)
			if _, ok := c.stopWords[w]; ok || err != nil {
				continue
			}
			if val, ok := wordIndex[w]; ok {
				wordFreq[val].freq++
			} else {
				wordFreq = append(wordFreq, &Word{value: w, freq: 1})
				wordIndex[word] = len(wordFreq) - 1
			}
		}
	}

	//Sorted from biggest freq to lowest
	sort.SliceStable(wordFreq, func(i, j int) bool { return wordFreq[i].freq > wordFreq[j].freq })
	//keep only 10k words
	if len(wordFreq) > wordsToKeep {
		wordFreq = wordFreq[:wordsToKeep]
	}
	sort.SliceStable(wordFreq, func(i, j int) bool { return wordFreq[i].value < wordFreq[j].value })
	//add to dico and
	//remove accents and upper-cases
	for i := 0; i < len(wordFreq); i++ {
		c.wordDictionary = append(c.wordDictionary, wordFreq[i].value)
		fmt.Printf("[%s]\n", c.wordDictionary[len(c.wordDictionary)-1])
	}
	return nil
}

func removeAccents(s string) (string, error) {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, err := transform.String(t, s)
	if err != nil {
		return "", err
	}
	return output, nil
}

func formatWord(word string) (string, error) {
	word, _ = removeAccents(strings.ToLower(word))
	regex, err := regexp.Compile("[A-Za-zÀ-ÖØ-öø-ÿ]+")
	if err != nil {
		return "", err
	}
	tmp := regex.FindString(word)
	if tmp == "" {
		return "", fmt.Errorf("no word found.")
	}
	return tmp, nil
}

//doCorpus returns a string slice containing words of title and text of an extracted page
//need to format word and check if they are real words
func doCorpus(title, text string) []string {
	corpus := strings.Split(title, " ")
	tmp := strings.Replace(text, "\n", " ", -1)
	return append(corpus, strings.Split(tmp, " ")...)
}

//extractPage returns the title of the page and its content.
//Throw an error if it fails to read the token
func extractPage(decoder *xml.Decoder) (string, string, error) {
	title, err := parseutils.Extract("title", decoder)
	if err == io.EOF {
		return "", "", err
	}
	if err != nil {
		return "", "", fmt.Errorf("error occured in extractPage title: %v", err)
	}
	text, err := parseutils.Extract("text", decoder)
	if err != nil {
		return "", "", fmt.Errorf("error occured in extractPage text: %v", err)
	}
	return title, text, nil
}
