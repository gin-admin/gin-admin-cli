package parser

import (
	"bytes"
	"text/template"

	"github.com/gin-admin/gin-admin-cli/v10/internal/utils"
)

const (
	FileForMods         = "mods.go"
	FileForModuleMain   = "{{lower .ModuleName}}/main.go"
	FileForModuleWire   = "{{lower .ModuleName}}/wire.go"
	FileForModuleAPI    = "{{lower .ModuleName}}/api/{{lowerUnderline .StructName}}.api.go"
	FileForModuleBiz    = "{{lower .ModuleName}}/biz/{{lowerUnderline .StructName}}.biz.go"
	FileForModuleDAL    = "{{lower .ModuleName}}/dal/{{lowerUnderline .StructName}}.dal.go"
	FileForModuleSchema = "{{lower .ModuleName}}/schema/{{lowerUnderline .StructName}}.go"
)

func ParseFilePathFromTpl(moduleName, structName string, tpl string) (string, error) {
	t := template.Must(template.New("").Funcs(utils.FuncMap).Parse(tpl))
	buf := new(bytes.Buffer)
	if err := t.Execute(buf, map[string]interface{}{
		"ModuleName": moduleName,
		"StructName": structName,
	}); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func GetModuleMainFilePath(moduleName string) (string, error) {
	p, err := ParseFilePathFromTpl(moduleName, "", FileForModuleMain)
	if err != nil {
		return "", err
	}
	return p, nil
}

func GetModuleWireFilePath(moduleName string) (string, error) {
	p, err := ParseFilePathFromTpl(moduleName, "", FileForModuleWire)
	if err != nil {
		return "", err
	}
	return p, nil
}
