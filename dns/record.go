package dns

import (
	"encoding/binary"
	"fmt"
	"io"
)

type DNSRecord struct {
	name   []byte
	type_  uint16
	class_ uint16
	ttl    uint32
	data   interface{}
}

func NewDNSRecord(name []byte, type_ uint16, class_ uint16, ttl uint32, data interface{}) *DNSRecord {
	return &DNSRecord{
		name:   name,
		type_:  type_,
		class_: class_,
		ttl:    ttl,
		data:   data,
	}
}
func parseRecordWithCompression(reader io.ReadSeeker) (*DNSRecord, error) {
	domainName, err := decodeNameWithCompression(reader)
	if err != nil {
		return nil, fmt.Errorf("record name decode error: %v", err)
	}

	nextBytes := make([]byte, 10)
	_, err = reader.Read(nextBytes)
	if err != nil {
		return nil, fmt.Errorf("record metadata read error: %v", err)
	}

	type_ := binary.BigEndian.Uint16(nextBytes[0:2])
	class_ := binary.BigEndian.Uint16(nextBytes[2:4])
	ttl := binary.BigEndian.Uint32(nextBytes[4:8])
	dataLength := binary.BigEndian.Uint16(nextBytes[8:10])

	var data interface{}
	switch type_ {
	case TYPE_NS:
		data, err = decodeNameWithCompression(reader)
		if err != nil {
			return nil, fmt.Errorf("NS record data decode error: %v", err)
		}

	case TYPE_A:
		ipBytes := make([]byte, dataLength)
		_, err := reader.Read(ipBytes)
		if err != nil {
			return nil, fmt.Errorf("A record data read error: %v", err)
		}
		data = ipToString(ipBytes)
	default:
		data = make([]byte, dataLength)
		_, err := reader.Read(data.([]byte))
		if err != nil {
			return nil, fmt.Errorf("record data read error: %v", err)
		}
	}

	return NewDNSRecord(domainName, type_, class_, ttl, data), nil
}
