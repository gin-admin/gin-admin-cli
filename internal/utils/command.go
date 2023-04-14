package utils

import (
	"fmt"
	"os/exec"

	"go.uber.org/zap"
)

// Executes the goimports command on the given file
func ExecGoImports(name string) error {
	localPath, err := exec.LookPath("goimports")
	if err != nil {
		zap.S().Errorf("goimports not found: %v", err)
		return nil
	}

	zap.S().Infof(fmt.Sprintf("%s -w %s", localPath, name))
	cmd := exec.Command(localPath, "-w", name)
	return cmd.Run()
}
