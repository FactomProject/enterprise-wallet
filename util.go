package main

import (
	"fmt"

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
