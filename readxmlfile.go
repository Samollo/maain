package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

type Mediawiki struct {
	XMLName xml.Name `xml:"mediawiki"`
	Pages   []Page   `xml:"page"`
}

type Page struct {
	XMLName xml.Name `xml:"page"`
	Title   string   `xml:"title"`
	Text    string   `xml:"text"`
}

func main() {
	xmlFile, err := os.Open("~/Downloads/frwiki-debut.xml")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully Opened frwiki-debut.xml")
	defer xmlFile.Close()

	byteValue, _ := ioutil.ReadAll(xmlFile)

	var mediawiki Mediawiki
	xml.Unmarshal(byteValue, &mediawiki)

	for i := 0; i < len(mediawiki.Pages); i++ {
		fmt.Println("Title:", mediawiki.Pages[i].Title)
		fmt.Println("Text:", mediawiki.Pages[i].Text)
	}

}
