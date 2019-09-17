package main

import (
	"fmt"
	"net/url"
	"strings"

	"encoding/json"
)

type BoolHolder struct {
	Value bool
}

func (b *BoolHolder) MarshalBinary() ([]byte, error) {
	return json.Marshal(b)
}

func (b *BoolHolder) UnmarshalBinaryData(data []byte) (newdata []byte, err error) {
	err = json.Unmarshal(data, b)
	return
}

func (b *BoolHolder) UnmarshalBinary(data []byte) (err error) {
	_, err = b.UnmarshalBinaryData(data)
	return
}

func MarshalStringToBytes(str string, maxlength int) ([]byte, error) {
	if len(str) > maxlength {
		return nil, fmt.Errorf("Length of string is too long, found length is %d, max length is %d",
			len(str), maxlength)
	}

	data := []byte(str)
	for i := 0; i < len(data); i++ {
		if data[i] == 0x00 {
			// Naughty, Naughty, Naughty
			data[i] = 0x01
		}
	}
	data = append(data, 0x00)

	return data, nil
}

func UnmarshalStringFromBytes(data []byte, maxlength int) (resp string, err error) {
	resp, _, err = UnmarshalStringFromBytesData(data, maxlength)
	return
}

func UnmarshalStringFromBytesData(data []byte, maxlength int) (resp string, newData []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("A panic has occurred while unmarshaling: %s", r)
			return
		}
	}()

	newData = data
	end := -1
	if len(data)-1 < maxlength {
		maxlength = len(data) - 1
	}

	for i := 0; i <= maxlength; i++ {
		if newData[i] == 0x00 {
			// found null terminator
			end = i
			break
		}
	}

	if end == -1 {
		err = fmt.Errorf("Could not find a 0x00 byte before max length + 1")
		return
	}

	resp = string(newData[:end])
	newData = newData[end+1:]
	return
}

// SanitizeFactomdLocation sanitizes user input for the factdom endpoint
// Accepts any string and attempts to parse scheme, host, port, and path
// returns a well-formated URL of scheme://host[:port][/path], removing
// any trailing slash
func SanitizeFactomdLocation(input string) (string, error) {
	if strings.Index(input, "://") == -1 {
		input = "http://" + input
	}

	parsed, err := url.Parse(input)
	if err != nil {
		return "", err
	}

	if len(parsed.Path) > 0 && parsed.Path[len(parsed.Path)-1:] == "/" {
		parsed.Path = parsed.Path[0 : len(parsed.Path)-1]
	}

	return fmt.Sprintf("%s://%s%s", parsed.Scheme, parsed.Host, parsed.Path), nil
}
