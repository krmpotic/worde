package worde

import (
	_ "embed"
	"fmt"
	"log"
	"strings"
)

//go:embed list.txt
var listTxt string

var list []string  // this stays the same
var words []string // still possible
var bestFirst string

const colorOn = true

func init() {
	list = strings.Split(listTxt, "\n")
	words = strings.Split(listTxt, "\n")

	if len(words) == 0 {
		log.Fatal("No word list")
	}

	bestFirst = Best(2)
	fmt.Println("Best first word: ", bestFirst)
}

func Filter(try, hint string) {
	hint = fixHint(try, hint)
	for i := 0; i < len(words); i++ {
		if !ok(try, hint, words[i]) && len(words) > 0 {
			words = append(words[:i], words[i+1:]...)
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
	for i, h := range hint[:5] {
		W := word[i]
		T := try[i]

		switch {
		case h == '2':
			if W != T {
				return false
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
		case h == '.':
			for i, _ := range word {
				if word[i] == T {
					return false
				}
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

func Best(guessesLeft int) string {
	if len(words) == 0 {
		log.Fatal("Out of words")
	}

	if guessesLeft == 1 || len(words) < 3 {
		return words[0]
	}

	I := 0
	best := len(list)
	for i, guess := range list {
		z := worst(getRunes(guess))
		if z < best {
			best, I = z, i
		}
	}

	return list[I]
}

func worst(runes []rune) (r int) {
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

