package generate

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

const (
	delimiter = "\n"
)

var ErrFileExists = errors.New("file has been exists")

func getModuleHeader(moduleName string, imports ...string) *bytes.Buffer {
	buf := new(bytes.Buffer)

	buf.WriteString(fmt.Sprintf("package %s", moduleName))
	buf.WriteString(delimiter)
	buf.WriteString(delimiter)

	if len(imports) > 0 {
		buf.WriteString("import (")
		buf.WriteString(delimiter)

		for _, s := range imports {
			buf.WriteByte('\t')
			buf.WriteString(s)
			buf.WriteString(delimiter)
		}

		buf.WriteByte(')')
		buf.WriteString(delimiter)
		buf.WriteString(delimiter)
	}

	return buf
}

func createFile(ctx context.Context, name string, buf *bytes.Buffer) error {
	exists := true
	_, err := os.Stat(name)
	if err != nil {
		if os.IsNotExist(err) {
			exists = false
		} else {
			return err
		}
	}

	if exists {
		fmt.Printf("File has been exists: %s, skip", name)
		return ErrFileExists
	}

	os.MkdirAll(filepath.Dir(name), 0777)
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()

	_, _ = io.Copy(file, buf)
	return nil
}

func execGoFmt(name string) error {
	cmd := exec.Command("gofmt", "-w", name, name)
	return cmd.Run()
}

func execParseTpl(tpl string, data interface{}) (*bytes.Buffer, error) {
	t := template.Must(template.New("").Parse(tpl))

	buf := new(bytes.Buffer)
	err := t.Execute(buf, data)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func readFile(name string) (*bytes.Buffer, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buf := new(bytes.Buffer)
	_, _ = io.Copy(buf, file)
	return buf, nil
}

func writeFile(name string, buf *bytes.Buffer) error {
	file, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	_, _ = io.Copy(file, buf)
	return nil
}

// insertContent 插入文件内容
// fn 回调当前行数据，返回-1为插入当前行之前，1为插入当前行之后
func insertContent(name string, fn func(string) (string, int, bool)) error {
	buf, err := readFile(name)
	if err != nil {
		return err
	}

	nbuf := new(bytes.Buffer)
	scanner := bufio.NewScanner(buf)

	for scanner.Scan() {
		cline := scanner.Text()
		data, flag, ok := fn(cline)
		if ok {
			if flag == -1 {
				nbuf.WriteString(data)
				nbuf.WriteString(delimiter)
				nbuf.WriteString(cline)
				nbuf.WriteString(delimiter)
				continue
			}
			nbuf.WriteString(cline)
			nbuf.WriteString(delimiter)
			nbuf.WriteString(data)
			nbuf.WriteString(delimiter)
			continue
		}
		nbuf.WriteString(cline)
		nbuf.WriteString(delimiter)
	}

	return writeFile(name, nbuf)
}
