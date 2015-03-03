package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/odysseus/stopwatch"
	"os"
)

func main() {
	sw := stopwatch.New()

	file, err := os.Open("./shakes.txt")
	check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	limit := 10
	current := 0
	s := ""
	for scanner.Scan() {
		if current == limit {
			break
		}
		s += fmt.Sprintln(scanner.Text())
		current++
	}

	m := WordCountAlpha(s, false)
	j := WordMapToJSON(m, true)
	err = WriteJSONToFile(j, "./shakes_count.json")
	check(err)

	fmt.Println(sw)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func WordCount(str string) map[string]int {
	m := make(map[string]int)
	wordBreaks := []rune(" ,.;:!?\n\t(){}[]")

	var buf [2048]rune
	i := 0
	for _, c := range str {
		if hasChar(wordBreaks, c) {
			if i > 0 {
				m[string(buf[:i])]++
				i = 0
			} else {
				continue
			}
		} else {
			buf[i] = c
			i++
		}
	}
	return m
}

func hasChar(chars []rune, c rune) bool {
	for _, r := range chars {
		if c == r {
			return true
		}
	}
	return false
}

// Word count, alphabetic characters only
// str: The string of text to be word-counted
// caseSensitive: If false, all letters are downcased
func WordCountAlpha(str string, caseSensitive bool) map[string]int {
	m := make(map[string]int)

	var buf [2048]rune
	i := 0
	for _, c := range str {
		if c >= 65 && c <= 90 {
			if caseSensitive {
				buf[i] = c + 32
			} else {
				buf[i] = c
			}
			i++
		} else if c >= 97 && c <= 122 {
			buf[i] = c
			i++
		} else if i > 0 {
			m[string(buf[:i])]++
			i = 0
		} else {
			continue
		}
	}

	return m
}

func MapMergeInto(a, b map[string]int) {
	for k, v := range b {
		a[k] += v
	}
}

func WordMapToJSON(m map[string]int, humanReadable bool) []byte {
	var j []byte
	var err error
	if humanReadable {
		j, err = json.MarshalIndent(m, "", "\t")
	} else {
		j, err = json.Marshal(m)
	}
	check(err)

	return j
}

func WriteJSONToFile(js []byte, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(js)
	if err != nil {
		return err
	}

	return nil
}
