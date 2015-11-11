package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/miekg/dns"
)

var (
	printf   *bool
	debug   *bool
)

func handleHosts(w dns.ResponseWriter, r *dns.Msg) {
	firstQuestion := r.Question[0].Name
	if *debug {
		fmt.Printf("// question[0]: %v\n", firstQuestion)
	}

	var (
		rr dns.RR
		a net.IP
	)

	m := new(dns.Msg)
	m.SetReply(r)

	addrs, _ := net.LookupHost(firstQuestion)
	if *debug {
		fmt.Printf("// hosts: %v\n", addrs)
	}
	if (len(addrs) > 0) {
		match := addrs[0]
		if (len(addrs) > 1 && net.ParseIP(match).To4() == nil) {
			// TODO choose IPv4 by default (iterate over all addrs, select the first IPv4 address)
			match = addrs[1]
		}
		if *printf {
			fmt.Printf("// match: %v -> %v\n", firstQuestion, match)
		}
		a = net.ParseIP(match)
		if *debug {
			fmt.Printf("// addr: %v\n", a.To4())
		}

		rr = new(dns.A)
		rr.(*dns.A).Hdr = dns.RR_Header{
			Name: firstQuestion,
			Rrtype: dns.TypeA,
			Class: dns.ClassINET,
			Ttl: 0}
		rr.(*dns.A).A = a.To4()

		m.Answer = append(m.Answer, rr)
		m.Authoritative = true
	}

	if *debug {
		fmt.Printf("// response\n%v\n", m.String())
	}

	w.WriteMsg(m)
}

func main() {
	port := flag.Int("port", 5353, "port to run on")
	printf = flag.Bool("print", true, "print matched addresses")
	debug = flag.Bool("debug", false, "print debug logs")
	flag.Parse()

	dns.HandleFunc(".", handleHosts)

	serve := func(net string, port int) {
		server := &dns.Server{Addr: ":" + strconv.Itoa(port), Net: net}
		err := server.ListenAndServe()
		if err != nil {
			fmt.Printf("Failed to setup the " + net + " server: %s\n", err.Error())
		}
	}

	go serve("tcp", *port)
	go serve("udp", *port)

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	forever:
	for {
		select {
		case s := <-sig:
			log.Fatalf("Signal (%d) received, stopping\n", s)
			break forever
		}
	}
}
