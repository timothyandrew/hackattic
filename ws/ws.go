package ws

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/net/websocket"
)

type Response struct {
	Token string `json:"token"`
}

func computeInterval(d time.Duration) int {
	intervals := []int64{700, 1500, 2000, 2500, 3000}
	duration := d.Milliseconds()

	min := math.MaxInt64
	minIndex := -1

	for i, interval := range intervals {
		diff := math.Abs(float64(duration - interval))
		if int(diff) < min {
			min = int(diff)
			minIndex = i
		}
	}

	return int(intervals[minIndex])
}

func Run(token string) error {
	response, err := http.Get(fmt.Sprintf("https://hackattic.com/challenges/websocket_chit_chat/problem?access_token=%s", token))
	if err != nil {
		return err
	}

	data := Response{}

	jsonDecoder := json.NewDecoder(response.Body)
	err = jsonDecoder.Decode(&data)
	if err != nil {
		return err
	}

	origin := "https://hackattic.com/"
	url := fmt.Sprintf("wss://hackattic.com/_/ws/%s", data.Token)
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		return err
	}

	start := time.Now()

	for {
		var msg = make([]byte, 512)
		n, err := ws.Read(msg)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		log.Println("recv:", string(msg[:n]))

		if string(msg[:n]) == "ping!" {
			elapsed := time.Since(start)
			interval := computeInterval(elapsed)
			intervalStr := strconv.FormatInt(int64(interval), 10)
			start = time.Now()

			log.Println("send:", intervalStr)

			_, err := ws.Write([]byte(intervalStr))
			if err != nil {
				return err
			}
		}

	}

	return nil
}
