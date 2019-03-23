package request

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
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

func (r *Request) Intersection(sentence string) []string {
	words, _ := parseutils.FormatWord(sentence)
	bestResults := make(map[string]float64)
	pagesFound := make([][]string, 0)
	totalPagesFound := 0

	//we get all the pages concerned by the request
	for _, v := range words {
		if value, ok := r.wordPageRelation[v]; ok {
			pagesFound = append(pagesFound, value)
			totalPagesFound += len(value)
		}
	}

	//We iterate through the pages found and give them a score
	//based on how many words do they have and their position in the pageranking
	for _, pages := range pagesFound {
		for rank, page := range pages {
			if _, ok := bestResults[page]; ok {
				bestResults[page] += 1 / float64(totalPagesFound)
			} else {
				bestResults[page] = (1 / float64(totalPagesFound)) + (float64(totalPagesFound)-float64(rank))/float64(totalPagesFound)
			}
		}
	}

	pagesScore := mapToSlice(bestResults)
	return SortWords(pagesScore)
}

func SortWords(words []*parseutils.Word) []string {
	sortedWords := make([]string, 0)

	//Sorted from biggest freq to lowest
	sort.SliceStable(words, func(i, j int) bool { return words[i].Frequence() > words[j].Frequence() })
	//keep only 10k words
	//if len(words) > constants.WordsToKeep {
	//	words = words[:constants.WordsToKeep]
	//} else {
	//	fmt.Printf("not enough words.\n")
	//}
	for i := 0; i < len(words); i++ {
		sortedWords = append(sortedWords, words[i].String())
	}
	return sortedWords
}

func mapToSlice(m map[string]float64) []*parseutils.Word {
	pairs := make([]*parseutils.Word, len(m))
	index := 0
	for i, v := range m {
		pairs[index] = parseutils.NewWord(i, int(v)*100)
		index++
	}
	return pairs
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
