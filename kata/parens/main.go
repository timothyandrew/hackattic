package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func processLine(line string) {
	line = strings.TrimSpace(line)
	if line == "" {
		return
	}

	count := 0
	for _, c := range line {
		if c == '(' {
			count++
		}

		if c == ')' {
			count--
		}

		if count < 0 {
			fmt.Println("no")
			return
		}
	}

	if count != 0 {
		fmt.Println("no")
		return
	}

	fmt.Println("yes")
}

func main() {
	stdin, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(stdin), "\n")
	for _, line := range lines {
		processLine(line)
	}
}
