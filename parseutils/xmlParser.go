package parseutils

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
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

func Extract(tag string, decoder *xml.Decoder) (string, error) {
	isOn := false
	result := ""

	for {
		t, err := decoder.Token()
		if err != nil {
			fmt.Printf("decoder.Token() failed with '%s'\n", err)
			return "", err
		}

		switch v := t.(type) {
		case xml.StartElement:
			if v.Name.Local == tag {
				isOn = true
			}
		case xml.EndElement:
			if v.Name.Local == tag {
				isOn = false
				return result, nil
			}
		case xml.CharData:
			if isOn {
				result = string(v[:])
			}
		}
	}
	return result, nil
}
