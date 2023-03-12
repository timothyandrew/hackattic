package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Extra struct {
	Balance int `json:"balance"`
}

type Balance struct {
	Name    string
	Balance int   `json:"balance"`
	Extra   Extra `json:"extra"`
}

func (b *Balance) GetBalance() int {
	if b.Extra.Balance != 0 {
		return b.Extra.Balance
	}

	return b.Balance
}

func (b *Balance) GetBalanceFormatted() string {
	balance := b.GetBalance()
	balanceStr := strconv.FormatInt(int64(balance), 10)

	retRevBuilder := strings.Builder{}
	ret := strings.Builder{}
	count := 0

	for i := len(balanceStr) - 1; i >= 0; i-- {
		retRevBuilder.WriteRune(rune(balanceStr[i]))
		if count%3 == 2 && i > 0 {
			retRevBuilder.WriteRune(',')
		}
		count++
	}

	retRev := retRevBuilder.String()
	for i := len(retRev) - 1; i >= 0; i-- {
		ret.WriteRune(rune(retRev[i]))
	}

	return ret.String()
}

func read() ([]string, error) {
	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(bytes), "\n")
	return lines, nil

}

func main() {
	lines, err := read()
	if err != nil {
		panic(err)
	}

	data := []Balance{}

	for _, line := range lines {
		b := make(map[string]Balance)

		err = json.Unmarshal([]byte(line), &b)
		if err != nil {
			panic(err)
		}

		balance := Balance{}

		for k, v := range b {
			if k == "extra" {
				balance.Extra.Balance = v.Balance
			} else {
				balance.Name = k
				balance.Balance = v.Balance
			}
		}

		data = append(data, balance)
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i].GetBalance() < data[j].GetBalance()
	})

	for _, b := range data {
		fmt.Printf("%s: %d\n", b.Name, b.GetBalance())
		fmt.Printf("%s: %s\n", b.Name, b.GetBalanceFormatted())
	}
}
