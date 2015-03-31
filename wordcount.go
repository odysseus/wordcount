package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/odysseus/concordance"
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
var topwords int

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

	flag.IntVar(&topwords, "top", 0, "Maximium number of top words to store")
	flag.IntVar(&topwords, "t", 0, "Maximium number of top words to store")
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
		msg, err := processFile(v, i+1, nargs)
		if err != nil {
			if !quiet {
				fmt.Print(msg)
			}
			fail++
			continue
		}
		success++

		fmt.Print(msg)
	}

	// On completion output the number of successes and failures, this still
	// runs in --quiet mode but not in --silent mode
	if !silent {
		fmt.Printf("%v Completed. %v Failed. %v Total. Took: %v\n",
			success, fail, nargs, tot)
	}
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

// Converts a map into JSON with the option to output as human-readable with indents
func ConcordanceToJSON(c *concordance.Concordance, humanReadable bool) ([]byte, error) {
	var j []byte
	var err error
	if humanReadable {
		j, err = json.MarshalIndent(c, "", "\t")
	} else {
		j, err = json.Marshal(c)
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

func failMsg(err error, current, total int) string {
	return fmt.Sprintf("Failed %v of %v: %s\n----------------\n",
		current, total, err.Error())
}

func processFile(path string, n, total int) (string, error) {
	sw := stopwatch.Start()

	// Get the input filepath
	in, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	// If successful create the output filename
	var out string
	if outname == "wordcount.json" {
		name := fmt.Sprintf("%v_counts.json", RawFilename(in))
		out = filepath.Join(outpath, name)
	} else {
		out = filepath.Join(outpath, outname)
	}

	// Read the file and wrap it in a scanner
	infile, err := os.Open(in)
	if err != nil {
		return failMsg(err, n, total), err
	}
	scanner := bufio.NewScanner(infile)

	// Count the words
	counts := concordance.NewConcordance(scanner, caseSensitive, topwords)

	// Convert to JSON
	j, err := ConcordanceToJSON(counts, true)
	if err != nil {
		return failMsg(err, n, total), err
	}

	// Write that to the output file
	err = WriteJSONToFile(j, out)
	if err != nil {
		return failMsg(err, n, total), err
	}

	// Output success message with summary statistics about the file
	cwd, _ := os.Getwd()
	infilename, _ := filepath.Rel(cwd, in)
	msg := fmt.Sprintf("Completed %v of %v. Input: %s --> Output: %s\nWords: %v  Unique: %v\nTook: %v\n----------------\n",
		n, total, infilename, out, counts.Total, counts.Unique, sw)

	return msg, nil
}
