package miniminer

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
)

type Block struct {
	Data  [][2]interface{} `json:"data"`
	Nonce int              `json:"nonce"`
}

type Response struct {
	Difficulty int   `json:"difficulty"`
	Block      Block `json:"block"`
}

type Solution struct {
	Nonce int `json:"nonce"`
}

func calcDifficulty(block Block) (int, error) {
	ser, err := json.Marshal(block)
	if err != nil {
		return 0, err
	}

	sha := sha256.New()
	sha.Write(ser)
	sum := sha.Sum(nil)

	count := 0

	for _, b := range sum {
		for i := 7; i >= 0; i-- {
			if b&(1<<i) == 0 {
				count++
			} else {
				return count, nil
			}
		}
	}

	return 0, nil
}

func Run(token string) error {
	response, err := http.Get(fmt.Sprintf("https://hackattic.com/challenges/mini_miner/problem?access_token=%s", token))
	if err != nil {
		return err
	}

	data := Response{}

	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&data)
	if err != nil {
		return err
	}

	data.Block.Nonce = 0

	for {
		difficulty, err := calcDifficulty(data.Block)
		if err != nil {
			return err
		}

		if difficulty >= data.Difficulty {
			fmt.Println("RESULT NONCE: ", data.Block.Nonce)
			break
		}

		data.Block.Nonce++
	}

	solutionData, err := json.Marshal(Solution{Nonce: data.Block.Nonce})
	if err != nil {
		return err
	}

	response, err = http.Post(fmt.Sprintf("https://hackattic.com/challenges/mini_miner/solve?access_token=%s", token), "application/json", bytes.NewBuffer(solutionData))
	if err != nil {
		return err
	}

	fmt.Println(response.Status)
	return nil
}
