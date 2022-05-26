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

func main() {
	s := NewSolver()
	r := bufio.NewReader(os.Stdin)

	for i := 0; i < N; i++ {
		b := s.Best(N - i)
		fmt.Println(b)
	tryAgain:
		try, hint, err := userInput(r, b)
		if err != nil {
			fmt.Printf("userInput: %s\n", err)
			goto tryAgain
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
