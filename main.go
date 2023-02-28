package main

import (
	"io/ioutil"

	"github.com/timothyandrew/hackattic/ssl"
)

type Exercise interface {
	Run(token string) error
}

func main() {
	tokenData, err := ioutil.ReadFile("token")
	if err != nil {
		panic(err)
	}

	token := string(tokenData)

	// err = miniminer.Run(token)
	// err = unpack.Run(token)
	err = ssl.Run(token)

	if err != nil {
		panic(err)
	}
}
