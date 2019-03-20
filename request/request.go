package request

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strings"

	"github.com/Samollo/maain/parseutils"
)

type Request struct {
	stopWords        map[string]int
	wordPageRelation map[string][]string
}

func InitializeRequest(filePath string) Request {
	r := Request{make(map[string]int, 0), make(map[string][]string, 0)}
	r.parseRelationFile(filePath)
	return r
}

//NextLine returns the next line read in the associated file, or an error if EOF
func NextLine(reader *bufio.Reader) ([]byte, error) {
	token := make([]byte, 0)
	for {
		t, isPrefix, err := reader.ReadLine()
		if err != nil {
			return nil, fmt.Errorf("Error in GetLine() of: %v", err)
		}
		token = append(token, t...)
		if !isPrefix {
			break
		}
	}
	return token, nil
}

func (r *Request) parseRelationFile(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		token, err := NextLine(reader)
		if err != nil {
			break
		}
		k, v := parseLine(string(token))
		r.wordPageRelation[k] = v
	}
	//fmt.Printf("%v\n", r.wordPageRelation)
}

func parseLine(line string) (string, []string) {
	value := strings.Split(line, ":")
	word := value[0]
	return word, strings.Split(value[1], ",")
}

func (r *Request) ReturnFoundPages(sentence string) []string {
	words, _ := parseutils.FormatWord(sentence)
	minLength := math.MaxInt64
	minSlice := 0
	pagesConcern := make([][]string, 0)
	isChange := false

	cmp := 0
	for _, v := range words {
		if value, ok := r.wordPageRelation[v]; ok {
			pagesConcern = append(pagesConcern, value)
			if len(value) < minLength {
				isChange = true
				minLength = len(value)
				minSlice = cmp
			}
			cmp++
		}
	}
	if !isChange {
		return make([]string, 0)
	}

	minSizePage := make(map[string]int)
	for _, pages := range pagesConcern {
		for _, page := range pages {
			if _, ok := minSizePage[page]; !ok {
				minSizePage[page] = 1
			} else {
				minSizePage[page]++
			}
		}
	}
	mvpPage := make([]string, 0)
	for _, v := range pagesConcern[minSlice] {
		if minSizePage[v] == len(words) {
			mvpPage = append(mvpPage, v)
		}
	}
	return mvpPage
}
