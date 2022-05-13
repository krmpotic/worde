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
	const maxSolves = 100
	const maxTries = 6

	shuffle := make([]string, len(list))
	copy(shuffle, list)
	rand.Shuffle(len(shuffle), func(i, j int) {
		shuffle[i], shuffle[j] = shuffle[j], shuffle[i]
	})

	for i, goal := range shuffle {
		if i == maxSolves {
			break;
		}
		if tries, ok := auto(goal, maxTries); !ok {
			t.Fatalf("Couldn't guess %q in <= 6 tries: %v\n", goal, tries)
		}
	}
	return
}

func TestFixHint(t *testing.T) {
	inputs := []struct{try string; hint string}{{"ABBAA","21..."},}
	want := []string{"21111"}

	for i, in := range inputs {
		if r := fixHint(in.try, in.hint); r != want[i] {
			t.Fatalf("fixHint(%q,%q): want %q got %q", in.try, in.hint, want[i], r)
		}
	}
}

func BenchmarkAuto(b *testing.B) {
	const maxTries = 6
	for i:= 0; i<b.N; i++ {
		auto(list[rand.Intn(len(list))], maxTries)
	}
}

func BenchmarkWorst(b *testing.B) {
	m := make(map[string][]rune, len(list))
	for _, w := range list {
		m[w] = getRunes(w)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		worst(m[list[rand.Intn(len(list))]])
	}
}

func BenchmarkBest(b *testing.B) {
	for i :=0; i < b.N; i++ {
		Best(2)
	}
}

func auto(goal string, max int) (tries []string, ok bool) {
	words = make([]string, len(list))
	copy(words, list)
	t := bestFirst
	for i := 0; i < max; i++ {
		tries = append(tries, t)
		if t == goal {
			return tries, true
		}

		Filter(t, genHint(goal,t))
		t = Best(max-i)
	}
	return tries, false
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
