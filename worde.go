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
)

var list []string  // this stays the same
var words []string // still possible
var bestFirst string

var flagA = flag.Bool("a", false, "analyze the word list")

const N = 6
const colorOn = true

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

	bestFirst = best(2)
	fmt.Println("Best first word: ", bestFirst)
}

func filter(try, hint string) {
	for i := 0; i < len(words); i++ {
		if !ok(try, hint, words[i]) && len(words) > 0 {
			words = append(words[:i], words[i+1:]...)
			i--
		}
	}
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
				W := word[i]
				if W == T {
					have = true
				}
			}
			if !have {
				return false
			}
		case h == '.':
			if T == W {
				return false
			}

			// some implementations of the game don't repeat hints
			// if the letter (T) is somewhere in the try already,
			// don't delete words that contain this letter - you are left with nothing
			skip := false
			for i, h := range hint {
				if (h == '1' || h == '2') && try[i] == T {
					skip = true
				}
			}
			if skip {
				continue
			}

			for i, _ := range word {
				W := word[i]
				if W == T {
					return false
				}
			}
		}
	}
	return true
}

func getRunes(s string) (runes []rune) { // TODO: sort affects efficency of worst(,)
	rm := make(map[rune]bool)
	for _, r := range s {
		rm[r] = true
	}
	for k, _ := range rm {
		runes = append(runes, k)
	}
	return runes
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
	var index []int // [0,1,2,3... len(list)-1]
	for i, _ := range words {
		index = append(index, i)
	}
	for i, guess := range list {
		z := worst(index, getRunes(guess))
		if z < best {
			best = z
			I = i
		}
	}

	if best >= len(words) || best == 0 {
		return words[0]
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
			fmt.Printf("%4d/%d %v ::: ", i, len(list), goal)
			g := analyze(goal)
			fmt.Printf("[%d] %v\n", g, time.Since(start))
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

func analyze(goal string) int {
	words = make([]string, len(list))
	copy(words, list)
	b := bestFirst
	for i := 0; i < 6; i++ {
		hint := genHint(goal, b)
		filter(b, hint)
		if b == goal {
			fmt.Printf("%5s%s", hintColor(b,hint), strings.Repeat(" ", (13)*(6-i)-5)) // result & alignment
			return i+1
		}
		fmt.Printf("%5s [%3d]  ", hintColor(b,hint), len(words))
		b = best(2)
	}
	return -1
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
