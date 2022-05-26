package main

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

	shuffle := make([]string, len(listEmb))
	copy(shuffle, listEmb)
	rand.Shuffle(len(shuffle), func(i, j int) {
		shuffle[i], shuffle[j] = shuffle[j], shuffle[i]
	})

	for i := 0; i < N && i < len(shuffle); i++ {
		if tries, ok := auto(shuffle[i], maxTries); !ok {
			t.Fatalf("Couldn't guess %q in <= %d tries: %v\n", shuffle[i], maxTries, tries)
		}
	}
	return
}

func auto(goal string, max int) (tries []string, ok bool) {
	s := NewSolver()
	for i := 0; i < max; i++ {
		t := s.Best(max - i)
		hint := genHint(goal, t)
		tries = append(tries, hint)
		if t == goal {
			return tries, true
		}

		s.Filter(t, hint)
	}
	return tries, false
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
