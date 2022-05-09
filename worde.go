package main

import (
	"fmt"
	"bufio"
	"os"
	"log"
	"unicode"
	"strings"
)

var list []string // this stays the same
var words []string // still possible

const N = 6

func init() {
	f, _ := os.Open("list.txt")
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		words = append(words, line)
		list = append(list, line)
	}

	if (len(words) == 0) {
		log.Fatal("No word list")
	}
}

func filter(try, hint string) {
	for i:= 0; i < len(words); i++ {
		if !ok(try, hint, words[i]) {
			words = append(words[:i], words[i+1:]...)
			i--
		}
	}
}

func ok(try, hint, word string) bool {
	for i, h := range hint[:5] {
		W := word[i]

		switch {
		case unicode.IsUpper(h):
			H := h
			if int(W) != int(H) {
				return false
			}
		case unicode.IsLower(h):
			H := unicode.ToUpper(h)

			if int(W) == int(H) {
				return false
			}

			have := false
			for _, W := range word {
				if int(W) == int(H) {
					have = true
				}
			}
			if !have {
				return false
			}
		default:
			T := try[i]
			if int(T) == int(word[i]) {
				return false
			}
			// if the letter (T) is somewhere in the hint already, don't delete words that contain
			// this letter - you are left with nothing
			skip := false
			for _, h := range hint {
				if int(unicode.ToUpper(h)) == int(T) {
					skip = true
				}
			}
			if skip {
				continue
			}
			for _, W := range word {
				if int(W) == int(T) {
					return false
				}
			}
		}
	}
	return true
}

func getRunes(s string) (runes []rune) {
	rm := make(map[rune]bool)
	for _, r := range s {
		rm[r] = true
	}
	for k, _ := range rm {
		runes = append(runes, k)
	}
	return runes
}

func best(n int) string {
	if len(words) == 0 {
		log.Fatal("Out of words")
	}

	if (n == N || len(words) < 3) {
		return words[0]
	}

	I := 0
	best := len(list)
	for i, guess := range list {
		z := worst(words, getRunes(guess))
		if z < best {
			best = z
			I = i
		}
	}

	return list[I]
}

func worst(words []string, runes []rune) int {
	if len(runes) == 0 {
		return len(words)
	}
	var left, right []string
	for _, w := range words {
		if !strings.ContainsRune(w, runes[len(runes)-1]) {
			left = append(left, w)
		} else {
			right = append(right, w)
		}
	}
	runes = runes[:len(runes)-1]
	a := worst(left, runes)
	b := worst(right, runes)
	if a > b {
		return a
	}
	return b
}

func main() {
	for i:= 0; i < N; i++ {
		var try, hint string
		fmt.Printf("%s [%d]\n", best(i), len(words))
		fmt.Scanf("%s %s", &try, &hint)
		if try == hint {
			return
		}
		filter(try, hint)
	}
}
