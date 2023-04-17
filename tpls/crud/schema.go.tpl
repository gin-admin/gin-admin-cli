package schema

import (
	"time"

	"{{.UtilsImportPath}}"
)

{{$name := .Name}}

{{with .Comment}}// {{.}}{{else}}// Defining the `{{$name}}` struct.{{end}}
type {{$name}} struct {
    {{- range .Fields}}{{$fieldName := .Name}}
	{{$fieldName}} {{.Type}} `json:"{{.JSONTag}}"{{with .GormTag}} gorm:"{{.}}"{{end}}{{with .CustomTag}} {{raw .}}{{end}}`{{with .Comment}}// {{.}}{{end}}
	{{- end}}
}

{{- with .TableName}}
// Defining the name of the database table that corresponds to the `{{$name}}` struct.
func (a {{title $name}}) TableName() string {
	return "{{.}}"
}
{{- end}}

// Defining the query parameters for the `{{$name}}` struct.
type {{$name}}QueryParam struct {
	utils.PaginationParam
	{{- range .Fields}}{{$fieldName := .Name}}{{$type :=.Type}}
	{{- with .Query}}
	{{.Name}} {{$type}} `form:"{{with .FormTag}}{{.}}{{else}}-{{end}}"{{with .BindingTag}} binding:"{{.}}"{{end}}{{with .CustomTag}} {{raw .}}{{end}}`{{with .Comment}}// {{.}}{{end}}
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
type {{$name}}Form struct {
	{{- range .Fields}}{{$fieldName := .Name}}{{$type :=.Type}}
	{{- with .Form}}
	{{.Name}} {{$type}} `json:"{{.JSONTag}}"{{with .BindingTag}} binding:"{{.}}"{{end}}{{with .CustomTag}} {{raw .}}{{end}}`{{with .Comment}}// {{.}}{{end}}
	{{- end}}
	{{- end}}
}

// A validation function for the `{{$name}}Form` struct.
func (a *{{$name}}Form) Validate() error {
	return nil
}

func (a *{{$name}}Form) FillTo({{lowerCamel $name}} *{{$name}}) *{{$name}} {
	{{- range .Fields}}{{$fieldName := .Name}}
	{{- with .Form}}
	{{lowerCamel $name}}.{{$fieldName}} = a.{{.Name}}
	{{- end}}
    {{- end}}
	return {{lowerCamel $name}}
}
