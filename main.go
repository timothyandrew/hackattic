package main

import (
	"io/ioutil"

	"github.com/timothyandrew/hackattic/miniminer"
)

type Exercise interface {
	Run(token string)
}

func main() {
	tokenData, err := ioutil.ReadFile("token")
	if err != nil {
		panic(err)
	}
	token := string(tokenData)

	exercises := []Exercise{
		miniminer.Exercise{},
	}

	for _, exercise := range exercises {
		exercise.Run(token)
	}
}
