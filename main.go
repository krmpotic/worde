package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	N = 6
)

type limitPrint []string
func (lp limitPrint) String() (s string) {
	const max = 10
	s += "["
	for i, w := range lp {
		if i >= max {
			s += fmt.Sprintf(" +%d", len(lp)-max)
			break
		}
		if i != 0 {
			s += ", "
		}
		s += w
	}
	s += "]"
	return s
}

func main() {
	s := NewSolver()
	r := bufio.NewReader(os.Stdin)

	for i := 0; i < N; i++ {
		b := s.Best(N - i)
		if b == "" {
			fmt.Println("Out of words.")
			return
		}
		fmt.Printf("  %s %s\n", b, limitPrint(s.left))

		try, hint, err := userInput(r, b)
		if err != nil {
			fmt.Printf("userInput: %s\n", err)
			continue
		}

		s.Filter(try, hint)
	}
}

func userInput(r *bufio.Reader, best string) (try, hint string, err error) {
	in := []string{}
	fmt.Printf("> ")
	line, _ := r.ReadString('\n')
	line = strings.TrimSpace(line)
	in = strings.Split(line, " ")
	switch {
	case len(in) == 1:
		try, hint = best, in[0]
	case len(in) == 2:
		try, hint = in[0], in[1]
	default:
		return "", "", fmt.Errorf("expected [try] hint")
	}

	if len(try) != numLetters {
		return "", "", fmt.Errorf("try not %d letters", numLetters)
	}
	if len(hint) != numLetters {
		return "", "", fmt.Errorf("hint not %d letters", numLetters)
	}
	return try, hint, nil

}
