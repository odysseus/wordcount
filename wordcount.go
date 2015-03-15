package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/davecheney/profile"
	"github.com/odysseus/stopwatch"
	"os"
)

var err error

func main() {
	defer profile.Start(profile.CPUProfile).Stop()

	sw := stopwatch.Start()

	in := "./shakes.txt"
	out := "./shakes_count.json"
	limit := 0
	caseSensitive := false

	lines := ReadLines(&in, limit)
	fmt.Printf("Lines read: %v lines %v\n", limit, sw)

	counts := WordCount(lines, caseSensitive)
	fmt.Printf("Words counted: %v\n", sw)

	j := WordMapToJSON(&counts, true)
	err = WriteJSONToFile(j, &out)
	check(err)
	fmt.Printf("JSON parsed and written to file: %v\n", sw)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// Reads lines until EOF or the limit is reached and returns them as a string
// path: string path of the file to be read
// limit: maximum number of lines to be read, 0 or -1 will read all lines
func ReadLines(path *string, limit int) *string {
	file, err := os.Open(*path)
	check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	current := 0
	fileBuffer := make([]rune, 8192)

	for scanner.Scan() {
		if limit > 0 && current == limit {
			break
		}
		fileBuffer = append(fileBuffer, []rune(fmt.Sprintln(scanner.Text()))...)
		current++
	}

	contents := string(fileBuffer)
	return &contents
}

// Word count, alphabetic characters only
// str: The string of text to be word-counted
// caseSensitive: If false, all letters are downcased
// return: a map of the words and counts
func WordCount(str *string, caseSensitive bool) map[string]int {
	m := make(map[string]int, 27000)

	var buf [1024]rune
	i := 0
	for _, c := range *str {
		if c == 39 {
			if i == 0 {
				continue
			} else {
				buf[i] = c
				i++
			}
		} else if c >= 65 && c <= 90 {
			if caseSensitive {
				buf[i] = c
			} else {
				buf[i] = c + 32
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

func WordMapToJSON(m *map[string]int, humanReadable bool) []byte {
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

func WriteJSONToFile(js []byte, path *string) error {
	file, err := os.Create(*path)
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
