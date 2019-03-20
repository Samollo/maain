package main

import (
	"fmt"
	"os"

	"github.com/Samollo/maain/crawler"
)

func main() {
	fmt.Println("WEB Search Engine by Ansari & Metadjer")
	args := os.Args[1:]
	c := crawler.NewCrawler(args[0])
	c.Prepare()
	c.CLI.PageRank()
	//	cli := parseutils.NewCLI()
	//	cli.AddPage(0, []int{1, 3})
	//	cli.AddPage(1, []int{2})
	//	cli.AddPage(2, []int{3})
	//	cli.AddPage(3, []int{0, 1, 2})

	//	cli.PageRank(0)
	//	cli.PageRank(1)
	//	cli.PageRank(2)
	//	cli.PageRank(3)

	//	cli.PageRank()

}
