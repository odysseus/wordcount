# Word Count

Reads a file and counts unique words and their number of occurrences. Outputs a JSON file of that information to the current directory. Currently set up for English characters, but alphabet-specific tuning can be easily accomplished by modifying the `if/else` statement in `WordCount`.

### Options

- `--case-sensitive`/`-cs`: Setting either flag turns on case sensitive matching.
- `--quiet`/`-q`: Suppresses the verbose, file by file output.
- `--silent`/`-s`: Suppresses all STDOUT output of the program.
- `--name`/`-n`: Supply an alternate output file, this can be a different path, it does not append `.json` or anything else to the filename.
- `--path`/`-p`: Supplies an alternate path to store the output files in, this *will* work with multiple input files.

### Example Usage

Single file, no options
```bash
$ wordcount shakes.txt

Completed 1 of 1. Input: shakes.txt --> Output: shakes_counts.json
Lines: 124453  Words: 904625  Unique: 26675
Took: 0m 0s 569.478ms
----------------
1 Completed. 0 Failed. 1 Total. Took: 0m 0s 567.55ms
```

Alternate output file:
```bash
$ wordcount --out bard.json shakes.txt

Completed 1 of 1. Input: shakes.txt --> Output: bard.json
Lines: 124453  Words: 904625  Unique: 26675
Took: 0m 0s 583.765ms
----------------
1 Completed. 0 Failed. 1 Total. Took: 0m 0s 567.55ms
```

Multiple files:
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

Multiple files/ Alternate path:
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

Quiet
```bash
$ wordcount --quiet *.txt
3 Completed. 0 Failed. 3 Total. Took: 0m 0s 563.535ms
```

Silent
```bash
$ wordcount --silent *.txt
```
