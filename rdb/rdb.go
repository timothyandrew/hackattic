package rdb

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type Requirements struct {
	CheckTypeOf string `json:"check_type_of"`
}

type Stats struct {
	DbCount     int
	EmojiValue  string
	ExpiryMs    int
	TypeToCheck string
}

type Response struct {
	Data         string       `json:"rdb"`
	Requirements Requirements `json:"requirements"`
}

func DumpHex(bytes []byte) {
	for _, b := range bytes {
		fmt.Printf("%x ", b)
	}

	fmt.Println("")
}

func PrintKV(key, vs string, vi int) {
	if vs == "" {
		fmt.Printf("%s -> %d\n", key, vi)
	} else {
		fmt.Printf("%s -> %s\n", key, vs)
	}

}

func decodeInt(b *bufio.Reader) (int, error) {
	length, _ := b.ReadByte()

	switch (length & 0xC0) >> 6 {
	case 0:
		var value uint8
		binary.Read(bytes.NewBuffer([]byte{length}), binary.LittleEndian, &value)
		return int(value), nil

	case 1:
		panic("unimplemented 1")
	case 2:
		panic("unimplemented 2")
	case 3:
		panic("unimplemented 3")
	}

	return 0, fmt.Errorf("FAILED")
}

func decodeString(b *bufio.Reader) (string, int, error) {
	length, _ := b.ReadByte()

	switch (length & 0xC0) >> 6 {
	case 0:
		value := make([]byte, length)
		b.Read(value)
		return string(value), 0, nil

	case 1:
		next, _ := b.ReadByte()

		var fullLength uint16
		binary.Read(bytes.NewBuffer([]byte{length & 0x3F, next}), binary.BigEndian, &fullLength)

		value := make([]byte, fullLength)
		b.Read(value)
		return string(value), 0, nil

	case 2:
		panic("unimplemented")
	case 3:
		switch length & 0x3F {
		case 0:
			var data uint8
			binary.Read(b, binary.LittleEndian, &data)
			return "", int(data), nil
		case 1:
			var data uint16
			binary.Read(b, binary.LittleEndian, &data)
			return "", int(data), nil
		case 2:
			var data uint32
			binary.Read(b, binary.LittleEndian, &data)
			return "", int(data), nil
		}
	}

	return "", 0, fmt.Errorf("FAILED")

}

func Load(reader io.Reader, checkTypeOf string) (ret Stats) {
	b := bufio.NewReader(reader)

	header := make([]byte, 5)
	b.Read(header)

	version := make([]byte, 4)
	b.Read(version)

	for {
		next, _ := b.ReadByte()
		switch next {
		case 0xFA:
			fmt.Print("AUX: ")

			key, _, _ := decodeString(b)
			valueStr, valueInt, _ := decodeString(b)

			PrintKV(key, valueStr, valueInt)

		case 0xFE:
			fmt.Print("\nSelected database: ")
			dbNumber, _ := decodeInt(b)
			ret.DbCount++
			fmt.Println(dbNumber)

		case 0xFD:
			var ts uint32
			binary.Read(b, binary.LittleEndian, &ts)
			fmt.Printf("EXPIRE (sec): %d\n", ts)

		case 0xFC:
			var ts uint64
			binary.Read(b, binary.LittleEndian, &ts)
			ret.ExpiryMs = int(ts)
			fmt.Printf("EXPIRE (ms): %d\n", ts)

		case 0xFB:
			fmt.Println("\nRESIZEDB: ")

			hashTableSize, _ := decodeInt(b)
			expiryHashTableSize, _ := decodeInt(b)

			fmt.Printf("  hash table size -> %d\n", hashTableSize)
			fmt.Printf("  expiry hash table size -> %d\n\n", expiryHashTableSize)

		case 0x00:
			key, _, _ := decodeString(b)
			valueStr, valueInt, _ := decodeString(b)

			if key == checkTypeOf {
				ret.TypeToCheck = "string"
			}

			if len([]byte(key)) == 4 {
				ret.EmojiValue = valueStr
			}

			PrintKV(key, valueStr, valueInt)

		case 0x0A:
			key, _, _ := decodeString(b)
			decodeString(b)

			if key == checkTypeOf {
				ret.TypeToCheck = "list"
			}

			fmt.Printf("%s -> (ziplist)\n", key)

		case 0x0D:
			key, _, _ := decodeString(b)
			decodeString(b)

			if key == checkTypeOf {
				ret.TypeToCheck = "hash"
			}

			fmt.Printf("%s -> (hashmap)\n", key)

		case 0x0C:
			key, _, _ := decodeString(b)
			decodeString(b)

			if key == checkTypeOf {
				ret.TypeToCheck = "sortedset"
			}

			fmt.Printf("%s -> (sorted set)\n", key)

		case 0x0B:
			key, _, _ := decodeString(b)
			decodeString(b)

			if key == checkTypeOf {
				ret.TypeToCheck = "set"
			}

			fmt.Printf("%s -> (intset)\n", key)

		case 0xFF:
			fmt.Printf("DONE with %d databases\n", ret.DbCount)
			return

		default:
			fmt.Printf("Don't know: %x\n", next)
			return
		}
	}

}

func Run(token string) error {
	response, err := http.Get(fmt.Sprintf("https://hackattic.com/challenges/the_redis_one/problem?access_token=%s", token))
	if err != nil {
		return err
	}

	data := Response{}

	jsonDecoder := json.NewDecoder(response.Body)
	err = jsonDecoder.Decode(&data)
	if err != nil {
		return err
	}

	dbRaw, err := base64.StdEncoding.DecodeString(data.Data)
	if err != nil {
		return err
	}

	stats := Load(bytes.NewBuffer(dbRaw), data.Requirements.CheckTypeOf)

	solutionData, err := json.Marshal(map[string]interface{}{
		"db_count":                    stats.DbCount,
		"emoji_key_value":             stats.EmojiValue,
		"expiry_millis":               stats.ExpiryMs,
		data.Requirements.CheckTypeOf: stats.TypeToCheck,
	})
	if err != nil {
		return err
	}

	response, err = http.Post(fmt.Sprintf("https://hackattic.com/challenges/the_redis_one/solve?access_token=%s", token), "application/json", bytes.NewBuffer(solutionData))
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
