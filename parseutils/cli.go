package parseutils

import (
	"fmt"

	"github.com/Samollo/maain/constants"
)

type CLI struct {
	c []float32
	l []int
	i []int
}

func NewCLI() *CLI {
	return &CLI{
		c: make([]float32, 0),
		l: make([]int, 0, constants.PagesToExtract),
		i: make([]int, 0),
	}
}

func (cli *CLI) C() []float32 {
	return cli.c
}

func (cli *CLI) L() []int {
	return cli.l
}
func (cli *CLI) I() []int {
	return cli.i
}

func (cli *CLI) AddPage(pageId int, links []int) error {
	coef := float32(1) / float32(len(links))

	for _, value := range links {
		cli.c = append(cli.c, coef)
		cli.i = append(cli.i, value)
	}

	nbLink := 0
	if len(cli.l) == 0 {
		cli.l = append(cli.l, nbLink)
	}
	nbLink = cli.l[len(cli.l)-1] + len(links)

	cli.l = append(cli.l, nbLink)
	return nil
}

func (cli *CLI) transposer(v []float32) []float32 {
	n := len(cli.l) - 1
	result := make([]float32, n)

	for i := 0; i < n; i++ {
		for j := cli.l[i]; j < cli.l[i+1]-1; j++ {
			result[cli.i[j]] += cli.c[j] * v[i]
		}
	}

	return result
}

func (cli *CLI) PageRank(id int, step int) {
	n := len(cli.l)
	P := make([]float32, n)
	P[id] = 1

	for i := 0; i < step; i++ {
		PK := cli.transposer(P)
		//todo verification
		P = PK
		fmt.Printf("Probality of Page %v: %v", i, P[i])
	}
}
