package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

func Convert(s string) string {
	ret := strings.Builder{}
	prefix := regexp.MustCompile(`^(p|d|u\d\d|i\d\d)`)

	s = prefix.ReplaceAllString(s, "")

	for i, c := range s {
		if c < 'a' {
			if i > 0 {
				ret.WriteRune('_')
			}
			ret.WriteString(strings.ToLower(string(c)))
		} else {
			ret.WriteRune(c)
		}
	}

	return ret.String()
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}

		s := strings.TrimSpace(line)
		fmt.Println(Convert(s))
	}
}
