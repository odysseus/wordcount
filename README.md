# Word Count

A program for outputting word statistics for a text file. The standalone library for use within programs is called [concordance](https://github.com/odysseus/concordance)

## Installation
First ensure that Go is installed, and your `$GOPATH` and `$GOBIN` are properly set. Go can be installed easily with homebrew using `brew install go` on mac or the package manager on Linux. Help on setting the environment variables can be found [here](https://golang.org/doc/code.html).

Once that's done you simply need to run:

```bash
$ go get github.com/odysseus/wordcount
```

From any command line.

## Options

### Modifying Output Data
- `--case-sensitive` / `-cs`: Setting either flag turns on case sensitive matching.
- `--top` / `-t`: Specifies how many words to store in the 'most used' list. By default it stores all words. Storing every word can also be signaled by a `top` value <=0

### Suppressing Terminal Output
- `--quiet` / `-q`: Suppresses the verbose, file by file output.
- `--silent` / `-s`: Suppresses all STDOUT output of the program.

### Changing Output Path and/or Filename
- `--name` / `-n`: Supply an alternate output file, this can be a different path, it does not append `.json` or anything else to the filename.
- `--path` / `-p`: Supplies an alternate path to store the output files in, this *will* work with multiple input files.

## Example Usage

**Single file, no options:**
```bash
$ wordcount shakes.txt

Completed 1 of 1. Input: shakes.txt --> Output: shakes_counts.json
Lines: 124453  Words: 904625  Unique: 26675
Took: 0m 0s 569.478ms
----------------
1 Completed. 0 Failed. 1 Total. Took: 0m 0s 567.55ms
```

**Alternate output file:**
```bash
$ wordcount --out bard.json shakes.txt

Completed 1 of 1. Input: shakes.txt --> Output: bard.json
Lines: 124453  Words: 904625  Unique: 26675
Took: 0m 0s 583.765ms
----------------
1 Completed. 0 Failed. 1 Total. Took: 0m 0s 567.55ms
```

**Multiple files:**
```bash
$ wordcount shakes.txt jabberwocky.txt hollow_men.txt

Completed 1 of 3. Input: shakes.txt --> Output: shakes_counts.json
Lines: 124453  Words: 904625  Unique: 26675
Took: 0m 0s 544.653ms
----------------
Completed 2 of 3. Input: jabberwocky.txt --> Output: jabberwocky_counts.json
Lines: 34  Words: 167  Unique: 91
Took: 0m 0s 0.518ms
----------------
Completed 3 of 3. Input: hollow_men.txt --> Output: hollow_men_counts.json
Lines: 119  Words: 427  Unique: 194
Took: 0m 0s 0.791ms
----------------
3 Completed. 0 Failed. 3 Total. Took: 0m 0s 564.517ms
```

**Multiple files/ Alternate path:**
```bash
$ wordcount --path json/ *.txt

Completed 1 of 3. Input: hollow_men.txt --> Output: json/hollow_men_counts.json
Lines: 119  Words: 427  Unique: 194
Took: 0m 0s 2.482ms
----------------
Completed 2 of 3. Input: jabberwocky.txt --> Output: json/jabberwocky_counts.json
Lines: 34  Words: 167  Unique: 91
Took: 0m 0s 1.504ms
----------------
Completed 3 of 3. Input: shakes.txt --> Output: json/shakes_counts.json
Lines: 124453  Words: 904625  Unique: 26675
Took: 0m 0s 554.698ms
----------------
3 Completed. 0 Failed. 3 Total. Took: 0m 0s 558.760ms
```

**Quiet:**
```bash
$ wordcount --quiet *.txt
3 Completed. 0 Failed. 3 Total. Took: 0m 0s 563.535ms
```

**Silent:**
```bash
$ wordcount --silent *.txt
```
