package dns

import (
	"bytes"
	"fmt"
)

type DNSPacket struct {
	header      DNSHeader
	questions   []DNSQuestion
	answers     []DNSRecord
	authorities []DNSRecord
	additionals []DNSRecord
}

func NewDNSPacket(header DNSHeader, questions []DNSQuestion, answers []DNSRecord, authorities []DNSRecord, additionals []DNSRecord) *DNSPacket {
	return &DNSPacket{
		header:      header,
		questions:   questions,
		answers:     answers,
		authorities: authorities,
		additionals: additionals,
	}
}

func parseDNSPacket(response []byte) (*DNSPacket, error) {
	header, err := parseHeader(response)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(response[12:])

	questions := make([]DNSQuestion, header.NumQuestions)
	for i := 0; i < int(header.NumQuestions); i++ {
		if reader.Len() == 0 {
			return nil, fmt.Errorf("insufficient bytes for answers")
		}
		question, err := parseQuestion(reader)
		if err != nil {
			return nil, err
		}
		questions[i] = *question
	}

	answers := make([]DNSRecord, header.NumAnswers)
	for i := 0; i < int(header.NumAnswers); i++ {
		if reader.Len() == 0 {
			return nil, fmt.Errorf("insufficient bytes for answers")
		}
		record, err := parseRecordWithCompression(reader)
		if err != nil {
			return nil, err
		}
		answers[i] = *record
	}

	authorities := make([]DNSRecord, header.NumAuthorities)
	for i := 0; i < int(header.NumAuthorities); i++ {
		if reader.Len() == 0 {
			return nil, fmt.Errorf("insufficient bytes for answers")
		}
		record, err := parseRecordWithCompression(reader)
		if err != nil {
			return nil, err
		}
		authorities[i] = *record
	}

	additionals := make([]DNSRecord, header.NumAdditionals)
	for i := 0; i < int(header.NumAdditionals); i++ {
		if reader.Len() == 0 {
			return nil, fmt.Errorf("insufficient bytes for answers")
		}
		record, err := parseRecordWithCompression(reader)
		if err != nil {
			return nil, err
		}
		additionals[i] = *record
	}

	return NewDNSPacket(*header, questions, answers, authorities, additionals), nil
}
