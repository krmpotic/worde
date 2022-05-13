package worde

import (
	_ "embed"
	"log"
	"strings"
)

//go:embed list.txt
var listTxt string
var listEmb []string // this stays the same
var bestFirst string

const (
	colorOn    = true
	numLetters = 5

	hint2 = '2'
	hint1 = '1'
	hint0 = '.'
)

type Solver struct {
	list      []string
	left      []string
	bestFirst string
	first     bool
}

func init() {
	listEmb = strings.Split(listTxt, "\n")

	s := NewSolver()
	bestFirst = s.Best(2)

	if bestFirst == "" {
		log.Fatalf("bestFirst not initialized\n")
	}
}

func NewSolver() (s Solver) {
	s.list = make([]string, len(listEmb))
	copy(s.list, listEmb)
	s.left = make([]string, len(listEmb))
	copy(s.left, listEmb)

	if bestFirst != "" {
		s.bestFirst = bestFirst
		s.first = true
	}

	return
}

func (s *Solver) Filter(try, hint string) {
	hint = fixHint(try, hint)
	for i := 0; i < len(s.left); i++ {
		if !ok(try, hint, s.left[i]) && len(s.left) > 0 {
			s.left = append(s.left[:i], s.left[i+1:]...)
			i--
		}
	}
}

func fixHint(try, hint string) (out string) {
	m := make(map[byte]bool)
	for i, h := range hint {
		if h == hint1 || h == hint2 {
			m[try[i]] = true
		}
	}

	for i, h := range hint {
		if h == hint0 && m[try[i]] {
			out += string(hint1)
		} else {
			out += string(h)
		}
	}
	return
}

func ok(try, hint, word string) bool {
	for i, h := range hint {
		W, T := word[i], try[i]

		switch {
		case h == hint0 && strings.ContainsRune(word, rune(T)):
			return false
		case h == hint1 && (W == T || !strings.ContainsRune(word, rune(T))):
			return false
		case h == hint2 && W != T:
			return false
		}
	}
	return true
}

func (s *Solver) Best(guessesLeft int) string {
	if s.first {
		s.first = false
		return s.bestFirst
	}

	if len(s.left) == 0 {
		return "" // TODO: fail better
	}

	if guessesLeft == 1 || len(s.left) < 3 {
		return s.left[0]
	}

	I, score := 0, len(s.left)
	for i, guess := range s.list {
		if s := Worst(s.left, guess); s < score {
			score, I = s, i
		}
	}

	return s.list[I]
}

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
