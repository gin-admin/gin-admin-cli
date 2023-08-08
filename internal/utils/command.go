package utils

import (
	"fmt"
	"go/format"
	"os"
	"os/exec"
	"strings"

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
func ExecGoImports(dir, name string) error {
	localPath, err := exec.LookPath("goimports")
	if err != nil {
		if err := ExecGoInstall(dir, "golang.org/x/tools/cmd/goimports@latest"); err != nil {
			return nil
		}
	}

	cmd := exec.Command(localPath, "-w", name)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	return cmd.Run()
}

func ExecGoInstall(dir, path string) error {
	localPath, err := exec.LookPath("go")
	if err != nil {
		zap.S().Warn("Not found go command, please install go first")
		return nil
	}

	cmd := exec.Command(localPath, "install", path)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	return cmd.Run()
}

func ExecGoModTidy(dir string) error {
	localPath, err := exec.LookPath("go")
	if err != nil {
		zap.S().Warn("Not found go command, please install go first")
		return nil
	}

	cmd := exec.Command(localPath, "mod", "tidy")
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	return cmd.Run()
}

// Executes the wire command on the given file
func ExecWireGen(dir, path string) error {
	_, err := exec.LookPath("wire")
	if err != nil {
		if err := ExecGoInstall(dir, "github.com/google/wire/cmd/wire@latest"); err != nil {
			return nil
		}
	}

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
		if err := ExecGoInstall(dir, "github.com/swaggo/swag/cmd/swag@latest"); err != nil {
			return nil
		}
	}

	zap.S().Infof(fmt.Sprintf("%s init --parseDependency --generalInfo %s --output %s", localPath, generalInfo, output))
	cmd := exec.Command("swag", "init", "--parseDependency", "--generalInfo", generalInfo, "--output", output)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	return cmd.Run()
}

func ExecGitInit(dir string) error {
	localPath, err := exec.LookPath("git")
	if err != nil {
		zap.S().Warn("Not found git command, please install git first")
		return nil
	}

	cmd := exec.Command(localPath, "init")
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	return cmd.Run()
}

func ExecGitClone(dir, url, branch, name string) error {
	localPath, err := exec.LookPath("git")
	if err != nil {
		zap.S().Warn("Not found git command, please install git first")
		return nil
	}

	var args []string
	args = append(args, "clone")
	args = append(args, url)
	if branch != "" {
		args = append(args, "-b")
		args = append(args, branch)
	}
	if name != "" {
		args = append(args, name)
	}

	zap.S().Infof("git %s", strings.Join(args, " "))
	cmd := exec.Command(localPath, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	return cmd.Run()
}

func ExecTree(dir string) error {
	localPath, err := exec.LookPath("tree")
	if err != nil {
		return nil
	}

	cmd := exec.Command(localPath, "-L", "2", "-I", ".git")
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	return cmd.Run()
}
