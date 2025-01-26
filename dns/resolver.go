package dns

import (
	"fmt"
	"net"
	"strings"
)

const TYPE_A = 1
const TYPE_NS = 2

func lookupDomain(domainName string) (string, error) {
	query, err := buildQuery(domainName, TYPE_A)
	if err != nil {
		return "", err
	}

	conn, err := net.Dial("udp", "8.8.8.8.53")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	_, err = conn.Write(query)
	if err != nil {
		return "", err
	}
	response := make([]byte, 1024)
	n, err := conn.Read(response)
	if err != nil {
		return "", err
	}

	fmt.Printf("Received %d bytes: %x\n", n, response[:n])

	responseData, err := parseDNSPacket(response)
	if err != nil {
		return "", err
	}

	ipBytes := responseData.answers[0].data.([]byte)
	ipParts := make([]string, len(ipBytes))
	for i, b := range ipBytes {
		ipParts[i] = fmt.Sprintf("%d", b)
	}
	ipAddress := strings.Join(ipParts, ".")

	fmt.Println(ipAddress)
	return ipAddress, nil
}

func sendQuery(ipAddress string, domainName string, recordType int) (*DNSPacket, error) {
	query, err := buildQuery(domainName, int16(recordType))
	if err != nil {
		return nil, err
	}

	port := 53
	conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", ipAddress, port))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	_, err = conn.Write(query)
	if err != nil {
		return nil, err
	}

	response := make([]byte, 1024)
	n, err := conn.Read(response)
	if err != nil {
		return nil, err
	}

	// Trim the response to actual received bytes
	response = response[:n]

	// fmt.Printf("Received %d bytes\n", n)
	// fmt.Printf("Raw response: %x\n", response)

	return parseDNSPacket(response)
}

func getNameServer(packet DNSPacket) []byte {
	for _, authority := range packet.authorities {
		if authority.type_ == TYPE_NS {
			return authority.data.([]byte)
		}
	}
	return nil
}

func getAnswer(packet DNSPacket) []byte {
	for _, answer := range packet.answers {
		if answer.type_ == TYPE_A {
			switch data := answer.data.(type) {
			case []byte:
				return data
			case string:
				return []byte(data)
			default:
				fmt.Printf("Unexpected type in answer data: %T\n", data)
				return nil
			}
		}
	}
	return nil
}

func getNameServerIp(packet DNSPacket) []byte {
	for _, answer := range packet.additionals {
		if answer.type_ == TYPE_A {
			switch data := answer.data.(type) {
			case []byte:
				return data
			case string:
				return []byte(data)
			default:
				fmt.Printf("Unexpected type in additional data: %T\n", data)
				return nil
			}
		}
	}
	return nil
}

func Resolve(domainName string, recordType int) ([]byte, error) {
	nameServer := "1.1.1.1"
	for {
		fmt.Printf("Querying %s for %s\n", nameServer, domainName)
		response, err := sendQuery(nameServer, domainName, recordType)
		if err != nil {
			return nil, err
		}

		ip := getAnswer(*response)
		if ip != nil {
			return ip, nil
		}
		fmt.Println(string(ip))
		// Check if we have a nameserver IP to continue recursion
		if nsIP := getNameServerIp(*response); nsIP != nil {
			nameServer = string(nsIP)
		} else {
			return nil, fmt.Errorf("something went wrong")
		}
	}
}
