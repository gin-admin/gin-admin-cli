package parser

import (
	"context"
	"testing"
)

const dir = "/home/lyric/go/src/gin-admin"

func TestModifyModuleMainFile(t *testing.T) {
	buf, err := ModifyModuleMainFile(context.Background(), BasicArgs{
		Dir:        dir,
		ModuleName: "RBAC",
		StructName: "Role",
		Flag:       AstFlagGen,
	})
	if err != nil {
		t.Fatal(err)
	}
	_ = buf
	// fmt.Println(string(buf))
}

func TestModifyModuleWireFile(t *testing.T) {
	buf, err := ModifyModuleWireFile(context.Background(), BasicArgs{
		Dir:         dir,
		ModuleName:  "RBAC",
		StructName:  "Role",
		Flag:        AstFlagGen,
		GenPackages: []string{"dal", "biz", "api"},
	})
	if err != nil {
		t.Fatal(err)
	}

	_ = buf
	// fmt.Println(string(buf))
}

func TestModifyModsFile(t *testing.T) {
	buf, err := ModifyModsFile(context.Background(), BasicArgs{
		Dir:        dir,
		ModuleName: "Configurator",
		Flag:       AstFlagGen,
	})
	if err != nil {
		t.Fatal(err)
	}

	_ = buf
	// fmt.Println(string(buf))
}
