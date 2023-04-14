package utils

import (
	"fmt"
	"go/format"
	"os"
	"os/exec"

	"go.uber.org/zap"
)

// Formats the given Go source code file
func ExecGoFormat(name string) error {
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

// Executes the goimports command on the given file
func ExecGoImports(name string) error {
	localPath, err := exec.LookPath("goimports")
	if err != nil {
		// zap.S().Warnf("goimports not found: %v", err)
		return nil
	}

	// zap.S().Infof(fmt.Sprintf("%s -w %s", localPath, name))
	cmd := exec.Command(localPath, "-w", name)
	return cmd.Run()
}

// Executes the wire command on the given file
func ExecWireGen(dir, path string) error {
	localPath, err := exec.LookPath("wire")
	if err != nil {
		// zap.S().Warnf("wire not found: %v", err)
		return nil
	}

	zap.S().Infof(fmt.Sprintf("%s gen %s", localPath, path))
	cmd := exec.Command("wire", "gen", "./"+path)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	return cmd.Run()
}

// Executes the swag command on the given file
func ExecSwagGen(dir, generalInfo, output string) error {
	localPath, err := exec.LookPath("swag")
	if err != nil {
		// zap.S().Warnf("swag not found: %v", err)
		return nil
	}

	zap.S().Infof(fmt.Sprintf("%s init --parseDependency --generalInfo %s --output %s", localPath, generalInfo, output))
	cmd := exec.Command("swag", "init", "--parseDependency", "--generalInfo", generalInfo, "--output", output)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	return cmd.Run()
}
