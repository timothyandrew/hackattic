package ssl

import (
	"bytes"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"

	"github.com/biter777/countries"
)

type RequiredData struct {
	Domain       string `json:"domain"`
	SerialNumber string `json:"serial_number"`
	Country      string `json:"country"`
}

type Response struct {
	PrivateKey string       `json:"private_key"`
	Data       RequiredData `json:"required_data"`
}

type Solution struct {
	Certificate string `json:"certificate"`
}

func Run(token string) error {
	response, err := http.Get(fmt.Sprintf("https://hackattic.com/challenges/tales_of_ssl/problem?access_token=%s", token))
	if err != nil {
		return err
	}

	data := Response{}

	jsonDecoder := json.NewDecoder(response.Body)
	err = jsonDecoder.Decode(&data)
	if err != nil {
		return err
	}

	// privateKeyBytes, err := base64.StdEncoding.DecodeString(data.PrivateKey)
	// if err != nil {
	// 	return err
	// }

	block, _ := pem.Decode([]byte(fmt.Sprintf(`
-----BEGIN PRIVATE KEY-----
%s
-----END PRIVATE KEY-----
	`, data.PrivateKey)))
	if block == nil {
		return fmt.Errorf("no block")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return err
	}

	serial, err := strconv.ParseInt(data.Data.SerialNumber[2:], 16, 64)
	if err != nil {
		fmt.Println(data.Data.SerialNumber)
		return err
	}

	country := countries.ByName(data.Data.Country)
	if country == countries.Unknown {
		return fmt.Errorf("don't know this country: %s", data.Data.Country)
	}

	template := x509.Certificate{
		DNSNames:            []string{data.Data.Domain},
		PermittedDNSDomains: []string{data.Data.Domain},
		Subject: pkix.Name{
			Country:    []string{country.Alpha2()},
			CommonName: data.Data.Domain,
		},
		SerialNumber: big.NewInt(serial),
	}

	cert, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return err
	}

	out := base64.StdEncoding.EncodeToString(cert)
	fmt.Println(out)

	solutionData, err := json.Marshal(Solution{Certificate: out})
	if err != nil {
		return err
	}

	response, err = http.Post(fmt.Sprintf("https://hackattic.com/challenges/tales_of_ssl/solve?access_token=%s", token), "application/json", bytes.NewBuffer(solutionData))
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
