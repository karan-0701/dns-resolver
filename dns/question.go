package dns

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type DNSQuestion struct {
	name   []byte
	type_  int16
	class_ int16
}

func NewDNSQuestion(name []byte, type_ int16, class_ int16) *DNSQuestion {
	return &DNSQuestion{
		name:   name,
		type_:  type_,
		class_: class_,
	}
}

func questionToBytes(question DNSQuestion) ([]byte, error) {
	buf := new(bytes.Buffer)

	_, err := buf.Write(question.name)
	if err != nil {
		return nil, err
	}

	fields := []interface{}{
		question.class_,
		question.type_,
	}

	for _, field := range fields {
		err := binary.Write(buf, binary.BigEndian, field)
		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func parseQuestion(reader io.ReadSeeker) (*DNSQuestion, error) {
	name, err := decodeNameWithCompression(reader)
	if err != nil {
		return nil, fmt.Errorf("name decode error: %v", err)
	}

	data := make([]byte, 4)
	n, err := reader.Read(data)
	if err != nil {
		return nil, fmt.Errorf("failed to read question data: %v (read %d bytes)", err, n)
	}

	type_ := binary.BigEndian.Uint16(data[0:2])
	class_ := binary.BigEndian.Uint16(data[2:4])

	result := NewDNSQuestion(name, int16(type_), int16(class_))
	return result, nil
}
