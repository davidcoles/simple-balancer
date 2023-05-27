package main

import (
	"log"
	"time"

	"github.com/davidcoles/vc5"
	"github.com/davidcoles/vc5/bgp4"
	"github.com/davidcoles/vc5/config"
)

const ASN = 65000 // Local autonomous system number to use for BGP sessions

func main() {
	hostip := "10.11.12.13"                            // IP address of the server running as load-balancer
	virtual := "192.168.100.1"                         // Virtual IP address of the service
	bgppeers := []string{"10.11.12.10"}                // BGP peers to advertise virtual IP address to - can be empty
	backends := []string{"10.11.12.20", "10.11.12.21"} // Backend servers running a webserver
	interfaces := []string{"eth0"}                     // Interface to use on the load-balancer server

	service := config.Service{RIPs: map[string]config.Checks{}} // Create a service object

	for _, rip := range backends {
		service.RIPs[rip] = config.Checks{} // Add each backend server to the service with empty checks
	}

	conf := config.Config{
		VIPs: map[string]map[string]config.Service{
			virtual: map[string]config.Service{"tcp:80": service}, // configure VIP to send HTTP to the service
		},
	}

	pool := bgp4.Pool{
		Address: hostip,   // Outgoing BGP connections bind to this IP
		Peers:   bgppeers, // advertise VIP to these hosts - typically routers
		ASN:     ASN,      // Local autonomous system number
	}

	lb := vc5.LoadBalancer{
		Native:     false,      // If your ethernet card has driver support for XDP this can be set to true (FAST!)
		KillSwitch: 5,          // Stop all load-balancing activity after 3 minutes as a safety precaution
		Interfaces: interfaces, // Load XDP/eBPF code to these interfaces
	}

	// Start BGP sessions with peers - if any
	if !pool.Open() {
		log.Fatal("BGP")
	}

	manifest, err := vc5.Load(&conf) // Sanity check the configuration
	if err != nil {
		log.Fatal("Conf: ", err)
	}

	err = lb.Start(hostip, manifest) // Start the load-balancer!
	if err != nil {
		log.Fatal("Balancer: ", err)
	}

	defer lb.Close() // Unload all the XDP/eBPF code on exit

	pool.NLRI(map[string]bool{virtual: true}) // advertise the VIP

	time.Sleep(6 * time.Minute) // Kick our heels for a while ...
}
