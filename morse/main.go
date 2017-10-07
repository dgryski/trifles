package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
)

var nato = map[rune]string{
	'a': "alpha",
	'b': "bravo",
	'c': "charlie",
	'd': "delta",
	'e': "echo",
	'f': "foxtrot",
	'g': "golf",
	'h': "hotel",
	'i': "india",
	'j': "juliett",
	'k': "kilo",
	'l': "lima",
	'm': "mike",
	'n': "november",
	'o': "oscar",
	'p': "papa",
	'q': "quebec",
	'r': "romeo",
	's': "sierra",
	't': "tango",
	'u': "uniform",
	'v': "victor",
	'w': "whiskey",
	'x': "xray",
	'y': "yankee",
	'z': "zulu",
}

var morse = map[rune]string{
	'a': ".-",
	'b': "-...",
	'c': "-.-.",
	'd': "-..",
	'e': ".",
	'f': "..-.",
	'g': "--.",
	'h': "....",
	'i': "..",
	'j': ".---",
	'k': "-.-",
	'l': ".-..",
	'm': "--",
	'n': "-.",
	'o': "---",
	'p': ".--.",
	'q': "--.-",
	'r': ".-.",
	's': "...",
	't': "-",
	'u': "..-",
	'v': "...-",
	'w': ".--",
	'x': "-..-",
	'y': "-.--",
	'z': "--..",
}

var regexps = map[rune]string{
	'.': "[a-m]",
	'-': "[n-z]",
}

type target struct {
	r  rune
	re string
}

func slurp(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func main() {
	var targets []target
	for r := rune('a'); r <= 'z'; r++ {
		var re string
		for _, c := range morse[r] {
			re += regexps[c]
		}
		targets = append(targets, target{r: r, re: re})
	}

	var words []string

	words, err := slurp("/usr/share/dict/words")
	if err != nil {
		log.Fatalln(err)
	}

	for _, t := range targets {
		re := regexp.MustCompile("^" + t.re + "$")
		for _, v := range words {
			if re.MatchString(v) {
				fmt.Println(nato[t.r], v)
			}
		}

	}
}
