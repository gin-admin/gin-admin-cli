package schema

import (
	"time"

	{{if .TableName}}"{{.RootImportPath}}/internal/config"{{end}}
	"{{.UtilImportPath}}"
)

{{$name := .Name}}
{{$includeSequence := .Include.Sequence}}
{{$treeTpl := eq .TplType "tree"}}

{{with .Comment}}// {{.}}{{else}}// Defining the `{{$name}}` struct.{{end}}
type {{$name}} struct {
    {{- range .Fields}}{{$fieldName := .Name}}
	{{$fieldName}} {{.Type}} `json:"{{.JSONTag}}"{{with .GormTag}} gorm:"{{.}}"{{end}}{{with .CustomTag}} {{raw .}}{{end}}`{{with .Comment}}// {{.}}{{end}}
	{{- end}}
}

{{- if .TableName}}
func (a {{$name}}) TableName() string {
	return config.C.FormatTableName("{{.TableName}}")
}
{{- end}}

// Defining the query parameters for the `{{$name}}` struct.
type {{$name}}QueryParam struct {
	util.PaginationParam
	{{if $treeTpl}}InIDs []string `form:"-"`{{- end}}
	{{- range .Fields}}{{$fieldName := .Name}}{{$type :=.Type}}
	{{- with .Query}}
	{{.Name}} {{$type}} `form:"{{with .FormTag}}{{.}}{{else}}-{{end}}"{{with .BindingTag}} binding:"{{.}}"{{end}}{{with .CustomTag}} {{raw .}}{{end}}`{{with .Comment}}// {{.}}{{end}}
	{{- end}}
	{{- range .Queries}}
	{{- with .}}
	{{.Name}} {{$type}} `form:"{{with .FormTag}}{{.}}{{else}}-{{end}}"{{with .BindingTag}} binding:"{{.}}"{{end}}{{with .CustomTag}} {{raw .}}{{end}}`{{with .Comment}}// {{.}}{{end}}
	{{- end}}
	{{- end}}
	{{- end}}
}

// Defining the query options for the `{{$name}}` struct.
type {{$name}}QueryOptions struct {
	util.QueryOptions
}

// Defining the query result for the `{{$name}}` struct.
type {{$name}}QueryResult struct {
	Data       {{.Name}}List
	PageResult *util.PaginationResult
}

// Defining the slice of `{{$name}}` struct.
type {{.Name}}List []*{{$name}}

{{- if $includeSequence}}
func (a {{.Name}}List) Len() int {
	return len(a)
}

func (a {{.Name}}List) Less(i, j int) bool {
	if a[i].Sequence == a[j].Sequence {
		return a[i].CreatedAt.Unix() > a[j].CreatedAt.Unix()
	}
	return a[i].Sequence > a[j].Sequence
}

func (a {{.Name}}List) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
{{- end}}

{{- if $treeTpl}}
func (a {{.Name}}List) ToMap() map[string]*{{$name}} {
	m := make(map[string]*{{$name}})
	for _, item := range a {
		m[item.ID] = item
	}
	return m
}

func (a {{.Name}}List) SplitParentIDs() []string {
	parentIDs := make([]string, 0, len(a))
	idMapper := make(map[string]struct{})
	for _, item := range a {
		if _, ok := idMapper[item.ID]; ok {
			continue
		}
		idMapper[item.ID] = struct{}{}
		if pp := item.ParentPath; pp != "" {
			for _, pid := range strings.Split(pp, util.TreePathDelimiter) {
				if pid == "" {
					continue
				}
				if _, ok := idMapper[pid]; ok {
					continue
				}
				parentIDs = append(parentIDs, pid)
				idMapper[pid] = struct{}{}
			}
		}
	}
	return parentIDs
}

func (a {{.Name}}List) ToTree() {{.Name}}List {
	var list {{.Name}}List
	m := a.ToMap()
	for _, item := range a {
		if item.ParentID == "" {
			list = append(list, item)
			continue
		}
		if parent, ok := m[item.ParentID]; ok {
			if parent.Children == nil {
				children := {{.Name}}List{item}
				parent.Children = &children
				continue
			}
			*parent.Children = append(*parent.Children, item)
		}
	}
	return list
}
{{- end}}

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

// Convert `{{$name}}Form` to `{{$name}}` object.
func (a *{{$name}}Form) FillTo({{lowerCamel $name}} *{{$name}}) error {
	{{- range .Fields}}{{$fieldName := .Name}}
	{{- with .Form}}
	{{lowerCamel $name}}.{{$fieldName}} = a.{{.Name}}
	{{- end}}
    {{- end}}
	return nil
}
