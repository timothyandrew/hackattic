package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	stdin, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(stdin), "\n")
	for _, line := range lines {
		data, err := base64.StdEncoding.DecodeString(line)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(data))
	}
}
