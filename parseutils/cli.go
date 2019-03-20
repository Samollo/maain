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

func (cli *CLI) transposer(v []float64, zap bool) ([]float64, float64) {
	SumOfPR := 0.0
	n := len(cli.l) - 1
	result := make([]float64, n)

	for i := 0; i < n; i++ {
		for j := cli.l[i]; j < cli.l[i+1]; j++ {
			if zap {
				result[cli.i[j]] += cli.c[j] * v[i] * constants.DumpFactor
			} else {
				result[cli.i[j]] += cli.c[j] * v[i]
			}
		}
	}
	for i := 0; i < n; i++ {
		SumOfPR += result[i]
	}

	return result, SumOfPR
}

func (cli *CLI) PageRank(id ...int) {
	fmt.Println("PageRank..")
	n := len(cli.l) - 1

	P := make([]float64, n)

	if len(id) > 0 {
		P[id[0]] = 1
	} else {
		for i := range P {
			P[i] = float64(1) / float64(n)
		}
	}

	//	epsilon := 0.0001
	//	var delta float64
	//	delta = 1

	compt := 1

	for i := 0; i < 3; i++ {
		//	for delta > epsilon {
		PK, SumOfPR := cli.transposer(P, true)
		_ = sum(P, PK, SumOfPR)
		P = PK
		fmt.Printf("---%v tours---\n", compt)
		fmt.Printf("Probability of vector P is %v\n\n", P)
		compt++
	}
}

func sum(a, b []float64, SumOfPR float64) float64 {
	delta := 0.0
	for i := range b {
		b[i] = b[i] + ((1 - SumOfPR) / float64(len(b)))
		delta += b[i] - a[i]
	}
	return math.Abs(delta)
}
