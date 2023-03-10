package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func processLine(line string) {
	line = strings.TrimSpace(line)
	if line == "" {
		return
	}

	elements := strings.Split(line, " ")
	var total int64

	for _, elem := range elements {
		var n int64
		var err error

		if strings.HasPrefix(elem, "0x") {
			n, err = strconv.ParseInt(elem[2:], 16, 64)
		} else if strings.HasPrefix(elem, "0o") {
			n, err = strconv.ParseInt(elem[2:], 8, 64)
		} else if strings.HasPrefix(elem, "0b") {
			n, err = strconv.ParseInt(elem[2:], 2, 64)
		} else if len(elem) == 1 && (elem[0] < '0' || elem[0] > '9') {
			n = int64(elem[0])
		} else {
			n, err = strconv.ParseInt(elem, 10, 64)
		}

		total += n

		if err != nil {
			panic(err)
		}
	}

	fmt.Println(total)
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			processLine(line)
			break
		} else if err != nil {
			panic(err)
		}

		processLine(line)
	}
}
