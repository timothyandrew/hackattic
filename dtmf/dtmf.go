package dtmf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/Hallicopter/go-dtmf/dtmf"

	"github.com/go-audio/wav"
)

type Solution struct {
	Sequence string `json:"sequence"`
}

func Solve(token string) error {
	s := "<solution>"

	solutionData, err := json.Marshal(Solution{Sequence: s})
	if err != nil {
		return err
	}

	response, err := http.Post(fmt.Sprintf("https://hackattic.com/challenges/touch_tone_dialing/solve?access_token=%s", token), "application/json", bytes.NewBuffer(solutionData))
	if err != nil {
		return err
	}

	fmt.Println(response.Status)

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body))
	return nil
}

func Run(token string) error {
	f, err := os.Open("/Users/tim/Downloads/touch_tone.wav")
	if err != nil {
		return err
	}
	defer f.Close()

	decoder := wav.NewDecoder(f)

	for {
		chunk, err := decoder.NextChunk()
		if err != nil {
			break
		}

		chunkData, _ := ioutil.ReadAll(chunk)

		result, err := dtmf.DecodeDTMFFromBytes(chunkData, float64(decoder.SampleRate), 5)
		if err != nil {
			return err
		}

		fmt.Println(result)
	}

	return nil
}
