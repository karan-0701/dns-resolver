package dns

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

func encodeDnsName(domainName string) []byte {
	var buffer bytes.Buffer
	domainSplit := strings.Split(domainName, ".")

	for i, part := range domainSplit {
		if i == 0 {
			continue
		}
		buffer.WriteByte(byte(len(part)))
		buffer.WriteString(part)
	}
	buffer.WriteByte(0)

	return buffer.Bytes()
}

func decodeNameWithCompression(reader io.ReadSeeker) ([]byte, error) {
	var parts [][]byte
	seenPointers := make(map[int64]bool)

	for {
		lengthByte := make([]byte, 1)
		initialPos, _ := reader.Seek(0, io.SeekCurrent)

		_, err := reader.Read(lengthByte)
		if err != nil {
			return nil, fmt.Errorf("length read error at position %d: %v", initialPos, err)
		}

		length := lengthByte[0]
		if length == 0 {
			break
		}

		if length&0b1100_0000 != 0 {
			pointer := uint16(length&0b00111111) << 8

			nextByte := make([]byte, 1)
			_, err := reader.Read(nextByte)
			if err != nil {
				return nil, err
			}
			pointer |= uint16(nextByte[0])

			if seenPointers[int64(pointer)] {
				return nil, fmt.Errorf("circular name compression detected")
			}
			seenPointers[int64(pointer)] = true

			currentPos, _ := reader.Seek(0, io.SeekCurrent)

			_, err = reader.Seek(int64(pointer), io.SeekStart)
			if err != nil {
				return nil, err
			}

			compressedPart, err := decodeNameWithCompression(reader)
			if err != nil {
				return nil, err
			}
			parts = append(parts, compressedPart)

			_, err = reader.Seek(currentPos, io.SeekStart)
			if err != nil {
				return nil, err
			}
			break
		} else {
			part := make([]byte, length)
			_, err = io.ReadFull(reader, part)
			if err != nil {
				return nil, err
			}
			parts = append(parts, part)
		}
	}

	return bytes.Join(parts, []byte(".")), nil
}

func decodeCompressedName(length byte, reader io.ReadSeeker) ([]byte, error) {
	pointer := uint16(length&0b00111111) << 8

	nextByte := make([]byte, 1)
	_, err := reader.Read(nextByte)
	if err != nil {
		return nil, fmt.Errorf("compressed pointer read error: %v", err)
	}
	pointer |= uint16(nextByte[0])

	currentPos, err := reader.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, fmt.Errorf("current position error: %v", err)
	}

	_, err = reader.Seek(int64(pointer), io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf("seek to pointer error: %v", err)
	}

	result, err := decodeNameWithCompression(reader)
	if err != nil {
		return nil, fmt.Errorf("recursive name decode error: %v", err)
	}

	_, err = reader.Seek(int64(currentPos), io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf("seek back error: %v", err)
	}

	return result, nil
}

func ipToString(ipBytes []byte) string {
	ipParts := make([]string, len(ipBytes))

	for i, b := range ipBytes {
		ipParts[i] = fmt.Sprintf("%d", b)
	}
	ipAddress := strings.Join(ipParts, ".")

	return ipAddress
}
