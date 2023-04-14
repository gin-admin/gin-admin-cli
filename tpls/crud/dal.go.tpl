package dal

import (
	"context"

	"{{.RootImportPath}}/internal/library/utils"
	"{{.ModuleImportPath}}/schema"
	"{{.RootImportPath}}/pkg/errors"
	"gorm.io/gorm"
)

{{$name := .Name}}

// Get {{lowerSpace .Name}} storage instance
func Get{{$name}}DB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return utils.GetDB(ctx, defDB).Model(new(schema.{{$name}}))
}

{{with .Comment}}// {{.}}{{else}}// {{$name}} data access object{{end}}
type {{$name}} struct {
	DB *gorm.DB
}

// Query {{lowerPlural .Name}} from the database based on the provided parameters and options.
func (a *{{$name}}) Query(ctx context.Context, params schema.{{$name}}QueryParam, opts ...schema.{{$name}}QueryOptions) (*schema.{{$name}}QueryResult, error) {
	var opt schema.{{$name}}QueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	db := Get{{$name}}DB(ctx, a.DB)

    {{range .Fields -}}
        {{$type := .Type}}
        {{with .Query -}}
            if v := params.{{.Name}}; {{with .IfCond}}{{.}}{{else}}{{convIfCond $type}}{{end}} {
		        db = db.Where("{{lowerUnderline .Name}} {{.OP}} ?", {{if eq .OP "LIKE"}}"%"+v+"%"{{else}}v{{end}})
	        }
        {{- end -}}
    {{- end -}}

	var list schema.{{plural .Name}}
	pageResult, err := utils.WrapPageQuery(ctx, db, params.PaginationParam, opt.QueryOptions, &list)
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
	ok, err := utils.FindOne(ctx, Get{{$name}}DB(ctx, a.DB).Where("id=?", id), opt.QueryOptions, item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}
	return item, nil
}

// Exist checks if the specified {{lowerSpace .Name}} exists in the database.
func (a *{{$name}}) Exists(ctx context.Context, id string) (bool, error) {
	ok, err := utils.Exists(ctx, Get{{$name}}DB(ctx, a.DB).Where("id=?", id))
	return ok, errors.WithStack(err)
}

// Create a new {{lowerSpace .Name}}.
func (a *{{$name}}) Create(ctx context.Context, item *schema.{{$name}}) error {
	result := Get{{$name}}DB(ctx, a.DB).Create(item)
	return errors.WithStack(result.Error)
}

// Update the specified {{lowerSpace .Name}} in the database.
func (a *{{$name}}) Update(ctx context.Context, item *schema.{{$name}}) error {
	result := Get{{$name}}DB(ctx, a.DB).Where("id=?", item.ID).Select("*").Omit("created_at").Updates(item)
	return errors.WithStack(result.Error)
}

// Delete the specified {{lowerSpace .Name}} from the database.
func (a *{{$name}}) Delete(ctx context.Context, id string) error {
	result := Get{{$name}}DB(ctx, a.DB).Where("id=?", id).Delete(new(schema.{{$name}}))
	return errors.WithStack(result.Error)
}
