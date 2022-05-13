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
	const N = 10

	for i := 0; i < N; i++ {
		if goal, tries, ok := auto(maxTries); !ok {
			t.Fatalf("Couldn't guess %q in <= %d tries: %v\n", goal, maxTries, tries)
		}
	}
	return
}

func TestFixHint(t *testing.T) {
	inputs := []struct {
		try  string
		hint string
	}{{"ABBAA", "21..."}}
	want := []string{"21111"}

	for i, in := range inputs {
		if r := fixHint(in.try, in.hint); r != want[i] {
			t.Fatalf("fixHint(%q,%q): want %q got %q", in.try, in.hint, want[i], r)
		}
	}
}

func auto(max int) (goal string, tries []string, ok bool) {
	s := NewSolver()
	goal = s.list[rand.Intn(len(s.list))]
	for i := 0; i < max; i++ {
		t := s.Best(max - i)
		hint := genHint(goal, t)
		tries = append(tries, hintColor(t, hint))
		if t == goal {
			return goal, tries, true
		}

		s.Filter(t, hint)
	}
	return goal, tries, false
}

func hintColor(try, hint string) (str string) {
	const (
		Black   = "\033[1;30m"
		Red     = "\033[1;31m"
		Green   = "\033[1;32m"
		Yellow  = "\033[1;33m"
		Purple  = "\033[1;34m"
		Magenta = "\033[1;35m"
		Teal    = "\033[1;36m"
		White   = "\033[1;37m"
		Reset   = "\033[0m"
	)

	if !colorOn {
		return try
	}

	for i, h := range hint {
		T := string(try[i])
		switch {
		case h == '1':
			str += Yellow + T + Reset
		case h == '2':
			str += Green + T + Reset
		default:
			str += T
		}
	}
	return str
}

func genHint(goal, try string) (hint string) {
	for i, r := range try {
		switch {
		case r == rune(goal[i]):
			hint += "2"
		case strings.ContainsRune(goal, r):
			hint += "1"
		default:
			hint += "."
		}
	}
	return hint
}
