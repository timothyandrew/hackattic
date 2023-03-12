package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func main() {
	stdin, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	input := strings.Split(string(stdin), " ")

	start, err := strconv.ParseInt(strings.TrimSpace(input[0]), 10, 64)
	if err != nil {
		panic(err)
	}

	end, err := strconv.ParseInt(strings.TrimSpace(input[1]), 10, 64)
	if err != nil {
		panic(err)
	}

	for i := start; i <= end; i++ {
		if i%3 == 0 && i%5 == 0 {
			fmt.Println("FizzBuzz")
		} else if i%3 == 0 {
			fmt.Println("Fizz")
		} else if i%5 == 0 {
			fmt.Println("Buzz")
		} else {
			fmt.Println(i)
		}
	}
}
