package fs

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"

	"github.com/gin-admin/gin-admin-cli/v10/internal/utils"
)

var Ins FS = NewEmbedFS()

func SetIns(ins FS) {
	Ins = ins
}

type FS interface {
	ReadFile(name string) ([]byte, error)
	ParseTplData(name string, data interface{}) ([]byte, error)
}

type osFS struct {
	dir string
}

func NewOSFS(dir string) FS {
	return &osFS{dir: dir}
}

func (fs *osFS) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(filepath.Join(fs.dir, name))
}

func (fs *osFS) ParseTplData(name string, data interface{}) ([]byte, error) {
	tplBytes, err := fs.ReadFile(name)
	if err != nil {
		return nil, err
	}
	return parseTplData(string(tplBytes), data)
}

func parseTplData(text string, data interface{}) ([]byte, error) {
	t := template.Must(template.New("").Funcs(utils.FuncMap).Parse(text))
	buf := new(bytes.Buffer)
	if err := t.Execute(buf, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
