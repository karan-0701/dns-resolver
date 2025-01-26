package dns

import (
	"bytes"
	"math/rand"
)

const CLASS_IN = 1

func buildQuery(domainName string, recordType int16) ([]byte, error) {
	seed := int64(1)
	r := rand.New(rand.NewSource(seed))
	id := r.Intn(65535)

	nameEncoded := encodeDnsName(domainName)
	RECURSION_DESIRED := 0
	header := NewDNSHeader(uint16(id), uint16(RECURSION_DESIRED), 1, 0, 0, 0)
	question := NewDNSQuestion(nameEncoded, recordType, CLASS_IN)

	buf := new(bytes.Buffer)
	dataHeader, err := headerToBytes(*header)
	if err != nil {
		return nil, err
	}
	buf.Write(dataHeader)
	dataQuestion, err := questionToBytes(*question)
	if err != nil {
		return nil, err
	}
	buf.Write(dataQuestion)

	return buf.Bytes(), nil
}
