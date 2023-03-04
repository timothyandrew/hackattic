package dns

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/miekg/dns"
)

type Record struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Data string `json:"data"`
}

func (r *Record) Match(n, t string) bool {
	if strings.HasPrefix(r.Name, "*.") {
		if !strings.HasSuffix(n, r.Name[2:]) && !strings.HasSuffix(n, fmt.Sprintf("%s.", r.Name[2:])) {
			return false
		}
	} else {
		if r.Name != n && fmt.Sprintf("%s.", r.Name) != n {
			return false
		}
	}

	if r.Type != t {
		return false
	}

	return true
}

type Server struct {
	Records []Record `json:"records"`
}

func (s *Server) Start() {
	server := dns.Server{Addr: ":53", Net: "udp"}
	dns.HandleFunc(".", s.handleRequest)
	log.Println("DNS server running on :53")
	server.ListenAndServe()
}

func (s *Server) handleRequest(w dns.ResponseWriter, r *dns.Msg) {
	name := r.Question[0].Name
	queryType := dns.Type(r.Question[0].Qtype).String()

	var result Record

	for _, record := range s.Records {
		fmt.Println(record.Name)

		if record.Match(name, queryType) {
			result = record
		}
	}

	m := dns.Msg{}
	m.SetReply(r)

	if result == (Record{}) {
		log.Printf("FAIL for %s %s", name, queryType)
		w.WriteMsg(&m)
		return
	}

	log.Println("PASS", result)

	var answer dns.RR

	switch result.Type {
	case "AAAA":
		answer = &dns.AAAA{
			AAAA: net.ParseIP(result.Data),
			Hdr:  dns.RR_Header{Name: name, Rrtype: dns.TypeAAAA, Class: dns.ClassINET, Ttl: 60},
		}
	case "A":
		answer = &dns.A{
			A:   net.ParseIP(result.Data),
			Hdr: dns.RR_Header{Name: name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
		}
	case "RP":
		answer = &dns.RP{
			Txt:  ".",
			Mbox: fmt.Sprintf("%s.", result.Data),
			Hdr:  dns.RR_Header{Name: name, Rrtype: dns.TypeRP, Class: dns.ClassINET, Ttl: 60},
		}
	case "TXT":
		answer = &dns.TXT{
			Txt: []string{result.Data},
			Hdr: dns.RR_Header{Name: name, Rrtype: dns.TypeTXT, Class: dns.ClassINET, Ttl: 60},
		}
	}

	fmt.Println(answer)

	m.Authoritative = true
	m.Answer = append(m.Answer, answer)

	w.WriteMsg(&m)
}

func Run(token string) error {
	response, err := http.Get(fmt.Sprintf("https://hackattic.com/challenges/serving_dns/problem?access_token=%s", token))
	if err != nil {
		return err
	}

	server := Server{}

	jsonDecoder := json.NewDecoder(response.Body)
	err = jsonDecoder.Decode(&server)
	if err != nil {
		return err
	}

	server.Start()

	return nil
}
