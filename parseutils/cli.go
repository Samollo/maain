package parseutils

import (
	"fmt"
	"math"

	"github.com/Samollo/maain/constants"
)

type CLI struct {
	c []float64
	l []int
	i []int
}

func NewCLI() *CLI {
	return &CLI{
		c: make([]float64, 0),
		l: make([]int, 0, constants.PagesToExtract),
		i: make([]int, 0),
	}
}

func (cli *CLI) C() []float64 {
	return cli.c
}

func (cli *CLI) L() []int {
	return cli.l
}
func (cli *CLI) I() []int {
	return cli.i
}

func (cli *CLI) AddPage(links []int) error {
	coef := float64(1) / float64(len(links))

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

func (cli *CLI) transposer(v []float64) []float64 {
	n := len(cli.l) - 1
	result := make([]float64, n)

	for i := 0; i < n; i++ {
		for j := cli.l[i]; j < cli.l[i+1]; j++ {
			result[cli.i[j]] += cli.c[j] * v[i]
		}
	}

	return result
}

func (cli *CLI) PageRank(id int) {
	n := len(cli.l) - 1

	P := make([]float64, n)
	P[id] = 1
	epsilon := 0.0004
	var delta float64
	delta = 1
	for delta > epsilon {
		PK := cli.transposer(P)

		//	fmt.Println(len(P))
		//	fmt.Println(len(PK))
		delta = math.Abs(sum(P, PK))
		//	fmt.Printf("delta: %v\n", delta)
		//	fmt.Printf("Probality of P: %v\n", P)
		//	fmt.Printf("Probality of PK: %v\n", PK)
		P = PK
	}
	fmt.Printf("Probability of vector P(%v) is %v\n", id, P)
}

func sum(a, b []float64) float64 {
	delta := 0.0
	for i, _ := range b {
		delta += b[i] - a[i]
	}
	return delta
}
