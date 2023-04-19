package utils

import (
	"bufio"
	"bytes"
	"io"
	"os"

	jsoniter "github.com/json-iterator/go"
	"gopkg.in/yaml.v2"
)

// Scanner scans the given reader line by line and calls the given function
func Scanner(r io.Reader, fn func(string) string) *bytes.Buffer {
	buf := new(bytes.Buffer)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		buf.WriteString(fn(line))
		buf.WriteString("\n")
	}
	return buf
}

// ExistsFile checks if the given file exists
func ExistsFile(name string) (bool, error) {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

// WriteFile writes the given data to the given file
func WriteFile(name string, data []byte) error {
	return os.WriteFile(name, data, 0644)
}

// Parses the given JSON file
func ParseJSONFile(name string, obj interface{}) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()

	return jsoniter.NewDecoder(f).Decode(obj)
}

// Parses the given YAML file
func ParseYAMLFile(name string, obj interface{}) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()

	return yaml.NewDecoder(f).Decode(obj)
}

// Checks if the given path is a directory
func IsDir(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		return false
	}
	return info.IsDir()
}
