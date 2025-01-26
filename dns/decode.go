package dns

import (
	"bytes"
	"encoding/binary"
)

func DecodeName(data []byte, offset int) (string, int, error) {
	var name string
	reader := bytes.NewReader(data[offset:])
	for {
		var length byte
		err := binary.Read(reader, binary.BigEndian, &length)
		if err != nil {
			return "", offset, err
		}
		if length == 0 {
			break
		}
		part := make([]byte, length)
		_, err = reader.Read(part)
		if err != nil {
			return "", offset, err
		}
		name += string(part) + "."
	}
	return name[:len(name)-1], offset + int(reader.Size()) - int(reader.Len()), nil
}
