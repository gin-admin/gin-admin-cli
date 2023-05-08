package parser

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-admin/gin-admin-cli/v10/internal/utils"
	"golang.org/x/mod/modfile"
)

const (
	StructNamingAPI = "API"
	StructNamingBIZ = "BIZ"
	StructNamingDAL = "DAL"
)

const (
	StructPackageAPI    = "api"
	StructPackageBIZ    = "biz"
	StructPackageDAL    = "dal"
	StructPackageSchema = "schema"
)

var StructPackages = []string{
	StructPackageSchema,
	StructPackageDAL,
	StructPackageBIZ,
	StructPackageAPI,
}

var StructPackageTplPaths = map[string]string{
	StructPackageAPI:    FileForModuleAPI,
	StructPackageBIZ:    FileForModuleBiz,
	StructPackageDAL:    FileForModuleDAL,
	StructPackageSchema: FileForModuleSchema,
}

func GetStructAPIName(structName string) string {
	return structName + StructNamingAPI
}

func GetStructBIZName(structName string) string {
	return structName + StructNamingBIZ
}

func GetStructDALName(structName string) string {
	return structName + StructNamingDAL
}

func GetStructRouterVarName(structName string) string {
	return utils.ToLowerCamel(structName)
}

func GetStructRouterGroupName(structName string) string {
	return utils.ToLowerHyphensPlural(structName)
}

func GetModuleImportName(moduleName string) string {
	return strings.ToLower(moduleName)
}

func GetRootImportPath(dir string) string {
	modBytes, err := os.ReadFile(filepath.Join(dir, "go.mod"))
	if err != nil {
		return ""
	}
	return modfile.ModulePath(modBytes)
}

func GetModuleImportPath(dir, modulePath, moduleName string) string {
	return GetRootImportPath(dir) + "/" + modulePath + "/" + GetModuleImportName(moduleName)
}

func GetUtilImportPath(dir, modulePath string) string {
	return GetRootImportPath(dir) + "/pkg/util"
}
