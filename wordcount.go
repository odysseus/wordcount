package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/odysseus/stopwatch"
	"log"
	"os"
	"path/filepath"
)

var caseSensitive bool
var silent bool
var limit int
var outpath string

func init() {
	flag.IntVar(&limit, "limit", 0, "Maximum number of lines to read")
	flag.BoolVar(&caseSensitive, "case-sensitive", false,
		"Sets word count to case sensitive")
	flag.BoolVar(&caseSensitive, "cs", false,
		"Sets word count to case sensitive (shorthand)")
	flag.BoolVar(&silent, "s", false, "Suppresses verbose output")
	flag.BoolVar(&silent, "silent", false, "Suppresses verbose output")
	flag.StringVar(&outpath, "out", "wordcount.json", "Name of the outputted file")
}

func main() {
	flag.Parse()
	args := flag.Args()
	nargs := len(args)
	if nargs == 0 {
		log.Fatal("Error: Missing input file")
	} else if nargs > 1 && outpath != "wordcount.json" {
		log.Fatal("Error: output file can only be specified when running on a single file")
	}

	success := 0
	fail := 0

	for i, v := range args {
		sw := stopwatch.Start()

		in, err := filepath.Abs(v)
		if err != nil {
			if !silent {
				fmt.Println(failMsg(err, i+1, nargs))
			}
			fail++
			continue
		}

		out := fmt.Sprintf("%v_counts.json", FilenameSansExt(in))

		lines, linesRead, err := ReadLines(in, limit)
		if err != nil {
			if !silent {
				fmt.Println(failMsg(err, i+1, nargs))
			}
			fail++
			continue
		}

		counts, totalWords := WordCount(lines, caseSensitive)
		uniqueWords := len(counts)

		j, err := WordMapToJSON(counts, true)
		if err != nil {
			if !silent {
				fmt.Println(failMsg(err, i+1, nargs))
			}
			fail++
			continue
		}

		err = WriteJSONToFile(j, out)
		if err != nil {
			if !silent {
				fmt.Println(failMsg(err, i+1, nargs))
			}
			fail++
			continue
		}

		success++

		if !silent {
			infile, outfile := filepath.Base(in), filepath.Base(out)
			fmt.Printf("Completed %v of %v. Input: %s --> Output: %s\nLines: %v  Words: %v  Unique: %v\nTook: %v\n----------------\n",
				i+1, nargs, infile, outfile, linesRead, totalWords, uniqueWords, sw)
		}
	}
	fmt.Printf("%v Completed. %v Failed. %v Total\n", success, fail, nargs)
}

func failMsg(err error, current, total int) string {
	return fmt.Sprintf("Failed %v of %v: %s\n----------------\n",
		current, total, err.Error())
}

// Removes the filetype extension from a filename
func FilenameSansExt(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '.' {
			return path[:i]
		}
	}
	return path
}

// Reads lines until EOF or the limit is reached and returns them as a string
// path: string path of the file to be read
// limit: maximum number of lines to be read, 0 or -1 will read all lines
func ReadLines(path string, limit int) (string, int, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	current := 0
	fileBuffer := make([]rune, 4096)

	for scanner.Scan() {
		if limit > 0 && current == limit {
			break
		}
		fileBuffer = append(fileBuffer, []rune(fmt.Sprintln(scanner.Text()))...)
		current++
	}

	contents := string(fileBuffer)
	return contents, current, nil
}

// Word count, alphabetic characters only
// str: The string of text to be word-counted
// caseSensitive: If false, all letters are downcased
// return: a map of the words and counts
func WordCount(str string, caseSensitive bool) (map[string]int, int) {
	m := make(map[string]int, 4096)

	var buf [1024]rune
	i := 0
	total := 0
	for _, c := range str {
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
			total++
			i = 0
		} else {
			continue
		}
	}

	return m, total
}

// Converts a map into JSON with the option to output as human-readable with indents
func WordMapToJSON(m map[string]int, humanReadable bool) ([]byte, error) {
	var j []byte
	var err error
	if humanReadable {
		j, err = json.MarshalIndent(m, "", "\t")
	} else {
		j, err = json.Marshal(m)
	}
	if err != nil {
		return nil, err
	}

	return j, nil
}

// Takes a byte slice of JSON and the string path of the file to be written to.
// Opens, writes, and then closes the file
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
