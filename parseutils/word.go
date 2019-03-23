package parseutils

type Pair struct {
	Id  int
	Val float64
}

type Word struct {
	value string
	freq  int
}

func NewWord(value string, i ...int) *Word {
	v := 1
	if len(i) > 0 {
		v = i[0]
	}
	return &Word{
		value: value,
		freq:  v,
	}
}

func (w *Word) Increment() {
	w.freq++
}

func (w *Word) Frequence() int {
	return w.freq
}

func (w *Word) Set(f int) {
	w.freq = f
}

func (w *Word) String() string {
	return w.value
}
