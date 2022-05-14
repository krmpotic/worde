package worde

import (
	_ "embed"
	"strings"
)

//go:embed list.txt
var listTxt string
var listEmb []string // this stays the same
var bestFirst string

const (
	numLetters = 5
)

const (
	byteHere = '2'
	byteYes = '1'
	byteNo = '.'

	hintNo = 0
	hintYes = 1
	hintHere = 2
)

type hint [numLetters]int

type Solver struct {
	list      []string
	left      []string
	bestFirst string
	first     bool
}

func init() {
	listEmb = strings.Split(listTxt, "\n")
	bestFirst = best(listEmb,listEmb)
}

func NewSolver() (s Solver) {
	s.list = make([]string, len(listEmb))
	copy(s.list, listEmb)
	s.left = make([]string, len(listEmb))
	copy(s.left, listEmb)

	s.bestFirst = bestFirst
	s.first = true

	return
}

func (s *Solver) Filter(try, byte string) {
	h := getHint(try, byte)
	for i := 0; i < len(s.left); i++ {
		if !ok(try, s.left[i], h) && len(s.left) > 0 {
			s.left = append(s.left[:i], s.left[i+1:]...)
			i--
		}
	}
}

func getHint(try, hintStr string) (h hint) {
	m := make(map[byte]bool)
	for i, b := range hintStr {
		if b != byteNo {
			m[try[i]] = true
		}
	}

	for i, b := range hintStr {
		switch {
		case b == byteYes || (b == byteNo && m[try[i]]):
			h[i] = hintYes
		case b == byteHere:
			h[i] = hintHere
		case b == byteNo:
			h[i] = hintNo
		}
	}
	return
}

func ok(try, word string, hnt hint) bool {
	for i, h := range hnt {
		W, T := word[i], try[i]

		switch {
		case h == hintNo && strings.ContainsRune(word, rune(T)):
			return false
		case h == hintYes && (W == T || !strings.ContainsRune(word, rune(T))):
			return false
		case h == hintHere && W != T:
			return false
		}
	}
	return true
}

func (s *Solver) Best(t int) string {
	if s.first {
		s.first = false
		return s.bestFirst
	}

	if len(s.left) == 0 {
		return "" // TODO: fail better
	}

	if t == 1 || len(s.left) < 3 {
		return s.left[0]
	}

	return best(s.left, s.list)
}

// which guess is the best in the worst case scenario
func best(left, list []string) string {
	I, score := 0, len(left)
	for i, guess := range list {
		if s := Worst(left, guess); s < score {
			score, I = s, i
		}
	}
	return list[I]
}

// number of words left in the worst case scenario
func Worst(words []string, guess string) (r int) {
	a := make([]int, 1<<numLetters)
	for _, w := range words {
		i_ := 0
		for i, r := range guess {
			if strings.ContainsRune(w, r) {
				i_ += 1 << i
			}
		}
		if a[i_]++; r < a[i_] {
			r = a[i_]
		}
	}
	return
}
