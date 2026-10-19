package main

import (
	"fmt"
	"github.com/miekg/dns"
	"net"
	"strings"
)

// Temporary registry for TonOS .dapp domains (will be replaced by TON smart-contract queries)
var dappRegistry = map[string]string{
	"telegram.dapp.": "127.0.0.1",
	"chat.dapp.":     "127.0.0.1",
	"wallet.dapp.":   "127.0.0.1",
}

func handleDnsRequest(w dns.ResponseWriter, r *dns.Msg) {
	msg := new(dns.Msg)
	msg.SetReply(r)
	msg.Authoritative = true

	for _, q := range r.Question {
		domain := strings.ToLower(q.Name)
		fmt.Printf("[TonOS DNS] Query received: %s\n", domain)

		if ip, exists := dappRegistry[domain]; exists {
			fmt.Printf("-> Found in TonOS registry! Routing to %s\n", ip)
			if q.Qtype == dns.TypeA {
				rr := &dns.A{
					Hdr: dns.RR_Header{Name: q.Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
					A:   net.ParseIP(ip),
				}
				msg.Answer = append(msg.Answer, rr)
			}
		} else if strings.HasSuffix(domain, ".dapp.") {
			fmt.Println("-> .dapp domain is not registered.")
			msg.SetRcode(r, dns.RcodeNameError)
		} else {
			// Forward regular traffic to Cloudflare public DNS
			client := new(dns.Client)
			resp, _, err := client.Exchange(r, "1.1.1.1:53")
			if err == nil && resp != nil {
				w.WriteMsg(resp)
				return
			}
		}
	}

	w.WriteMsg(msg)
}

func main() {
	port := ":53"
	server := &dns.Server{Addr: port, Net: "udp"}
	dns.HandleFunc(".", handleDnsRequest)

	fmt.Printf("TonOS DNS daemon started on port %s...\n", port)
	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("DNS server error: %s\n", err.Error())
	}
}
