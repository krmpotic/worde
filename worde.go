package main

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
	byteYes  = '1'
	byteNo   = '.'

	codeNo   = 0
	codeYes  = 1
	codeHere = 2
)

type Solver struct {
	list      []string
	left      []string
	bestFirst string
	first     bool
}

type Hint struct {
	try string
	code [numLetters]int
}

func init() {
	listEmb = strings.Split(listTxt, "\n")
	bestFirst = best(listEmb, listEmb)
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

func (s *Solver) Filter(try, codeStr string) {
	h := getHint(try, codeStr)
	for i := 0; i < len(s.left); i++ {
		if !wordOk(s.left[i], h) {
			s.left = append(s.left[:i], s.left[i+1:]...)
			i--
		}
	}
}

func getHint(try, codeStr string) (Hint) {
	c := [numLetters]int{}
	m := make(map[byte]bool)
	for i, b := range codeStr {
		if b != byteNo {
			m[try[i]] = true
		}
	}

	for i, b := range codeStr {
		switch {
		case b == byteYes || (b == byteNo && m[try[i]]):
			c[i] = codeYes
		case b == byteHere:
			c[i] = codeHere
		case b == byteNo:
			c[i] = codeNo
		}
	}
	return Hint{try: try, code: c}
}

func wordOk(word string, h Hint) bool {
	for i, c := range h.code {
		W, T := word[i], h.try[i]

		switch {
		case c == codeNo && strings.ContainsRune(word, rune(T)):
			return false
		case c == codeYes && (W == T || !strings.ContainsRune(word, rune(T))):
			return false
		case c == codeHere && W != T:
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
		if s := worst(left, guess); s < score {
			score, I = s, i
		}
	}
	return list[I]
}

// number of words left in the worst case scenario
func worst(words []string, guess string) (r int) {
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
	return r
}
