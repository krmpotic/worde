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

var flagS = flag.Int("s", 0, "number of simulations (-1 for infinite, 0 for off)")
var flagQ = flag.Bool("q", false, "in simulation mode, just print stats")

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

func printStats(n int, info map[int]int, cd, ad time.Duration) {
	var keys []int
	for k, _ := range info {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	fmt.Printf("#%3d [ ", n)
	for _, k := range keys {
		fmt.Printf("%d:%d ", k, info[k])
	}
	fmt.Printf("] -- Time: %v Avg: %v\n", cd, ad) // TODO: fix print of info-map
}

func main() {
	flag.Parse()

	if *flagS != 0 {
		rand.Seed(time.Now().UnixNano())
		info := make(map[int]int)
		startTotal := time.Now()
		for i := 0; *flagS == -1 || i < *flagS; i++ {
			start := time.Now()
			g := simulate()
			info[g]++
			printStats(i+1, info, time.Since(start), getAvgTime(time.Since(startTotal), i+1))
		}
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

func simulate() int {
	goal := list[rand.Intn(len(list))]
	words = make([]string, len(list))
	copy(words, list)
	for i := 0; ; i++ {
		b := bestFirst
		if i != 0 {
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
