package main

import (
	"fmt"
	"bufio"
	"os"
	"log"
	"unicode"
)

var words []string

func init() {
	f, _ := os.Open("list.txt")
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		words = append(words, line)
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

func best() string {
	if len(words) == 0 {
		log.Fatal("Out of words")
	}
	return words[0]
}

func main() {
	for i:= 0; i < 6; i++ {
		var try, hint string
		fmt.Println(best())
		fmt.Scanf("%s %s", &try, &hint)
		if try == hint {
			return
		}
		filter(try, hint)
	}
}
