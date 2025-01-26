package main

import (
	"fmt"

	"dns-resolver/dns"
)

func main() {
	domainName := "www.google.com"
	recordType := dns.TYPE_A

	ip, err := dns.Resolve(domainName, recordType)
	if err != nil {
		fmt.Println("Error resolving domain:", err)
		return
	}

	// Print the resolved IP address
	fmt.Printf("Resolved IP address for %s: %s\n", domainName, string(ip))
}
