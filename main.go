package main

import (
	"io/ioutil"

	"github.com/timothyandrew/hackattic/jwt"
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
	// err = ssl.Run(token)

	// err = registry.Run(token)
	// err = registry.Solve(token)

	// err = dtmf.Run(token)
	// err = dtmf.Solve(token)

	// err = rdb.Run(token)
	// err = dns.Run(token)

	// err = ws.Run(token)
	// err = face.Run(token)

	err = jwt.Run(token)

	if err != nil {
		panic(err)
	}
}
