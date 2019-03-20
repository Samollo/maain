package parseutils

import "sort"

type WordsPagesRelation struct {
	words     []string
	pages     []string
	wordsID   map[string]int
	pagesID   map[string]int
	relations [][]*Pair
	// not enough, need to check frequency of a word into a page
}

func NewWordPagesRelation(words []string, pagesName ...string) *WordsPagesRelation {
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
	return &WordsPagesRelation{
		words:     words,
		pages:     pagesName,
		wordsID:   wIds,
		pagesID:   pIds,
		relations: nil,
	}
}

func (wpr *WordsPagesRelation) update(PageRank []float64) {
	for index, pages := range wpr.relations {
		pagesUpdated := make([]*Pair, 0, len(pages))
		for _, page := range pages {
			if page.val > 0.01 {
				pagesUpdated = append(pagesUpdated, &Pair{page.id, PageRank[page.id]})
			}
		}
		wpr.relations[index] = sortPagesByRank(pagesUpdated)
	}
}

func sortPagesByRank(pages []*Pair) []*Pair {
	sort.SliceStable(pages, func(i, j int) bool { return pages[i].val > pages[j].val })
	return pages
}

func (wpr *WordsPagesRelation) FindPages(word string) []int {
	id := wpr.pagesID[word]
	pages := make([]int, 0)
	for _, v := range wpr.relations[id] {
		pages = append(pages, v.id)
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

func (wpr *WordsPagesRelation) AddPage(title string, corpus string, stopWords map[string]int) {
	//if page not already processed, then add it
	if _, ok := wpr.pagesID[title]; !ok {
		wpr.pages = append(wpr.pages, title)
		wpr.pagesID[title] = len(wpr.pages) - 1
	}

	//update the relations
	words, content := extractWords(title, corpus, nil, nil, stopWords)
	for _, w := range words {
		//if word from page is into our hashmap then update the relations
		if index, ok := wpr.wordsID[w.String()]; !ok {
			page := &Pair{wpr.pagesID[title], float64(w.Frequence()) / float64(len(content))}
			wpr.relations[index] = append(wpr.relations[index], page)
			wpr.pages = append(wpr.pages, title)
			wpr.pagesID[title] = len(wpr.pages) - 1
		}
	}
}
