package parseutils

type Pair struct {
	Id  int
	Val float64
}

type Word struct {
	value string
	freq  int
}

func NewWord(value string) *Word {
	return &Word{
		value: value,
		freq:  1,
	}
}

func (w *Word) Increment() {
	w.freq++
}

func (w *Word) Frequence() int {
	return w.freq
}

func (w *Word) String() string {
	return w.value
}
