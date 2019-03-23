package parseutils

import (
	"fmt"
	"sort"
	"strings"
	"unicode"

	"github.com/Samollo/maain/constants"
)

type WordsPagesRelation struct {
	words       []string
	pages       []string
	pagesLength []int
	wordsID     map[string]int
	pagesID     map[string]int
	relations   [][]*Pair
	// not enough, need to check frequency of a word into a page
}

func NewWordPagesRelation(words []string, pagesL []int, pagesName ...string) *WordsPagesRelation {
	wIds := make(map[string]int)
	pIds := make(map[string]int)
	for i, v := range words {
		wIds[v] = i
	}
	if len(pagesName) > 0 {
		for i, v := range pagesName {
			pIds[v] = i
		}
	}
	r := make([][]*Pair, len(words))
	for i := range r {
		r[i] = make([]*Pair, 0)
	}
	return &WordsPagesRelation{
		words:       words,
		pages:       pagesName,
		pagesLength: pagesL,
		wordsID:     wIds,
		pagesID:     pIds,
		relations:   r,
	}
}

func isNumber(w string) bool {
	for _, v := range []rune(w) {
		if unicode.IsDigit(v) == false {
			return false
		}
	}
	return true
}

func (wpr *WordsPagesRelation) Update(PageRank []float64) {
	for index, pages := range wpr.relations {
		//pagesUpdated := make([]*Pair, 0)
		for i, page := range pages {
			if page.Val > 0.00001 {
				pr := PageRank[page.Id]
				if strings.Contains(wpr.pages[page.Id], wpr.words[index]) {
					pr *= constants.PageRankBoost
				}
				if isNumber(wpr.pages[page.Id]) {
					pr /= constants.PageRankBoost
				}
				wpr.relations[index][i].Val = pr
			} else {
				wpr.relations[index][i].Val = -1
			}
		}
		wpr.relations[index] = sortPagesByRank(wpr.relations[index])
	}
}

func (wpr *WordsPagesRelation) Print() {
	for i, v := range wpr.relations {
		fmt.Printf("%v ", wpr.words[i])
		for _, w := range v {
			fmt.Printf("{%v,%v,%v} ", w.Id, wpr.pagesLength[w.Id], w.Val)
		}
		fmt.Println()
	}
}

func sortPagesByRank(pages []*Pair) []*Pair {
	sort.SliceStable(pages, func(i, j int) bool { return pages[i].Val > pages[j].Val })
	return pages
}

func (wpr *WordsPagesRelation) FindPages(word string) []int {
	id := wpr.pagesID[word]
	pages := make([]int, 0)
	for _, v := range wpr.relations[id] {
		pages = append(pages, v.Id)
	}
	return pages
}

func (wpr *WordsPagesRelation) Words() []string {
	return wpr.words
}

func (wpr *WordsPagesRelation) Pages() []string {
	return wpr.pages
}

func (wpr *WordsPagesRelation) WordsID() map[string]int {
	return wpr.wordsID
}

func (wpr *WordsPagesRelation) PagesID() map[string]int {
	return wpr.pagesID
}

func (wpr *WordsPagesRelation) WordByID(id int) string {
	return wpr.words[id]
}

func (wpr *WordsPagesRelation) WordByValue(word string) int {
	return wpr.wordsID[word]
}

func (wpr *WordsPagesRelation) PageByID(id int) string {
	return wpr.pages[id]
}

func (wpr *WordsPagesRelation) PageByValue(title string) int {
	return wpr.pagesID[title]
}

func (wpr *WordsPagesRelation) Relations() [][]*Pair {
	return wpr.relations
}

func (wpr *WordsPagesRelation) AddPage(title string, corpus string, stopWords map[string]int) {
	//if page not already processed, then add it
	if _, ok := wpr.pagesID[title]; !ok {
		wpr.pages = append(wpr.pages, title)
		wpr.pagesID[title] = len(wpr.pages) - 1
	}

	//update the relations
	words, _, _ := extractWords(title, corpus, nil, nil, stopWords)
	for _, w := range words {
		//if word from page is into our hashmap then update the relations
		if index, ok := wpr.wordsID[w.String()]; ok {
			page := &Pair{wpr.pagesID[title], float64(w.Frequence()) / float64(wpr.pagesLength[wpr.pagesID[title]])}
			wpr.relations[index] = append(wpr.relations[index], page)
		}
	}
}
