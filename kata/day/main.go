package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

func processLine(line string) {
	line = strings.TrimSpace(line)
	if line == "" {
		return
	}

	start := time.Unix(0, 0).UTC()
	offset, err := strconv.ParseInt(line, 10, 64)
	if err != nil {
		panic(err)
	}

	end := start.Add(time.Duration(offset) * 24 * time.Hour)
	fmt.Println(end.Weekday())
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
