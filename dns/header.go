package dns

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type DNSHeader struct {
	ID             uint16
	Flags          uint16
	NumQuestions   uint16
	NumAnswers     uint16
	NumAuthorities uint16
	NumAdditionals uint16
}

func NewDNSHeader(id uint16, flags uint16, numQuestions uint16, numAnswers uint16, numAuthorities uint16, numAdditionals uint16) *DNSHeader {
	return &DNSHeader{
		ID:             id,
		Flags:          flags,
		NumQuestions:   numQuestions,
		NumAnswers:     numAnswers,
		NumAuthorities: numAuthorities,
		NumAdditionals: numAdditionals,
	}
}

func headerToBytes(header DNSHeader) ([]byte, error) {
	buf := new(bytes.Buffer)

	fields := []interface{}{
		header.ID,
		header.Flags,
		header.NumQuestions,
		header.NumAnswers,
		header.NumAuthorities,
		header.NumAdditionals,
	}

	for _, field := range fields {
		err := binary.Write(buf, binary.BigEndian, field)
		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func parseHeader(data []byte) (*DNSHeader, error) {
	if len(data) < 12 {
		return nil, fmt.Errorf("insufficient header data: %d bytes", len(data))
	}

	reader := bytes.NewReader(data)
	var header DNSHeader

	err := binary.Read(reader, binary.BigEndian, &header)
	if err != nil {
		return nil, fmt.Errorf("header read error: %v", err)
	}

	return &header, nil
}
