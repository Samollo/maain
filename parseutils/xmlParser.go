package parseutils

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

	"github.com/Samollo/maain/constants"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func GenerateDataset(input string, categories []string) ([]string, []string, []int, error) {
	pageProcessed := 0
	total := 0
	titles := make([]string, 0)
	wordFreq := make([]*Word, 0)
	wordIndex := make(map[string]int)
	stopWords := StopWords()
	pagesLength := make([]int, 0)

	xmlFile, err := os.Open(input)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("an error occured. os.Open(%v) in GenerateDaset(): %v", input, err)
	}

	outputFile, err := os.Create(constants.Output)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("an error occured. os.Create(%v) failed in GenerateDaset(): %v", constants.Output, err)
	}

	decoder := xml.NewDecoder(xmlFile)
	for {
		if pageProcessed == constants.PagesToExtract && !constants.TakeAll {
			break
		}

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
				total++
				title, _ := Extract("title", decoder)
				content, _ := Extract("text", decoder)
				if constants.TakeAll || contains(content, categories) {
					titles = append(titles, title)
					size := 0
					wordFreq, content, size = extractWords(title, content, wordFreq, wordIndex, stopWords)
					pagesLength = append(pagesLength, size)
					serialize(outputFile, title, content)
					pageProcessed++
				}
			}
		}
	}
	fmt.Printf("%v pages extracted on a total of %v read.\n", pageProcessed, total)

	return SortWords(wordFreq), titles, pagesLength, nil
}

func serialize(f *os.File, title string, content string) {
	f.WriteString("<page>\n<title>")
	f.WriteString(title)
	f.WriteString("</title>\n")
	f.WriteString("<text>")
	f.WriteString(content)
	f.WriteString("</text>\n</page>\n")
}

func splitWord(word string) []string {
	array := regexp.MustCompile("[ \n\r\n\t]+").Split(word, -1)
	return array
}

func StopWords() map[string]int {
	hmap := make(map[string]int)
	sw, err := os.Open(constants.Stopwords)
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

//This function extract the Words(value, frequence) and the number of word in the corpus
func extractWords(title, content string, wordFreq []*Word, wordIndex map[string]int, stopWords map[string]int) ([]*Word, string, int) {
	if wordFreq == nil {
		wordFreq = make([]*Word, 0)
	}
	if wordIndex == nil {
		wordIndex = make(map[string]int)
	}
	corpus, _ := FormatWord(content)
	size := 0
	//Iterate through corpus and generate list of Word with their freq
	for _, w := range corpus {
		if _, ok := stopWords[w]; ok {
			continue
		}
		if val, ok := wordIndex[w]; ok {
			wordFreq[val].Increment()
		} else {
			wordFreq = append(wordFreq, NewWord(w))
			wordIndex[w] = len(wordFreq) - 1
		}
		size++
	}
	regex, _ := regexp.Compile("[[].*?[]]]")
	//fmt.Println(len(regex.FindAllString(content, -1)))
	tmp := regex.FindAllString(content, -1)
	//fmt.Println(tmp)
	return wordFreq, strings.Join(tmp, " "), size
}

func SortWords(words []*Word) []string {
	sortedWords := make([]string, 0)

	//Sorted from biggest freq to lowest
	sort.SliceStable(words, func(i, j int) bool { return words[i].freq > words[j].freq })
	//keep only 10k words
	if len(words) > constants.WordsToKeep {
		words = words[:constants.WordsToKeep]
	} else {
		fmt.Printf("not enough words.\n")
	}
	sort.SliceStable(words, func(i, j int) bool { return words[i].value < words[j].value })
	//add to sorted slice of string
	for i := 0; i < len(words); i++ {
		sortedWords = append(sortedWords, words[i].value)
	}

	return sortedWords
}

func contains(text string, categories []string) bool {

	for _, v := range categories {
		if strings.Contains(text, v) {
			return true
		}
	}
	return false
}

type Links []string

func InternalLinks(corpus string, wpr *WordsPagesRelation) ([]int, error) {
	links := make(Links, 0)

	regex, err := regexp.Compile("[[].*?[]]]")
	if err != nil {
		return nil, fmt.Errorf("error occured while getting internal links: %v", err)
	}
	strMatched := regex.FindAllString(corpus, -1) //We extract all the link with the format [[link]]

	re, _ := regexp.Compile("[0-9A-Za-zÀ-ÖØ-öø-ÿ ]+")
	for _, v := range strMatched {
		if toAppend := re.FindAllString(v, -1); len(toAppend) > 0 {
			tmp := re.FindAllString(v, -1)
			for _, v := range tmp {
				links = append(links, v) //We extract all the link within the tag [[...]]

			}
		}
	}

	//If link is not in our dataset, do not keep it.
	inDataset := make([]int, 0)
	for _, link := range links {
		if id, ok := wpr.pagesID[link]; ok {
			inDataset = append(inDataset, id)
		}
	}

	return inDataset, nil
}

func Extract(tag string, decoder *xml.Decoder) (string, error) {
	isOn := false
	result := make([]byte, 0)
	buf := bytes.NewBuffer(result)

	for {
		t, err := decoder.Token()
		if err == io.EOF {
			return "", err
		}
		if err != nil {
			return "", fmt.Errorf("error occured in Extract: %v", err)
		}

		switch v := t.(type) {
		case xml.StartElement:
			if v.Name.Local == tag {
				isOn = true
			}
		case xml.EndElement:
			if v.Name.Local == tag {
				isOn = false
				return buf.String(), nil
			}
		case xml.CharData:
			if isOn {
				xml.EscapeText(buf, v)
			}
		}
	}
	return buf.String(), nil
}

func removeAccents(s string) (string, error) {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, err := transform.String(t, s)
	if err != nil {
		return "", err
	}
	return output, nil
}

func FormatWord(word string) ([]string, error) {
	word, _ = removeAccents(strings.ToLower(word))
	regex, err := regexp.Compile(constants.RegexWord)
	if err != nil {
		return nil, err
	}
	tmp := regex.FindAllString(word, -1)
	return tmp, nil
}

//doCorpus returns a string slice containing words of title and text of an extracted page
//need to format word and check if they are real words
func DoCorpus(title, text string) []string {
	t := splitWord(title)
	a := splitWord(text)
	return append(t, a...)
}

//extractPage returns the title of the page and its content.
//Throw an error if it fails to read the token
func ExtractPage(decoder *xml.Decoder) (string, string, error) {
	title, err := Extract("title", decoder)
	if err == io.EOF {
		return "", "", err
	}
	if err != nil {
		return "", "", fmt.Errorf("error occured in extractPage title: %v", err)
	}
	text, err := Extract("text", decoder)
	if err != nil {
		return "", "", fmt.Errorf("error occured in extractPage text: %v", err)
	}
	return title, text, nil
}
