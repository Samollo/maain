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
				total++
				title := extract("title", decoder)
				content := extract("text", decoder)
				if contains(content) {
					outputFile.WriteString(title)
					outputFile.WriteString(content)
					pageProcessed++
				}
			}
		}
	}

	fmt.Printf("Nombre de pages extraites: %v sur un total de %v\n", pageProcessed, total)
	return nil
}

func contains(text string) bool {
	for _, v := range categories {
		if strings.Contains(text, v) {
			return true
		}
	}
	return false
}

func extract(tag string, decoder *xml.Decoder) string {
	isOn := false
	result := ""

	for {
		t, err := decoder.Token()
		if err != nil {
			fmt.Printf("decoder.Token() failed with '%s'\n", err)
			break
		}

		switch v := t.(type) {
		case xml.StartElement:
			if v.Name.Local == tag {
				isOn = true
			}
		case xml.EndElement:
			if v.Name.Local == tag {
				isOn = false
				return result
			}
		case xml.CharData:
			if isOn {
				result = string(v[:])
			}
		}
	}
	return result
}
