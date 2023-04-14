package parser

import (
	"bytes"
	"path/filepath"
	"text/template"

	"github.com/gin-admin/gin-admin-cli/v10/internal/utils"
)

const (
	ModsPrefix = "internal/mods"
)

const (
	FileForMods         = "mods.go"
	FileForModuleMain   = "{{lower .ModuleName}}/main.go"
	FileForModuleWire   = "{{lower .ModuleName}}/wire.go"
	FileForModuleAPI    = "{{lower .ModuleName}}/api/{{lowerUnderline .StructName}}.api.go"
	FileForModuleBiz    = "{{lower .ModuleName}}/biz/{{lowerUnderline .StructName}}.biz.go"
	FileForModuleDAL    = "{{lower .ModuleName}}/dal/{{lowerUnderline .StructName}}.dal.go"
	FileForModuleSchema = "{{lower .ModuleName}}/schema/{{lowerUnderline .StructName}}.schema.go"
)

func ParseFilePathFromTpl(moduleName, structName string, tpls ...string) ([]string, error) {
	var paths []string
	for _, tpl := range tpls {
		t := template.Must(template.New("").Funcs(utils.FuncMap).Parse(tpl))
		buf := new(bytes.Buffer)
		if err := t.Execute(buf, map[string]interface{}{
			"ModuleName": moduleName,
			"StructName": structName,
		}); err != nil {
			return nil, err
		}
		paths = append(paths, filepath.Join(ModsPrefix, buf.String()))
	}

	return paths, nil
}

func GetModuleMainFilePath(moduleName string) (string, error) {
	paths, err := ParseFilePathFromTpl(moduleName, "", FileForModuleMain)
	if err != nil {
		return "", err
	}

	return paths[0], nil
}

func GetModuleWireFilePath(moduleName string) (string, error) {
	paths, err := ParseFilePathFromTpl(moduleName, "", FileForModuleWire)
	if err != nil {
		return "", err
	}

	return paths[0], nil
}

func GetModsFilePath() string {
	return filepath.Join(ModsPrefix, FileForMods)
}
