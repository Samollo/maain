package main

import (
	"fmt"
	"os"

	"github.com/Samollo/maain/crawler"
	"github.com/Samollo/maain/parseutils"
)

func main() {
	fmt.Println("WEB Search Engine by Ansari & Metadjer")
	args := os.Args[1:]
	_ = crawler.NewCrawler(args[0])
	parseutils.ParseXMLFile(args[0], "output.xml")
}
