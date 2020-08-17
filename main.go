package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"regexp"
	"github.com/miekg/dns"
)

const (
	MAX_SUBDOMAIN_LENGTH = 63
	MAX_DOMAIN_LENGTH = 253
)

type RecordParser func (string, string) (string, error)

type DNSRecord struct {
	name string
	identifier string
}

func validDomain(domain string) bool {
	if len(domain) > MAX_DOMAIN_LENGTH{
		return false
	}
	subdomains := strings.Split(domain, ".")
	for _, subdomain := range subdomains{
		if len(subdomain) > MAX_SUBDOMAIN_LENGTH {
			return false
		}
	}
	return true
}

func parseAAAARecords(m *dns.Msg, domain string) {
	subdomains := strings.Split(domain, ".")
	for ind, subdomain := range subdomains {
		if subdomain == AAAA.identifier {
			if (ind + 9 >= len(subdomains)){
				// incorrect format, would cause error
				continue
			}
			aaaaRecord := ""
			for i := 0; i < 8; i++ {
				aaaaRecord += subdomains[ind+i+1] + ":"
			}
			aaaaRecord = aaaaRecord[:len(aaaaRecord)-1] //remove last ":"
			addRecord(m, domain, AAAA.name, aaaaRecord)
		}
	}
}

func parseARecords(m *dns.Msg, domain string){
	subdomains := strings.Split(domain, ".")
	for ind, subdomain := range subdomains {
		if subdomain == A.identifier {
			if (ind + 5 >= len(subdomains)){
				// incorrect format, would cause error
				continue
			}
			aRecord := ""
			for i := 0; i < 4; i++ {
				aRecord += subdomains[ind+i+1] + "."
			}
			aRecord = aRecord[:len(aRecord)-1] //remove last "."
			addRecord(m, domain, A.name, aRecord)
		}
	}
}

func parseCNameAndAliasRecords(m *dns.Msg, domain string) {
	subdomains := strings.Split(domain, ".")
	for ind, subdomain := range subdomains {
		identifier := ""
		name := ""

		cNameMatch, _ := regexp.Match(`cname\-record\-\d+`, []byte(subdomain))
		aliasMatch, _ :=  regexp.Match(`alias\-record\-\d+`, []byte(subdomain))
		if cNameMatch {
			identifier = CNAME.identifier
			name = CNAME.name	
		} else {
			identifier = ALIAS.identifier
			name = ALIAS.name
		}

		if aliasMatch || cNameMatch {
			// get count of subdomains in cname record
			// example: cname-record-4.this.is.my.cname.example.com
			// would return this.is.my.cname
			
			subCount, err := strconv.Atoi(subdomain[len(identifier):])
			if err != nil {
				continue
			}
			
			record := "" 
			
			if (ind + subCount + 1 >= len(subdomains)){
				// incorrect format, would cause error
				continue
			}
			for i := 0; i < subCount; i++ {
				record += subdomains[ind+i+1] + "."
			}
			
			addRecord(m, domain, name, record)

		} 

	}
}

func addRecord(m *dns.Msg, name string,  recordType string, value string){
	rr, err := dns.NewRR(fmt.Sprintf("%s %s %s", name, recordType, value))
	if err == nil {
		m.Answer = append(m.Answer, rr)
	}
}

func parseQuery(m *dns.Msg) {
	for _, q := range m.Question {
		switch q.Qtype {
		case dns.TypeA:
			//log.Printf("Query for %s\n", q.Name)
			if !validDomain(q.Name) {
				return
			}
			
			parseCNameAndAliasRecords(m, q.Name)
			parseARecords(m, q.Name)
			parseAAAARecords(m, q.Name)
			

		}
	}
}

func handleDnsRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false
	switch r.Opcode {
	case dns.OpcodeQuery:
		parseQuery(m)
	}

	w.WriteMsg(m)
}

func main() {
	// attach request handler func
	dns.HandleFunc(".", handleDnsRequest)

	// start server
	port := 53
	server := &dns.Server{Addr: ":" + strconv.Itoa(port), Net: "udp"}
	log.Printf("Starting at %d\n", port)
	err := server.ListenAndServe()
	defer server.Shutdown()
	if err != nil {
		log.Fatalf("Failed to start server: %s\n ", err.Error())
	}
}