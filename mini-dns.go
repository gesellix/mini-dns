package main

import (
	"flag"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/golang/glog"
	"github.com/miekg/dns"
	"log"
)

func newRR(host string, addr net.IP) (dns.RR) {
	rr := new(dns.A)
	rr.Hdr = dns.RR_Header{
		Name: host,
		Rrtype: dns.TypeA,
		Class: dns.ClassINET,
		Ttl: 0}
	rr.A = addr.To4()
	return rr
}

func findFirstIPv4(addrs []string) (match string) {
	match = ""
	if len(addrs) == 0 {
		return
	}
	for _, match = range addrs {
		if net.ParseIP(match).To4() != nil {
			return
		}
	}
	return
}

func handleFirstHost(w dns.ResponseWriter, r *dns.Msg) {
	firstQuestion := r.Question[0].Name
	glog.Infof("first question: %v\n", firstQuestion)

	m := new(dns.Msg)
	m.SetReply(r)

	addrs, _ := net.LookupHost(firstQuestion)
	glog.Infof("hosts: %v\n", addrs)
	log.Printf("hosts(%v): %v\n", firstQuestion, addrs)
	match := findFirstIPv4(addrs)
	if (match != "") {
		glog.Infof("match: %v -> %v\n", firstQuestion, match)
		a := net.ParseIP(match)
		glog.Infof("addr: %v\n", a.To4())

		rr := newRR(firstQuestion, a)
		m.Answer = append(m.Answer, rr)
		m.Authoritative = true
		//	} else {
		//		import 	"github.com/bogdanovich/dns_resolver"
		//		glog.Infof("fallback via resolv.conf...")
		//		resolver, _ := dns_resolver.NewFromResolvConf("/etc/resolv.conf")
		//		result, _ := resolver.LookupHost(firstQuestion)
		//		glog.Infof("...result: %v\n", result)
		//		if len(result) > 0 {
		//			rr = newRR(firstQuestion, result[0])
		//			m.Answer = append(m.Answer, rr)
		//			m.Authoritative = true
		//		}
	}

	glog.Infof("response\n%v\n", m.String())

	w.WriteMsg(m)
}

func main() {
	port := flag.Int("port", 5353, "port to run on")
	flag.Parse()

	dns.HandleFunc(".", handleFirstHost)

	serve := func(net string, port int) {
		server := &dns.Server{Addr: ":" + strconv.Itoa(port), Net: net}
		err := server.ListenAndServe()
		if err != nil {
			glog.Infof("Failed to setup the " + net + " server: %s\n", err.Error())
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
