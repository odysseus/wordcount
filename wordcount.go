package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	file, err := os.Open("./shakes.txt")
	checkErr(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	limit := 5
	current := 0
	s := ""
	for scanner.Scan() {
		if current == limit {
			break
		}
		s += fmt.Sprintln(scanner.Text())
		current++
	}

	var buf [1024]rune
	i := 0
	for _, c := range s {
		if c == ' ' {
			fmt.Println(string(buf[:i]))
			i = 0
		} else {
			buf[i] = c
			i++
		}
	}

}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
