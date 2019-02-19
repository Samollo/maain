package parseutils

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"
)

var categories = [7]string{"ingenierie", "voiture", "pilot", "moteur", "auto", "mobile", "constructeur"}

func ParseXMLFile(input string, output string) error {
	pageProcessed := 0
	total := 0
	xmlFile, err := os.Open(input)
	if err != nil {
		return err
	}
	outputFile, err := os.Create(output)
	if err != nil {
		return err
	}

	decoder := xml.NewDecoder(xmlFile)
	for {
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
					//	fmt.Println("Titre : " + string(v))
					outputFile.Write([]byte("<title>"))
					outputFile.Write([]byte(v))
					outputFile.Write([]byte("</title>\n"))
				}
				if isOnContent {
					//	fmt.Println("Contenu : " + string(v))
					outputFile.Write([]byte("<text>"))
					outputFile.Write([]byte(v))
					outputFile.Write([]byte("</text>\n"))
				}
			}
		}
	}
	return nil
}

func decode(path string) {
	xmlFile, err := os.Open(path)
	if err != nil {
		return
	}
	decoder := xml.NewDecoder(xmlFile)

	for {
		// Read tokens from the XML document in a stream.
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		// Inspect the type of the token just read.
		switch se := t.(type) {
		case xml.StartElement:
			// If we just read a StartElement token
			// ...and its name is "page"
			if se.Name.Local == "page" {
				var p Page
				// decode a whole chunk of following XML into the
				// variable p which is a Page (se above)
				decoder.DecodeElement(&p, &se)
				// Do some stuff with the page.
			}
		}
	}
}
