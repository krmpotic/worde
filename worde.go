package worde

import (
	_ "embed"
	"log"
	"strings"
)

//go:embed list.txt
var listTxt string
var listEmb []string  // this stays the same
var words []string // still possible
var bestFirst string

const colorOn = true
const numLetters = 5

type Solver struct {
	list []string
	left []string
	bestFirst string
	first bool
}


func init() {
	listEmb = strings.Split(listTxt, "\n")
	words = strings.Split(listTxt, "\n")

	if len(words) == 0 {
		log.Fatal("No word list")
	}

	s := NewSolver()
	bestFirst = s.best(2)

	if bestFirst == "" {
		log.Fatalf("bestFirst not initialized\n")
	}
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
		if h == '1' || h == '2' {
			m[try[i]] = true
		}
	}

	for i, h := range hint {
		if h == '.' && m[try[i]] {
			out += "1"
		} else {
			out += string(h)
		}
	}
	return
}

func ok(try, hint, word string) bool {
	for i, h := range hint[:numLetters] {
		W := word[i]
		T := try[i]

		switch {
		case h == '.':
			for i, _ := range word {
				if word[i] == T {
					return false
				}
			}
		case h == '1':
			if W == T {
				return false
			}

			have := false
			for i, _ := range word {
				if word[i] == T {
					have = true
				}
			}
			if !have {
				return false
			}
		case h == '2':
			if W != T {
				return false
			}
		}
	}
	return true
}

func getRunes(s string) (runes []rune) {
	for _, r := range s {
		f := true
		for _, R := range runes {
			if r == R {
				f = false
			}
		}
		if f {
			runes = append(runes, r)
		}
	}
	return runes
}

func (s *Solver) Best(guessesLeft int) string {
	if s.first {
		s.first = false
		return s.bestFirst
	}

	return s.best(guessesLeft)
}

func (s *Solver) best(guessesLeft int) string {
	if len(s.left) == 0 {
		return ""
	}

	if guessesLeft == 1 || len(s.left) < 3 {
		return s.left[0]
	}

	I := 0
	score := len(s.left)
	for i, guess := range s.list {
		if s := Worst(s.left, getRunes(guess)); s < score {
			score, I = s, i
		}
	}

	return s.list[I]
}

func Worst(words []string, runes []rune) (r int) {
	a := make([]int, 1<<len(runes))
	for _, w := range words {
		i_ := 0
		for i, r := range runes {
			if strings.ContainsRune(w, r) {
				i_ += 1 << i
			}
		}
		a[i_]++
		if r < a[i_] {
			r = a[i_]
		}
	}
	return
}
