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
var quiet bool
var silent bool
var outname string
var outpath string

func init() {
	flag.BoolVar(&caseSensitive, "case-sensitive", false, "Sets word count to case sensitive")
	flag.BoolVar(&caseSensitive, "cs", false, "Sets word count to case sensitive (shorthand)")

	flag.BoolVar(&quiet, "quiet", false, "Suppresses verbose output")
	flag.BoolVar(&quiet, "q", false, "Suppresses verbose output")

	flag.BoolVar(&silent, "silent", false, "Suppresses ALL stdout output of the program")
	flag.BoolVar(&silent, "s", false, "Suppresses ALL stdout output of the program")

	flag.StringVar(&outname, "name", "wordcount.json", "Name of the outputted file")
	flag.StringVar(&outname, "n", "wordcount.json", "Name of the outputted file")

	flag.StringVar(&outpath, "path", "./", "Alternate path for output files")
	flag.StringVar(&outpath, "p", "./", "Alternate path for output files")
}

func main() {
	tot := stopwatch.Start()
	flag.Parse()
	// The silent flag also sets the quiet flag
	if silent {
		quiet = true
	}
	// Args contains the filenames to be processed
	args := flag.Args()
	nargs := len(args)
	if nargs == 0 {
		log.Fatal("Error: Missing input file")
	} else if nargs > 1 && outname != "wordcount.json" {
		log.Fatal("Error: output file name can only be specified when running on a single file")
	}

	success := 0
	fail := 0
	if !quiet {
		fmt.Println()
	}

	// Now process each filename passed
	for i, v := range args {
		sw := stopwatch.Start()

		// Get the input filepath
		in, err := filepath.Abs(v)
		if err != nil {
			if !quiet {
				fmt.Println(failMsg(err, i+1, nargs))
			}
			fail++
			continue
		}

		// If successful create the output filename
		var out string
		if outname == "wordcount.json" {
			name := fmt.Sprintf("%v_counts.json", RawFilename(in))
			out = filepath.Join(outpath, name)
		} else {
			out = filepath.Join(outpath, outname)
		}

		// Read the file
		lines, linesRead, err := ReadLines(in)
		if err != nil {
			if !quiet {
				fmt.Println(failMsg(err, i+1, nargs))
			}
			fail++
			continue
		}

		// Count the words
		counts, totalWords := WordCount(lines, caseSensitive)
		uniqueWords := len(counts)

		// Convert to JSON
		j, err := WordMapToJSON(counts, true)
		if err != nil {
			if !quiet {
				fmt.Println(failMsg(err, i+1, nargs))
			}
			fail++
			continue
		}

		// Write that to the output file
		err = WriteJSONToFile(j, out)
		if err != nil {
			if !quiet {
				fmt.Println(failMsg(err, i+1, nargs))
			}
			fail++
			continue
		}

		// Register success
		success++

		// Output success message with summary statistics about the file
		if !quiet {
			cwd, _ := os.Getwd()
			infile, _ := filepath.Rel(cwd, in)
			fmt.Printf("Completed %v of %v. Input: %s --> Output: %s\nLines: %v  Words: %v  Unique: %v\nTook: %v\n----------------\n",
				i+1, nargs, infile, out, linesRead, totalWords, uniqueWords, sw)
		}
	}
	// On completion output the number of successes and failures, this still
	// runs in --quiet mode but not in --silent mode
	if !silent {
		fmt.Printf("%v Completed. %v Failed. %v Total. Took: %v\n",
			success, fail, nargs, tot)
	}
}

func failMsg(err error, current, total int) string {
	return fmt.Sprintf("Failed %v of %v: %s\n----------------\n",
		current, total, err.Error())
}

// Returns the raw filename with no path and no extension
func RawFilename(path string) string {
	_, file := filepath.Split(path)
	for i := len(file) - 1; i >= 0; i-- {
		if file[i] == '.' {
			return file[:i]
		}
	}
	return file
}

// Reads lines until EOF or the limit is reached and returns them as a string
// path: string path of the file to be read
// limit: maximum number of lines to be read, 0 or -1 will read all lines
func ReadLines(path string) (string, int, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	current := 0
	fileBuffer := make([]rune, 4096)

	for scanner.Scan() {
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
