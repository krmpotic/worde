package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"
	"unicode"
)

var list []string  // this stays the same
var words []string // still possible
var bestFirst string

var flagA = flag.Bool("a", false, "analyze the word list")
var flagQ = flag.Bool("q", false, "in analyze mode, just print stats")

const N = 6

func init() {
	f, _ := os.Open("list.txt")
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		words = append(words, line)
		list = append(list, line)
	}

	if len(words) == 0 {
		log.Fatal("No word list")
	}

	bestFirst = best(len(words))
	fmt.Println("Best first word: ", bestFirst)
}

func filter(try, hint string) {
	for i := 0; i < len(words); i++ {
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

func initIndexList() (index []int) {
	for i, _ := range words {
		index = append(index, i)
	}
	return index
}

func best(guessesLeft int) string {
	if len(words) == 0 {
		log.Fatal("Out of words")
	}

	if guessesLeft == 1 || len(words) < 3 {
		return words[0]
	}

	I := 0
	best := len(list)
	for i, guess := range list {
		z := worst(initIndexList(), getRunes(guess))
		if z < best {
			best = z
			I = i
		}
	}

	return list[I]
}

func worst(indexes []int, runes []rune) int {
	if len(runes) == 0 {
		return len(indexes)
	}
	var left, right []int
	for _, i := range indexes {
		if !strings.ContainsRune(words[i], runes[len(runes)-1]) {
			left = append(left, i)
		} else {
			right = append(right, i)
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

func getAvgTime(d time.Duration, i int) time.Duration {
	return time.Duration(int64(d) / int64(i))
}

func printStats(info map[int]int, time time.Duration) {
	var keys []int
	for k, _ := range info {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	fmt.Printf("[ ")
	for _, k := range keys {
		fmt.Printf("%d:%d ", k, info[k])
	}
	fmt.Printf("] -- Avg Time: %v\n", time) // TODO: fix print of info-map
}

func main() {
	flag.Parse()

	if *flagA {
		rand.Seed(time.Now().UnixNano())
		info := make(map[int]int)
		start := time.Now()
		for i, goal := range list {
			start := time.Now()
			g := analyze(goal)
			fmt.Printf("%d/%d %q [%d] %v\n", i, len(list), goal, g, time.Since(start))
			info[g]++
		}
		printStats(info, getAvgTime(time.Since(start), len(list)))
		return
	}

	for i := 0; i < N; i++ {
		var try, hint string
		fmt.Printf("%s [%d/%d]\n", best(N-i), len(words), len(list))
		fmt.Scanf("%s %s", &try, &hint) // TODO: add option for hint only, which means try == best
		if try == hint {
			return
		}
		filter(try, hint)
	}
}

func genHint(goal, try string) (hint string) {
	for i, r := range try {
		if r == rune(goal[i]) {
			hint += string(r)
			continue
		}
		if strings.ContainsRune(goal, r) {
			hint += string(unicode.ToLower(r))
		} else {
			hint += "."
		}
	}
	return hint
}

func analyze(goal string) int {
	words = make([]string, len(list))
	copy(words, list)
	for i := 1; ; i++ {
		b := bestFirst
		if i != 1 {
			b = best(len(words))
		}
		hint := genHint(goal, b)
		filter(b, hint)
		if !*flagQ {
			fmt.Printf("%s %s [%d/%d]\n", b, hint, len(words), len(list))
		}
		if b == goal {
			return i
		}
	}
}
