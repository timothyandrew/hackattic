package unpack

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Response struct {
	Bytes string `json:"bytes"`
}

type Solution struct {
	Int             int32   `json:"int"`
	Uint            uint32  `json:"uint"`
	Short           int16   `json:"short"`
	Float           float64 `json:"float"`
	Double          float64 `json:"double"`
	BigEndianDouble float64 `json:"big_endian_double"`
}

func Run(token string) error {
	response, err := http.Get(fmt.Sprintf("https://hackattic.com/challenges/help_me_unpack/problem?access_token=%s", token))
	if err != nil {
		return err
	}

	data := Response{}

	jsonDecoder := json.NewDecoder(response.Body)
	err = jsonDecoder.Decode(&data)
	if err != nil {
		return err
	}

	decoder := base64.NewDecoder(base64.StdEncoding, bytes.NewBuffer([]byte(data.Bytes)))
	decoded, err := ioutil.ReadAll(decoder)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(decoded)

	var regularInt int32
	var unsignedInt uint32
	var short int16
	var padding int16

	var float float32
	var double float64
	var doubleBigEndian float64

	binary.Read(buf, binary.LittleEndian, &regularInt)
	binary.Read(buf, binary.LittleEndian, &unsignedInt)
	binary.Read(buf, binary.LittleEndian, &short)
	binary.Read(buf, binary.LittleEndian, &padding)

	binary.Read(buf, binary.LittleEndian, &float)

	binary.Read(buf, binary.LittleEndian, &double)
	binary.Read(buf, binary.BigEndian, &doubleBigEndian)

	solutionData, err := json.Marshal(Solution{
		Int:             regularInt,
		Uint:            unsignedInt,
		Short:           short,
		Float:           float64(float),
		Double:          double,
		BigEndianDouble: doubleBigEndian,
	})
	if err != nil {
		return err
	}

	response, err = http.Post(fmt.Sprintf("https://hackattic.com/challenges/help_me_unpack/solve?access_token=%s", token), "application/json", bytes.NewBuffer(solutionData))
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
