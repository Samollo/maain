package main

import (
	"fmt"

	"github.com/Samollo/maain/parseutils"
)

func main() {
	fmt.Println("WEB Search Engine by Ansari & Metadjer")
	//args := os.Args[1:]
	//c := crawler.NewCrawler(args[0])
	//c.Prepare()
	cli := parseutils.NewCLI()
	cli.AddPage([]int{1, 3})
	cli.AddPage([]int{2})
	cli.AddPage([]int{3})
	cli.AddPage([]int{0, 1, 2})

	cli.PageRank(0)
	cli.PageRank(1)
	cli.PageRank(2)
	cli.PageRank(3)

}
