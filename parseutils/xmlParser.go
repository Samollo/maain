package parseutils

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

func GenerateDataset(input string, output string, pagesToExtract int, categories []string) error {
	pageProcessed := 0
	total := 0
	xmlFile, err := os.Open(input)
	if err != nil {
		return fmt.Errorf("an error occured. os.Open(%v) in GenerateDaset(): %v", input, err)
	}
	outputFile, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("an error occured. os.Create(%v) failed in GenerateDaset(): %v", output, err)
	}

	decoder := xml.NewDecoder(xmlFile)
	for {
		if pageProcessed == pagesToExtract {
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
				if contains(content, categories) {
					outputFile.WriteString("<title>")
					outputFile.WriteString(title)
					outputFile.WriteString("</title>\n")
					outputFile.WriteString("<text>")
					outputFile.WriteString(content)
					outputFile.WriteString("</text>\n")
					pageProcessed++
				}
			}
		}
	}
	fmt.Printf("%v pages extracted on a total of %v\n", pageProcessed, total)
	return nil
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

func internalLinks(corpus string) (Links, error) {
	l := make(Links, 0)
	regex, err := regexp.Compile("[[].*?[]]]")
	if err != nil {
		return l, fmt.Errorf("error occured while getting internal links: %v", err)
	}
	return regex.FindStringSubmatch("[0-9A-Za-zÀ-ÖØ-öø-ÿ ]+"), nil
}

type WordsPagesRelation struct {
	words     []string
	pages     []string
	wordIDs   map[string]int
	pagesID   map[string]int
	relations [][]int
}

func NewWordPagesRelation(words []string, pages ...string) *WordsPagesRelation {
	wIds := make(map[string]int)
	pIds := make(map[string]int)
	for i, v := range words {
		wIds[v] = i
	}
	if len(pages) > 0 {
		for i, v := range pages {
			pIds[v] = i
		}
	}
	return &WordsPagesRelation{wordIDs: nil, relations: nil}
}

func WordsPagesRelation(words []string, text []string) {
	wpr := make([][]int, len(words))
	for i, w := range wpr {
		w = make([]int, 0)
	}
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
