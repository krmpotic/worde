package main

import (
	_ "embed"
	"strings"
)

//go:embed list.txt
var listTxt string
var list []string // this stays the same

const (
	numLetters = 5
)

const (
	byteNo   = '.'
	byteYes  = '1'
	byteHere = '2'

	codeNo   = 0
	codeYes  = 1
	codeHere = 2
)

type Solver struct {
	left      []string
}

type Hint struct {
	try string
	code [numLetters]int
}

func init() {
	list = strings.Split(listTxt, "\n")
}

func NewSolver() (s Solver) {
	s.left = make([]string, len(list))
	copy(s.left, list)

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
	m := make(map[byte]bool)
	for i, b := range codeStr {
		if b != byteNo {
			m[try[i]] = true
		}
	}

	c := [numLetters]int{}
	for i, b := range codeStr {
		switch {
		case b == byteYes || (b == byteNo && m[try[i]]):
			// interpret partial input correctly
			// e.g. AREAS: 1.... becomes 1..1.
			// the rest of the code relies on it
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
	if len(s.left) == 0 {
		return ""
	}

	if t == 1 || len(s.left) < 3 {
		return s.left[0]
	}

	return best(s.left)
}

// which guess is the best in the worst case scenario
func best(left []string) (word string) {
	score := len(left)
	for _, guess := range list {
		if s := worst(left, guess); s < score {
			word, score = guess, s
		}
	}
	return
}

// number of words left in the worst case scenario
func worst(words []string, guess string) (score int) {
	a := make([]int, 1<<numLetters)
	for _, w := range words {
		i_ := 0
		for i, r := range guess {
			if strings.ContainsRune(w, r) {
				i_ += 1 << i
			}
		}
		if a[i_]++; a[i_] > score {
			score = a[i_]
		}
	}
	return score
}
