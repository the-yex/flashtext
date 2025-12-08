package main

import (
	"math/rand"
	"os"
	"time"
)

const SIZE uint = 100000
const CORPUS_WORDS_SIZE uint = 5000
const WORDS_FILE_PATH = "./words_benchmark_test.txt"
const CORPUS_FILE_PATH = "./corpus_benchmark_test.txt"

func getRandomWord() (string, string) {
	rand.Seed(time.Now().UnixNano())
	var lenghts []int = []int{2, 4, 6, 8, 10, 12, 16, 20, 33}
	var randLen int = lenghts[rand.Intn(len(lenghts))]
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	var word []rune = make([]rune, randLen+1)
	var word_c []rune = make([]rune, randLen+1)
	for i := 0; i < randLen; i++ {
		var t rune = letters[rand.Intn(len(letters))]
		word[i] = t
		word_c[i] = t
	}
	word[randLen] = '\n'
	word_c[randLen] = ' '
	return string(word), string(word_c)
}

func gen_test_file(wordsFilePath, corpusFilePath string, corpus_words_size uint) error {
	var (
		err         error
		filewords   *os.File
		filewords_c *os.File
	)
	filewords, err = os.OpenFile(wordsFilePath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	filewords_c, err = os.OpenFile(corpusFilePath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	var (
		word   string
		word_c string
	)
	var count uint = 0
	for i := 0; i < int(SIZE); i++ {
		word, word_c = getRandomWord()
		filewords.WriteString(word)
		if count < corpus_words_size && rand.Intn(2) == 1 {
			filewords_c.WriteString(word_c)
		}
	}
	filewords.Close()
	filewords_c.Close()
	return nil
}

func main() {
	gen_test_file(WORDS_FILE_PATH, CORPUS_FILE_PATH, CORPUS_WORDS_SIZE)
}
