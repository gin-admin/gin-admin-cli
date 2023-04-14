package schema

import (
	"fmt"
	"testing"

	jsoniter "github.com/json-iterator/go"
	"gopkg.in/yaml.v2"
)

func TestMarshalS(t *testing.T) {
	s := &S{
		Name: "Menu",
		Fields: []*Field{
			{
				Name:    "Name",
				Type:    "string",
				Comment: "Display name of menu",
				GormTag: &FieldGormTag{Index: true, Size: 128},
				Query: &FieldQuery{
					Name:    "LikeName",
					InQuery: true,
					FormTag: "name",
					OP:      "LIKE",
				},
				Form: &FieldForm{
					BindingTag: &FieldBindingTag{},
				},
			},
			{
				Name:    "Description",
				Type:    "string",
				Comment: "Details about menu",
				GormTag: &FieldGormTag{Size: 1024},
				Form:    &FieldForm{},
			},
			{
				Name:    "Sequence",
				Type:    "int",
				Comment: "Sequence for sorting",
				GormTag: &FieldGormTag{Index: true},
				Form:    &FieldForm{},
			},
			{
				Name:    "Type",
				Type:    "string",
				Comment: "Type of menu (group/menu/button)",
				GormTag: &FieldGormTag{Size: 20},
				Form:    &FieldForm{},
			},
		},
		DisablePagination: true,
	}

	buf, err := jsoniter.MarshalIndent([]*S{s}, "", "")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("\n" + string(buf) + "\n")

	fmt.Println("=====================================")

	ybuf, err := yaml.Marshal([]*S{s})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("\n" + string(ybuf) + "\n")
}
