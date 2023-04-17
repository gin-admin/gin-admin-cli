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
				GormTag: "size:128;index",
				Query: &FieldQuery{
					Name:    "LikeName",
					InQuery: true,
					FormTag: "name",
					OP:      "LIKE",
				},
				Form: &FieldForm{
					BindingTag: "required",
				},
			},
			{
				Name:    "Description",
				Type:    "string",
				Comment: "Details about menu",
				GormTag: "type:text",
				Form: &FieldForm{
					Name:    "Desc",
					JSONTag: `,omitempty`,
				},
				CustomTag: `form:"desc" validate:"max=255"`,
			},
			{
				Name:    "Sequence",
				Type:    "int",
				Comment: "Sequence for sorting",
				GormTag: "index;",
				Form:    &FieldForm{},
			},
			{
				Name:    "Type",
				Type:    "string",
				Comment: "Type of menu (group/menu/button)",
				GormTag: "size:32;index;",
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
