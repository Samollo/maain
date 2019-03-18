package parseutils

type CLI struct {
	c []float32
	l []int
	i []int
}

func (cli *CLI) GetC() []float32 {
	return cli.c
}

func (cli *CLI) GetL() []int {
	return cli.l
}
func (cli *CLI) GetI() []int {
	return cli.i
}

func NewCLI() *CLI {
	return &CLI{
		c: make([]float32, 0),
		l: make([]int, 0),
		i: make([]int, 0),
	}
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
