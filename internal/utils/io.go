package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"go/format"
	"io"
	"os"
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

// FmtGoFile formats the given Go source code file
func FmtGoFile(name string) error {
	// read the contents of the file
	content, err := os.ReadFile(name)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	// format the source code
	formatted, err := format.Source(content)
	if err != nil {
		return fmt.Errorf("error formatting file: %v", err)
	}

	// overwrite the existing file with the formatted code
	err = WriteFile(name, formatted)
	if err != nil {
		return fmt.Errorf("error writing formatted file: %v", err)
	}

	return nil
}
