package main

import (
	"encoding/base64"
	"fmt"
)

func main() {
	msg := "This is totally fun get hands-on learning it from the ground up. Thank you for sharing this info with me and helping me"
	encoded := encode(msg)
	fmt.Println(encoded)

	s, err := decode(encoded)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(s)
}

func encode(msg string) string {
	encoded := base64.URLEncoding.EncodeToString([]byte(msg))
	return string(encoded)
}

func decode(encoded string) (string, error) {
	s, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("Couldn't decode string: %w", err)
	}

	return string(s), nil
}
