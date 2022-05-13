package worde

import (
	"math/rand"
	"strings"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func TestAutoSolve(t *testing.T) {
	const maxTries = 6
	const N = 50

	for i := 0; i < N; i++ {
		if goal, tries, ok := auto(maxTries); !ok {
			t.Fatalf("Couldn't guess %q in <= %d tries: %v\n", goal, maxTries, tries)
		}
	}
	return
}

func auto(max int) (goal string, tries []string, ok bool) {
	s := NewSolver()
	goal = s.list[rand.Intn(len(s.list))]
	for i := 0; i < max; i++ {
		t := s.Best(max - i)
		hint := genHint(goal, t)
		tries = append(tries, hint)
		if t == goal {
			return goal, tries, true
		}

		s.Filter(t, hint)
	}
	return goal, tries, false
}

func genHint(goal, try string) (hint string) {
	for i, r := range try {
		switch {
		case r == rune(goal[i]):
			hint += string(byteHere)
		case strings.ContainsRune(goal, r):
			hint += string(byteYes)
		default:
			hint += string(byteNo)
		}
	}
	return hint
}
