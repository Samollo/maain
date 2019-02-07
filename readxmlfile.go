package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

type XMLContent struct {
	Pages []Page `xml:"page"`
}

type Page struct {
	Title string `xml:"title"`
	Text  string `xml:"text"`
}

func main() {
	/*
	xmlFile, err := os.Open("/Users/mohammed/Downloads/frwiki-20190120-pages-articles.xml")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully Opened frwiki-debut.xml")
	defer xmlFile.Close()

	byteValue, _ := ioutil.ReadAll(xmlFile)

	var content XMLContent
	xml.Unmarshal(byteValue, &content)

	for i := 0; i < len(content.Pages); i++ {
		fmt.Println("Title:", content.Pages[i].Title)
		fmt.Println("Text:", content.Pages[i].Text)
	}*/

	ParseXMLFile(os.Args[1], "")

}

func ParseXMLFile(filePath string, outputPath string) error {
	xmlFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	isOnPage := false
	isOnTitle := false
	isOnContent := false

	forloopTime := 1

	decoder := xml.NewDecoder(xmlFile)
	for {
		i := 0

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
				isOnPage = true
			}
			if v.Name.Local == "title" {
				isOnTitle = true
			}
			if v.Name.Local == "text" {
				isOnContent = true
			}

		case xml.EndElement:
			if v.Name.Local == "page" {
				isOnPage = false
				i++
			}
			if v.Name.Local == "title" {
				isOnTitle = false
			}
			if v.Name.Local == "text" {
				isOnContent = false
			}

		case xml.CharData:
			if isOnPage {
				if isOnTitle {
					fmt.Println("Titre : " + string(v))
				}
				if isOnContent {
					fmt.Println("Contenu : " + string(v))
				}
			}
		}

		if i >= forloopTime {
			break
		}
	}

	return nil
}
