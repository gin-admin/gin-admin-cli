package generate

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

const (
	delimiter = "\n"
)

// 获取模块头
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

// 写入文件
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
		return fmt.Errorf("文件[%s]已经存在", name)
	}

	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()

	io.Copy(file, buf)
	return nil
}

// 执行go代码格式化
func execGoFmt(name string) error {
	cmd := exec.Command("gofmt", "-w", name, name)
	return cmd.Run()
}

// 执行解析模板
func execParseTpl(tpl string, data interface{}) (*bytes.Buffer, error) {
	t := template.Must(template.New("").Parse(tpl))

	buf := new(bytes.Buffer)
	err := t.Execute(buf, data)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

// 读取文件内容
func readFile(name string) (*bytes.Buffer, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buf := new(bytes.Buffer)
	io.Copy(buf, file)
	return buf, nil
}

// 写入文件
func writeFile(name string, buf *bytes.Buffer) error {
	file, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	io.Copy(file, buf)
	return nil
}

// 插入文件内容
func insertFileContent(name string, startPrefix, endContain, content string, excludeEnds ...string) error {
	buf, err := readFile(name)
	if err != nil {
		return err
	}

	pstart := false
	start := false
	end := false
	complete := false

	var pline string

	nbuf := new(bytes.Buffer)
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		cline := scanner.Text()

		if pstart {
			nbuf.WriteString(pline)
			nbuf.WriteString(delimiter)
		}
		pstart = true
		pline = cline

		if complete {
			continue
		}

		tline := strings.TrimSpace(cline)
		if !start && strings.HasPrefix(tline, startPrefix) {
			start = true
		}

		exclude := tline == ""
		if !exclude {
			for _, e := range excludeEnds {
				if strings.HasPrefix(tline, e) {
					exclude = true
					break
				}
			}
		}

		if !end && start && !exclude && strings.Contains(tline, endContain) {
			end = true
		}

		if !(!end || exclude || strings.Contains(tline, endContain)) {
			nbuf.WriteString(content)
			complete = true
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	nbuf.WriteString(pline)

	return writeFile(name, nbuf)
}
