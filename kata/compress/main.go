package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func processLine(line string) {
	line = strings.TrimSpace(line)
	if line == "" {
		return
	}

	count := 0
	result := strings.Builder{}

	for i := 0; i < len(line); i++ {
		for i < len(line)-1 && line[i+1] == line[i] {
			i++
			count++
		}

		if count == 1 {
			result.WriteString(string(line[i]))
		}

		if count > 1 {
			result.WriteString(strconv.FormatInt(int64(count+1), 10))
		}

		result.WriteString(string(line[i]))
		count = 0
	}

	fmt.Println(result.String())
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
