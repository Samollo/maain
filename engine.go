package main

import (
	"fmt"
	"os"

	"github.com/Samollo/maain/crawler"
	"github.com/Samollo/maain/front"
)

//sylvain.perifel@irif.fr

func main() {
	fmt.Println("WEB Search Engine by Ansari & Metadjer")

	_, err := os.Open("wordpages")
	if err == nil {
		front.LaunchFront()
	} else {
		args := os.Args[1:]
		c := crawler.NewCrawler(args[0])
		c.Prepare()
		front.LaunchFront()
	}
}
