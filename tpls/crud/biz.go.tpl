package biz

import (
	"context"
	"time"

	"{{.UtilsImportPath}}"
	"{{.ModuleImportPath}}/dal"
	"{{.ModuleImportPath}}/schema"
	"{{.RootImportPath}}/pkg/errors"
	"{{.RootImportPath}}/pkg/idx"
)

{{$name := .Name}}

{{with .Comment}}// {{.}}{{else}}// {{$name}} business logic{{end}}
type {{$name}} struct {
	Trans       *utils.Trans
	{{$name}}DAL *dal.{{$name}}
}

// Query {{lowerPlural .Name}} from the data access object based on the provided parameters and options.
func (a *{{$name}}) Query(ctx context.Context, params schema.{{$name}}QueryParam) (*schema.{{$name}}QueryResult, error) {
	params.Pagination = {{if .DisablePagination}}false{{else}}true{{end}}

	result, err := a.{{$name}}DAL.Query(ctx, params, schema.{{$name}}QueryOptions{
		QueryOptions: utils.QueryOptions{
			OrderFields: []utils.OrderByParam{
                {{- range .Fields}}{{$fieldName := .Name}}
				{{- with .Order}}
				{Field: "{{lowerUnderline $fieldName}}", Direction: {{if eq .Direction "DESC"}}utils.DESC{{else}}utils.ASC{{end}}},
				{{- end}}
                {{- end}}
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Get the specified {{lowerSpace .Name}} from the data access object.
func (a *{{$name}}) Get(ctx context.Context, id string) (*schema.{{$name}}, error) {
	{{lowerCamel $name}}, err := a.{{$name}}DAL.Get(ctx, id)
	if err != nil {
		return nil, err
	} else if {{lowerCamel $name}} == nil {
		return nil, errors.NotFound("", "{{titleSpace $name}} not found")
	}
	return {{lowerCamel $name}}, nil
}

// Create a new {{lowerSpace .Name}} in the data access object.
func (a *{{$name}}) Create(ctx context.Context, sitem *schema.{{$name}}Save) (*schema.{{$name}}, error) {
	{{lowerCamel $name}} := &schema.{{$name}}{
		ID:          idx.NewXID(),
		CreatedAt:   time.Now(),
	}

    {{- range .Fields}}{{$fieldName := .Name}}
	{{- with .Form}}
	{{lowerCamel $name}}.{{$fieldName}} = sitem.{{$fieldName}}
	{{- end}}
    {{- end}}

	err := a.Trans.Exec(ctx, func(ctx context.Context) error {
		if err := a.{{$name}}DAL.Create(ctx, {{lowerCamel $name}}); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return {{lowerCamel $name}}, nil
}

// Update the specified {{lowerSpace .Name}} in the data access object.
func (a *{{$name}}) Update(ctx context.Context, id string, sitem *schema.{{$name}}Save) error {
	{{lowerCamel $name}}, err := a.{{$name}}DAL.Get(ctx, id)
	if err != nil {
		return err
	} else if {{lowerCamel $name}} == nil {
		return errors.NotFound("", "{{titleSpace $name}} not found")
	}
    {{- range .Fields}}{{$fieldName := .Name}}
	{{- with .Form}}
	{{lowerCamel $name}}.{{$fieldName}} = sitem.{{$fieldName}}
	{{- end}}
    {{- end}}
    {{lowerCamel $name}}.UpdatedAt = time.Now()
	
	return a.Trans.Exec(ctx, func(ctx context.Context) error {
		if err := a.{{$name}}DAL.Update(ctx, {{lowerCamel $name}}); err != nil {
			return err
		}
		return nil
	})
}

// Delete the specified {{lowerSpace .Name}} from the data access object.
func (a *{{$name}}) Delete(ctx context.Context, id string) error {
	exists, err := a.{{$name}}DAL.Exists(ctx, id)
	if err != nil {
		return err
	} else if !exists {
		return errors.NotFound("", "{{titleSpace $name}} not found")
	}

	return a.Trans.Exec(ctx, func(ctx context.Context) error {
		if err := a.{{$name}}DAL.Delete(ctx, id); err != nil {
			return err
		}
		return nil
	})
}
