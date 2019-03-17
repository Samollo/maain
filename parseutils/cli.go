package parseutils

import "github.com/Samollo/maain/constants"

type CLI struct {
	c []float32
	l []int
	i []int
}

func NewCLI() *CLI {
	return &CLI{
		c: make([]float32, 0),
		l: make([]int, constants.PagesToExtract),
		i: make([]int, 0)
	}
}

func (cli *CLI) AddPage(pageId int, links Links) error {
	l := removeDuplicates(links)

	coef := float32(1) / float32(len(l))

	for _, value := range l {
		cli.c = append(cli.c, coef)
		cli.i = append(cli.i, value)
	}

	nbLink := 0
	if len(cli.l) > 1 {
		nbLink = cli.l[len(cli.l)-1] + 1
	}

	cli.l = append(cli.l, nbLink)
	return nil
}

func removeDuplicates(elements Links, addToRemove ...int) Links {
	encountered := make(map[int]bool)
	for _, value := range addToRemove {
		encountered[value] = true
	}

	result := make(Links, 0)

	for v := range elements {
		if !encountered[elements[v]] {
			encountered[elements[v]] = true
			result = append(result, elements[v])
		}
	}
	return result
}
