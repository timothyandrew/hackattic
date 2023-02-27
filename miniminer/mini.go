package miniminer

import "fmt"

type Exercise struct {
}

func (e Exercise) Run(token string) {
	fmt.Println(token)
}
