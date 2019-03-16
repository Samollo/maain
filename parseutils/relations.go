package parseutils

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
	return &WordsPagesRelation{words: words, pages: pagesName, wordsID: wIds, pagesID: pIds, relations: nil}
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

func (wpr *WordsPagesRelation) addPage(title string, corpus string) {
	//if page not already processed, then add it
	if _, ok := wpr.pagesID[title]; !ok {
		wpr.pages = append(wpr.pages, title)
		wpr.pagesID[title] = len(wpr.pages) - 1
	}

	//update the relations
	words, total := extractWords(title, corpus, nil, nil)
	for _, w := range words {
		//if word from page is into our hashmap then update the relations
		if index, ok := wpr.wordsID[w.String()]; !ok {
			page := &Pair{wpr.pagesID[title], w.Frequence() / total}
			wpr.relations[index] = append(wpr.relations[index], page)
			wpr.pages = append(wpr.pages, title)
			wpr.pagesID[title] = len(wpr.pages) - 1
		}
	}
}
