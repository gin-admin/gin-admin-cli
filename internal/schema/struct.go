package schema

type S struct {
	RootImportPath    string   `yaml:"root_import_path" json:"root_import_path"`
	ModuleName        string   `yaml:"module_name" json:"module_name"`
	ModuleImportPath  string   `yaml:"module_import_path" json:"module_import_path"`
	Name              string   `yaml:"name" json:"name"`
	TableName         string   `yaml:"table_name" json:"table_name"`
	Comment           string   `yaml:"comment" json:"comment"`
	Fields            []*Field `yaml:"fields" json:"fields"`
	DisablePagination bool     `yaml:"disable_pagination" json:"disable_pagination"`
}

type Field struct {
	Name       string       `yaml:"name" json:"name"`
	ColumnName string       `yaml:"column_name" json:"column_name"`
	Type       string       `yaml:"type" json:"type"`
	GormTag    FieldGormTag `yaml:"gorm_tag" json:"gorm_tag"` // {{with .GormTag}}gorm:"{{if .Tag}}{{.Tag}}{{else}}{{if .Index}};index{{end}}{{if .Size}};size:{{.Size}}{{end}}{{end}}" {{end}}
	JSONTag    FieldJSONTag `yaml:"json_tag" json:"json_tag"` // {{with .JSONTag}}json:"{{if .Tag}}{{.Tag}}{{else}}{{lowerUnderline $.Name}}{{if .OmitEmpty}},omitempty{{end}}{{end}}"{{end}}
	Comment    string       `yaml:"comment" json:"comment"`   // {{with .Comment}}// {{.}}{{end}}
	Query      *FieldQuery  `yaml:"query" json:"query"`
	Order      *FieldOrder  `yaml:"order" json:"order"`
	Form       *FieldForm   `yaml:"form" json:"form"`
}

type FieldGormTag struct {
	Tag   string `yaml:"tag" json:"tag"`
	Size  int    `yaml:"size" json:"size"`
	Index bool   `yaml:"index" json:"index"`
}

type FieldJSONTag struct {
	Tag       string `yaml:"tag" json:"tag"`
	OmitEmpty bool   `yaml:"omit_empty" json:"omit_empty"`
}

type FieldBindingTag struct {
	Tag      string `yaml:"tag" json:"tag"`
	Required bool   `yaml:"required" json:"required"`
}

type FieldQuery struct {
	Name       string           `yaml:"name" json:"name"`
	InQuery    bool             `yaml:"in_query" json:"in_query"`
	FormTag    string           `yaml:"form_tag" json:"form_tag"`
	BindingTag *FieldBindingTag `yaml:"binding_tag" json:"binding_tag"` // {{with .BindingTag}}binding:"{{if .Tag}}{{.Tag}}{{else}}{{if .Required}}required{{end}}{{end}}"{{end}}
	Comment    string           `yaml:"comment" json:"comment"`         // {{with .Comment}}// {{.}}{{end}}
	IfCond     string           `yaml:"cond" json:"cond"`               // {{with .IfCond}}{{.}}{{end}}
	OP         string           `yaml:"op" json:"op"`                   // LIKE/=/</>/<=/>=/<>
}

type FieldOrder struct {
	Direction string `yaml:"direction" json:"direction"`
}

type FieldForm struct {
	JSONTag    FieldJSONTag     `yaml:"json_tag" json:"json_tag"`
	BindingTag *FieldBindingTag `yaml:"binding_tag" json:"binding_tag"`
	Comment    string           `yaml:"comment" json:"comment"`
}
