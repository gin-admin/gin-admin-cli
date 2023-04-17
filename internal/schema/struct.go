package schema

import "github.com/gin-admin/gin-admin-cli/v10/internal/utils"

type S struct {
	RootImportPath    string   `yaml:"-" json:"-"`
	ModuleName        string   `yaml:"-" json:"-"`
	ModuleImportPath  string   `yaml:"-" json:"-"`
	UtilsImportPath   string   `yaml:"-" json:"-"`
	Name              string   `yaml:"name,omitempty" json:"name,omitempty"`
	TableName         string   `yaml:"table_name,omitempty" json:"table_name,omitempty"`
	Comment           string   `yaml:"comment,omitempty" json:"comment,omitempty"`
	Fields            []*Field `yaml:"fields,omitempty" json:"fields,omitempty"`
	DisablePagination bool     `yaml:"disable_pagination,omitempty" json:"disable_pagination,omitempty"`
	Outputs           []string `yaml:"outputs,omitempty" json:"outputs,omitempty"`
	TplType           string   `yaml:"tpl_type,omitempty" json:"tpl_type,omitempty"` // crud/tree
}

func (a *S) Format() *S {
	for i, item := range a.Fields {
		a.Fields[i] = item.Format()
	}
	return a
}

type Field struct {
	Name       string        `yaml:"name,omitempty" json:"name,omitempty"`
	ColumnName string        `yaml:"column_name,omitempty" json:"column_name,omitempty"`
	Type       string        `yaml:"type,omitempty" json:"type,omitempty"`
	Comment    string        `yaml:"comment,omitempty" json:"comment,omitempty"`   // {{with .Comment}}// {{.}}{{end}}
	GormTag    *FieldGormTag `yaml:"gorm_tag,omitempty" json:"gorm_tag,omitempty"` // {{with .GormTag}}gorm:"{{if .Tag}}{{.Tag}}{{else}}{{if .Index}};index{{end}}{{if .Size}};size:{{.Size}}{{end}}{{end}},omitempty" {{end}}
	JSONTag    FieldJSONTag  `yaml:"json_tag,omitempty" json:"json_tag,omitempty"` // {{with .JSONTag}}json:"{{if .Tag}}{{.Tag}}{{else}}{{lowerUnderline $.Name}}{{if .OmitEmpty}},omitempty{{end}}{{end}}"{{end}}
	Query      *FieldQuery   `yaml:"query,omitempty" json:"query,omitempty"`
	Order      *FieldOrder   `yaml:"order,omitempty" json:"order,omitempty"`
	Form       *FieldForm    `yaml:"form,omitempty" json:"form,omitempty"`
}

func (a *Field) Format() *Field {
	if a.Query != nil {
		if a.Query.Name == "" {
			a.Query.Name = a.Name
		}
		if a.Query.Comment == "" {
			a.Query.Comment = a.Comment
		}
		if a.Query.InQuery && a.Query.FormTag == "" {
			a.Query.FormTag = utils.ToLowerCamel(a.Name)
		}
		if a.Query.OP == "" {
			a.Query.OP = "="
		}
	}
	if a.Form != nil {
		if a.Form.Comment == "" {
			a.Form.Comment = a.Comment
		}
	}
	return a
}

type FieldGormTag struct {
	Tag   string `yaml:"tag,omitempty" json:"tag,omitempty"`
	Size  int    `yaml:"size,omitempty" json:"size,omitempty"`
	Index bool   `yaml:"index,omitempty" json:"index,omitempty"`
	Type  string `yaml:"type,omitempty" json:"type,omitempty"`
}

type FieldJSONTag struct {
	Tag       string `yaml:"tag,omitempty" json:"tag,omitempty"`
	OmitEmpty bool   `yaml:"omit_empty,omitempty" json:"omit_empty,omitempty"`
}

type FieldBindingTag struct {
	Tag string `yaml:"tag,omitempty" json:"tag,omitempty"`
}

type FieldQuery struct {
	Name       string           `yaml:"name,omitempty" json:"name,omitempty"`
	InQuery    bool             `yaml:"in_query,omitempty" json:"in_query,omitempty"`
	FormTag    string           `yaml:"form_tag,omitempty" json:"form_tag,omitempty"`
	BindingTag *FieldBindingTag `yaml:"binding_tag,omitempty" json:"binding_tag,omitempty"` // {{with .BindingTag}}binding:"{{if .Tag}}{{.Tag}}{{else}}{{if .Required}}required{{end}}{{end}}"{{end}}
	Comment    string           `yaml:"comment,omitempty" json:"comment,omitempty"`         // {{with .Comment}}// {{.}}{{end}}
	IfCond     string           `yaml:"cond,omitempty" json:"cond,omitempty"`               // {{with .IfCond}}{{.}}{{end}}
	OP         string           `yaml:"op,omitempty" json:"op,omitempty"`                   // LIKE/=/</>/<=/>=/<>
}

type FieldOrder struct {
	Direction string `yaml:"direction,omitempty" json:"direction,omitempty"`
}

type FieldForm struct {
	JSONTag    FieldJSONTag     `yaml:"json_tag,omitempty" json:"json_tag,omitempty"`
	BindingTag *FieldBindingTag `yaml:"binding_tag,omitempty" json:"binding_tag,omitempty"`
	Comment    string           `yaml:"comment,omitempty" json:"comment,omitempty"`
}
