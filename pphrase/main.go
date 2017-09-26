package main

import (
	"bufio"
	"crypto/rand"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {

	wordlist := flag.String("w", "words.txt", "word list to use")
	phraseLen := flag.Int("n", 5, "length of passphrase")
	bracketed := flag.Bool("bracketed", false, "return passphrase inside brackets")

	flag.Parse()

	f, err := os.Open(*wordlist)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(f)

	var words []string

	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	if len(words) != 2048 {
		log.Fatal("wordlist must contain 2048 words")
	}

	var phrase []string
	for i := 0; i < *phraseLen; i++ {
		var buf [2]byte
		_, err := rand.Read(buf[:])
		if err != nil {
			log.Fatal("error reading crypto/rand:", err)
		}
		word := int(buf[0]&0x07)<<8 + int(buf[1])
		phrase = append(phrase, words[word])
	}

	if *bracketed {
		fmt.Println(phrase)
	} else {
		for _, word := range phrase {
			fmt.Print(word + " ")
		}
		fmt.Println()
	}

}
