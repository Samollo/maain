package crawler

import (
	"fmt"
	"math"

	"github.com/Samollo/maain/constants"
)

type CLI struct {
	c          []float64
	l          []int
	i          []int
	dumpFactor float64
	zap        bool
}

func NewCLI(beta float64, z bool) *CLI {
	return &CLI{
		c:          make([]float64, 0),
		l:          make([]int, 0, constants.PagesToExtract),
		i:          make([]int, 0),
		dumpFactor: beta,
		zap:        z,
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

func (cli *CLI) PageRank(id ...int) []float64 {
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

	epsilon := constants.Epsilon
	delta := 1.0
	count := 1
	for delta > epsilon {
		PK := cli.transpose(P, true)
		sum := 0.0
		if cli.zap {
			sum = cli.pageRankSum(PK)
		}
		delta = cli.updateDelta(P, PK, sum)
		//	fmt.Printf("Probability of vector P is %v\n\n", P)

		P = PK
		count++
	}

	fmt.Printf("---%v tours---\n", count)
	return P
}

func (cli *CLI) transpose(v []float64, zap bool) []float64 {
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
	return result
}

func (cli *CLI) pageRankSum(P []float64) float64 {
	sum := 0.0
	for i := 0; i < len(P); i++ {
		sum += P[i]
	}
	return sum
}

func (cli *CLI) updateDelta(a, b []float64, SumOfPR float64) float64 {
	sum := 0.0
	for i := range b {
		if SumOfPR != 0 {
			b[i] = b[i] + ((1 - SumOfPR) / float64(len(b)))
		}
		sum += math.Pow(b[i]-a[i], 2)
	}
	return sum / float64(len(b))
}
