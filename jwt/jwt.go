package jwt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type MyCustomClaims struct {
	Append string `json:"append"`
	jwt.RegisteredClaims
}

type Request struct {
	AppUrl string `json:"app_url"`
}

type Response struct {
	Solution string `json:"solution"`
}

func getHandler(key string) http.HandlerFunc {
	s := strings.Builder{}

	return func(w http.ResponseWriter, req *http.Request) {
		token, err := ioutil.ReadAll(req.Body)
		if err != nil {
			panic(err)
		}

		fmt.Println(string(token))

		parsed, err := jwt.ParseWithClaims(string(token), &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(key), nil
		})
		if err != nil {
			w.WriteHeader(422)
			return
		}

		if claims, ok := parsed.Claims.(*MyCustomClaims); ok && parsed.Valid {
			if claims.Append == "" {
				solutionData, err := json.Marshal(Response{Solution: s.String()})
				if err != nil {
					panic(err)
				}

				w.Write(solutionData)
				return
			} else {
				s.WriteString(claims.Append)
				w.WriteHeader(200)
				fmt.Println("SO FAR", s.String())
				return
			}
		} else {
			panic(err)
		}
	}
}

func Run(token string) error {
	stdin, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return err
	}

	key := strings.TrimSpace(string(stdin))

	http.HandleFunc("/", getHandler(key))

	go func() {
		solutionData, err := json.Marshal(Request{AppUrl: "http://jwt.timothyandrew.net:9999"})
		if err != nil {
			panic(err)
		}

		response, err := http.Post(fmt.Sprintf("https://hackattic.com/challenges/jotting_jwts/solve?access_token=%s", token), "application/json", bytes.NewBuffer(solutionData))
		if err != nil {
			panic(err)
		}

		fmt.Println("GOT", response.Status)
	}()

	http.ListenAndServe(":9999", nil)
	return nil
}
