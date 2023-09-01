package dal

import (
	"context"

	"{{.UtilImportPath}}"
	"{{.ModuleImportPath}}/schema"
	"{{.RootImportPath}}/pkg/errors"
	"gorm.io/gorm"
)

{{$name := .Name}}
{{$includeCreatedAt := .Include.CreatedAt}}
{{$includeStatus := .Include.Status}}
{{$treeTpl := eq .TplType "tree"}}

// Get {{lowerSpace .Name}} storage instance
func Get{{$name}}DB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return util.GetDB(ctx, defDB).Model(new(schema.{{$name}}))
}

// Get {{lowerSpace .Name}} table name
func Get{{$name}}TableName(defDB *gorm.DB) string {
	stat := gorm.Statement{DB: defDB}
	stat.Parse(&schema.{{$name}}{})
	return stat.Table
}

{{with .Comment}}// {{.}}{{else}}// Defining the `{{$name}}` data access object.{{end}}
type {{$name}} struct {
	DB *gorm.DB
}

// Query {{lowerSpacePlural .Name}} from the database based on the provided parameters and options.
func (a *{{$name}}) Query(ctx context.Context, params schema.{{$name}}QueryParam, opts ...schema.{{$name}}QueryOptions) (*schema.{{$name}}QueryResult, error) {
	var opt schema.{{$name}}QueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	db := Get{{$name}}DB(ctx, a.DB)

	{{- if $treeTpl}}
	if v:= params.InIDs; len(v) > 0 {
		db = db.Where("id IN ?", v)
	}
	{{- end}}

    {{- range .Fields}}{{$type := .Type}}{{$fieldName := .Name}}
    {{- with .Query}}
	if v := params.{{.Name}}; {{with .IfCond}}{{.}}{{else}}{{convIfCond $type}}{{end}} {
		db = db.Where("{{lowerUnderline $fieldName}} {{.OP}} ?", {{if .Args}}{{raw .Args}}{{else}}{{if eq .OP "LIKE"}}"%"+v+"%"{{else}}v{{end}}{{end}})
	}
    {{- end}}
    {{- end}}

	var list schema.{{plural .Name}}
	pageResult, err := util.WrapPageQuery(ctx, db, params.PaginationParam, opt.QueryOptions, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	queryResult := &schema.{{$name}}QueryResult{
		PageResult: pageResult,
		Data:       list,
	}
	return queryResult, nil
}

// Get the specified {{lowerSpace .Name}} from the database.
func (a *{{$name}}) Get(ctx context.Context, id string, opts ...schema.{{$name}}QueryOptions) (*schema.{{$name}}, error) {
	var opt schema.{{$name}}QueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	item := new(schema.{{$name}})
	ok, err := util.FindOne(ctx, Get{{$name}}DB(ctx, a.DB).Where("id=?", id), opt.QueryOptions, item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}
	return item, nil
}

// Exist checks if the specified {{lowerSpace .Name}} exists in the database.
func (a *{{$name}}) Exists(ctx context.Context, id string) (bool, error) {
	ok, err := util.Exists(ctx, Get{{$name}}DB(ctx, a.DB).Where("id=?", id))
	return ok, errors.WithStack(err)
}

// Create a new {{lowerSpace .Name}}.
func (a *{{$name}}) Create(ctx context.Context, item *schema.{{$name}}) error {
	result := Get{{$name}}DB(ctx, a.DB).Create(item)
	return errors.WithStack(result.Error)
}

// Update the specified {{lowerSpace .Name}} in the database.
func (a *{{$name}}) Update(ctx context.Context, item *schema.{{$name}}) error {
	result := Get{{$name}}DB(ctx, a.DB).Where("id=?", item.ID).Select("*"){{if $includeCreatedAt}}.Omit("created_at"){{end}}.Updates(item)
	return errors.WithStack(result.Error)
}

// Delete the specified {{lowerSpace .Name}} from the database.
func (a *{{$name}}) Delete(ctx context.Context, id string) error {
	result := Get{{$name}}DB(ctx, a.DB).Where("id=?", id).Delete(new(schema.{{$name}}))
	return errors.WithStack(result.Error)
}

{{- if $treeTpl}}
// Updates the parent path of the specified {{lowerSpace .Name}}.
func (a *{{$name}}) UpdateParentPath(ctx context.Context, id, parentPath string) error {
	result := Get{{$name}}DB(ctx, a.DB).Where("id=?", id).Update("parent_path", parentPath)
	return errors.WithStack(result.Error)
}

{{- if $includeStatus}}
// Updates the status of all {{lowerPlural .Name}} whose parent path starts with the provided parent path.
func (a *{{$name}}) UpdateStatusByParentPath(ctx context.Context, parentPath, status string) error {
	result := Get{{$name}}DB(ctx, a.DB).Where("parent_path like ?", parentPath+"%").Update("status", status)
	return errors.WithStack(result.Error)
}
{{- end}}
{{- end}}