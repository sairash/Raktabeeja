package helper

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strings"
)

func Base64Encode(output []byte) string {
	return strings.Replace(base64.StdEncoding.EncodeToString(output), "/", "-", -1)
}

func ReturnBase64(output string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(strings.Replace(output, "-", "/", -1))
}

func SplitBytes(m []byte) ([][]byte, error) {
	n := len(m)
	r := make([]byte, n)
	_, err := rand.Read(r)
	if err != nil {
		return nil, err
	}
	_s0 := make([]byte, n+1)
	_s1 := make([]byte, n+1)
	_s2 := make([]byte, n+1)
	// set tags
	_s0[0] = 0
	_s1[0] = 1
	_s2[0] = 2
	for i := 0; i < n; i++ {
		_s0[i+1] = ((m[i] & 0xf0) >> 4) ^ r[i]
		_s1[i+1] = ((m[i] & 0x0f) << 4) ^ r[i]
		_s2[i+1] = m[i] ^ r[i]
	}
	return [][]byte{_s0, _s1, _s2}, nil
}

func JoinBytes01New(a, b []byte) []byte {
	m := make([]byte, len(a))
	for i := 0; i < len(m); i++ {
		c := a[i] ^ b[i]
		m[i] = ((c << 4) & 0xf0) | ((c >> 4) & 0x0f)
	}
	return m
}

func JoinBytes(a, b []byte) []byte {
	if len(a) != len(b) {
		return nil
	}
	if len(a) < 1 {
		return nil
	}
	if a[0] > b[0] {
		a, b = b, a
	}
	m := make([]byte, len(a)-1)
	if a[0] == 0 && b[0] == 1 {
		joinBytes01(m, a[1:], b[1:])
	} else if a[0] == 0 && b[0] == 2 {
		joinBytes02(m, a[1:], b[1:])
	} else if a[0] == 1 && b[0] == 2 {
		joinBytes12(m, a[1:], b[1:])
	} else {
		return nil
	}
	return m
}

func joinBytes01(m, a, b []byte) {
	for i := 0; i < len(m); i++ {
		c := a[i] ^ b[i]
		m[i] = ((c << 4) & 0xf0) | ((c >> 4) & 0x0f)
	}
}

func joinBytes02(m, a, b []byte) {
	for i := 0; i < len(m); i++ {
		c := a[i] ^ b[i]
		m[i] = ((c & 0xf0) >> 4) ^ c
	}
}

func joinBytes12(m, a, b []byte) {
	for i := 0; i < len(m); i++ {
		c := a[i] ^ b[i]
		m[i] = ((c & 0x0f) << 4) ^ c
	}
}

func Hash(data string) string {
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}

func HashBytes(data []byte) string {
	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:])
}

func Xor(input string, key int) string {
	result := ""
	for _, char := range input {
		xoredChar := rune(char) ^ rune(key)
		result += fmt.Sprintf("%02x", xoredChar)
	}
	return result
}

func XorBack(result string, key int) string {
	decodedBytes, err := hex.DecodeString(result)
	if err != nil {
		return ""
	}

	originalChars := make([]rune, len(decodedBytes))
	for i, b := range decodedBytes {
		originalChars[i] = rune(b) ^ rune(key)
	}

	return string(originalChars)
}

func Int64ToBytes(n int64) ([]byte, error) {
	var buf bytes.Buffer

	err := binary.Write(&buf, binary.LittleEndian, n)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
