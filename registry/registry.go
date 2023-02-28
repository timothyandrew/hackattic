package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Credentials struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type Response struct {
	Credentials  Credentials `json:"credentials"`
	IgnitionKey  string      `json:"ignition_key"`
	TriggerToken string      `json:"trigger_token"`
}

type TriggerRequest struct {
	Host string `json:"registry_host"`
}

type SolutionRequest struct {
	Secret string `json:"secret"`
}

func Solve(token string) error {
	solution := ""

	solutionData, err := json.Marshal(SolutionRequest{Secret: solution})
	if err != nil {
		return err
	}

	response, err := http.Post(fmt.Sprintf("https://hackattic.com/challenges/dockerized_solutions/solve?access_token=%s", token), "application/json", bytes.NewBuffer(solutionData))
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
	urlData, err := ioutil.ReadFile("registry_url")
	if err != nil {
		panic(err)
	}

	url := string(urlData)
	fmt.Println(url)

	response, err := http.Get(fmt.Sprintf("https://hackattic.com/challenges/dockerized_solutions/problem?access_token=%s", token))
	if err != nil {
		return err
	}

	data := Response{}

	jsonDecoder := json.NewDecoder(response.Body)
	err = jsonDecoder.Decode(&data)
	if err != nil {
		return err
	}

	fmt.Println(data)

	triggerData, err := json.Marshal(TriggerRequest{Host: "registry.timothyandrew.net"})
	if err != nil {
		return err
	}

	response, err = http.Post(fmt.Sprintf("https://hackattic.com/_/push/%s", data.TriggerToken), "application/json", bytes.NewBuffer(triggerData))
	if err != nil {
		return err
	}

	fmt.Println(response.Status)

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body))
	fmt.Println("IGNITION", data.IgnitionKey)

	return nil
}
