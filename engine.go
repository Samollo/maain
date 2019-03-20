package main

import (
	"fmt"

	"github.com/Samollo/maain/constants"
	"github.com/Samollo/maain/crawler"
)

func main() {
	fmt.Println("WEB Search Engine by Ansari & Metadjer")
	/*	args := os.Args[1:]
		c := crawler.NewCrawler(args[0])
		c.Prepare()
		c.CLI.PageRank()
	*/
	cli := crawler.NewCLI(constants.DumpFactor, true)
	cli.AddPage([]int{1, 3})
	cli.AddPage([]int{3})
	cli.AddPage([]int{4})
	cli.AddPage([]int{2, 4})
	cli.AddPage([]int{})

	cli.PageRank()

}
