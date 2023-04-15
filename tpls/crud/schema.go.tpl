package schema

import (
	"time"

	"{{.RootImportPath}}/internal/utils"
)

{{$name := .Name}}

{{with .Comment}}// {{.}}{{else}}// Defining the `{{$name}}` struct.{{end}}
type {{$name}} struct {
	ID          string    `gorm:"size:20;primarykey;" json:"id"` // Unique ID
    {{- range .Fields}}{{$fieldName := .Name}}{{$columnName := .ColumnName}}
	{{$fieldName}} {{.Type}} `{{with .GormTag}}gorm:"{{with .Tag}}{{.}}{{else}}{{with $columnName}}column:{{.}};{{end}}{{with .Index}}index;{{end}}{{with .Size}}size:{{.}};{{end}}{{with .Type}}type:{{.}};{{end}}{{end}}"{{end}} {{- with .JSONTag}} json:"{{if .Tag}}{{.Tag}}{{else}}{{lowerUnderline $fieldName}}{{if .OmitEmpty}},omitempty{{end}}{{end}}"{{end}}`{{with .Comment}}// {{.}}{{end}}
	{{- end}}
	CreatedAt   time.Time `gorm:"index;" json:"created_at"`      // Create time
	UpdatedAt   time.Time `gorm:"index;" json:"updated_at"`      // Update time
}

{{with .TableName}}
// Defining the name of the database table that corresponds to the `{{$name}}` struct.
func (a {{title $name}}) TableName() string {
	return "{{.}}"
}
{{end}}

// Defining the query parameters for the `{{$name}}` struct.
type {{$name}}QueryParam struct {
	utils.PaginationParam
	{{- range .Fields}}{{$fieldName := .Name}}{{$type :=.Type}}
	{{- with .Query}}{{$inQuery := .InQuery}}{{$queryName := .Name}}
	{{.Name}} {{$type}} `form:"{{with .FormTag}}{{.}}{{else}}{{if $inQuery}}{{lowerCamel $queryName}}{{else}}-{{end}}{{end}}"{{with .BindingTag}} binding:"{{if .Tag}}{{.Tag}}{{else}}required{{end}}"{{end}}`{{with .Comment}}// {{.}}{{end -}}
	{{- end}}
	{{- end}}
}

// Defining the query options for the `{{$name}}` struct.
type {{$name}}QueryOptions struct {
	utils.QueryOptions
}

// Defining the query result for the `{{$name}}` struct.
type {{$name}}QueryResult struct {
	Data       {{plural .Name}}
	PageResult *utils.PaginationResult
}

// Defining the slice of `{{$name}}` struct.
type {{plural .Name}} []*{{$name}}

// Defining the data structure for creating a `{{$name}}` struct.
type {{$name}}Save struct {
	{{- range .Fields}}{{$fieldName := .Name}}{{$type :=.Type}}
	{{- with .Form}}
	{{$fieldName}} {{$type}} `{{with .JSONTag -}} json:"{{if .Tag}}{{.Tag}}{{else}}{{lowerUnderline $fieldName}}{{if .OmitEmpty}},omitempty{{end}}{{end}}"{{end}} {{- with .BindingTag}} binding:"{{if .Tag}}{{.Tag}}{{else}}required{{end}}"{{end}}`{{with .Comment}}// {{.}}{{end}}
	{{- end}}
	{{- end}}
}

// A validation function for the `{{$name}}Save` struct.
func (a *{{$name}}Save) Validate() error {
	return nil
}
