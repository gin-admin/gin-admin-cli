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
{{$includeID := .Include.ID}}
{{$includeCreatedAt := .Include.CreatedAt}}
{{$includeUpdatedAt := .Include.UpdatedAt}}
{{$includeStatus := .Include.Status}}

{{with .Comment}}// {{.}}{{else}}// Defining the `{{$name}}` business logic.{{end}}
type {{$name}} struct {
	Trans       *utils.Trans
	{{$name}}DAL *dal.{{$name}}
}

// Query {{lowerSpacePlural .Name}} from the data access object based on the provided parameters and options.
func (a *{{$name}}) Query(ctx context.Context, params schema.{{$name}}QueryParam) (*schema.{{$name}}QueryResult, error) {
	params.Pagination = {{if .DisablePagination}}false{{else}}true{{end}}

	result, err := a.{{$name}}DAL.Query(ctx, params, schema.{{$name}}QueryOptions{
		QueryOptions: utils.QueryOptions{
			OrderFields: []utils.OrderByParam{
                {{- range .Fields}}{{$fieldName := .Name}}
				{{- if .Order}}
				{Field: "{{lowerUnderline $fieldName}}", Direction: {{if eq .Order "DESC"}}utils.DESC{{else}}utils.ASC{{end}}},
				{{- end}}
                {{- end}}
			},
		},
	})
	if err != nil {
		return nil, err
	}
	result.Data = result.Data.ToTree()
	return result, nil
}

func (a *{{$name}}) appendChildren(ctx context.Context, data schema.{{plural .Name}}) (schema.{{plural .Name}}, error) {
	if len(data) == 0 {
		return data, nil
	}

	existsInData := func(id string) bool {
		for _, item := range data {
			if item.ID == id {
				return true
			}
		}
		return false
	}

	for _, item := range data {
		childResult, err := a.{{$name}}DAL.Query(ctx, schema.{{$name}}QueryParam{
			ParentPathPrefix: item.ParentPath + item.ID + utils.TreePathDelimiter,
		})
		if err != nil {
			return nil, err
		}
		for _, child := range childResult.Data {
			if existsInData(child.ID) {
				continue
			}
			data = append(data, child)
		}
	}

	parentIDs := data.SplitParentIDs()
	if len(parentIDs) > 0 {
		parentResult, err := a.{{$name}}DAL.Query(ctx, schema.{{$name}}QueryParam{
			InIDs: parentIDs,
		})
		if err != nil {
			return nil, err
		}
		for _, p := range parentResult.Data {
			if existsInData(p.ID) {
				continue
			}
			data = append(data, p)
		}
		sort.Sort(data)
	}

	return data, nil
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
func (a *{{$name}}) Create(ctx context.Context, formItem *schema.{{$name}}Form) (*schema.{{$name}}, error) {
	{{lowerCamel $name}} := &schema.{{$name}}{
		{{if $includeID}}ID:          idx.NewXID(),{{end}}
		{{if $includeCreatedAt}}CreatedAt:   time.Now(),{{end}}
	}

	if parentID := formItem.ParentID; parentID != "" {
		parent, err := a.{{$name}}DAL.Get(ctx, parentID)
		if err != nil {
			return nil, err
		} else if parent == nil {
			return nil, errors.NotFound("", "Parent not found")
		}
		{{lowerCamel $name}}.ParentPath = parent.ParentPath + parent.ID + utils.TreePathDelimiter
	}
	if err := formItem.FillTo({{lowerCamel $name}}); err != nil {
		return nil, err
	}

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
func (a *{{$name}}) Update(ctx context.Context, id string, formItem *schema.{{$name}}Form) error {
	{{lowerCamel $name}}, err := a.{{$name}}DAL.Get(ctx, id)
	if err != nil {
		return err
	} else if {{lowerCamel $name}} == nil {
		return errors.NotFound("", "{{titleSpace $name}} not found")
	}

	oldParentPath := {{lowerCamel $name}}.ParentPath
	{{- if $includeStatus}}
	oldStatus := {{lowerCamel $name}}.Status
	{{- end}}
	var childData schema.{{plural .Name}}
	if {{lowerCamel $name}}.ParentID != formItem.ParentID {
		if parentID := formItem.ParentID; parentID != "" {
			parent, err := a.{{$name}}DAL.Get(ctx, parentID)
			if err != nil {
				return err
			} else if parent == nil {
				return errors.NotFound("", "Parent not found")
			}
			{{lowerCamel $name}}.ParentPath = parent.ParentPath + parent.ID + utils.TreePathDelimiter
		} else {
			{{lowerCamel $name}}.ParentPath = ""
		}

		childResult, err := a.{{$name}}DAL.Query(ctx, schema.{{$name}}QueryParam{
			ParentPathPrefix: oldParentPath + {{lowerCamel $name}}.ID + utils.TreePathDelimiter,
		}, schema.{{$name}}QueryOptions{
			QueryOptions: utils.QueryOptions{
				SelectFields: []string{"id", "parent_path"},
			},
		})
		if err != nil {
			return err
		}
		childData = childResult.Data
	}
	if err := formItem.FillTo({{lowerCamel $name}}); err != nil {
		return err
	}
	{{if $includeUpdatedAt}}{{lowerCamel $name}}.UpdatedAt = time.Now(){{end}}
	
	return a.Trans.Exec(ctx, func(ctx context.Context) error {
		if err := a.{{$name}}DAL.Update(ctx, {{lowerCamel $name}}); err != nil {
			return err
		}

		{{- if $includeStatus}}
		if oldStatus != formItem.Status {
			opath := oldParentPath + {{lowerCamel $name}}.ID + utils.TreePathDelimiter
			if err := a.{{$name}}DAL.UpdateStatusByParentPath(ctx, opath, formItem.Status); err != nil {
				return err
			}
		}
		{{- end}}

		for _, child := range childData {
			opath := oldParentPath + {{lowerCamel $name}}.ID + utils.TreePathDelimiter
			npath := {{lowerCamel $name}}.ParentPath + {{lowerCamel $name}}.ID + utils.TreePathDelimiter
			err := a.{{$name}}DAL.UpdateParentPath(ctx, child.ID, strings.Replace(child.ParentPath, opath, npath, 1))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// Delete the specified {{lowerSpace .Name}} from the data access object.
func (a *{{$name}}) Delete(ctx context.Context, id string) error {
	{{lowerCamel $name}}, err := a.{{$name}}DAL.Get(ctx, id)
	if err != nil {
		return err
	} else if {{lowerCamel $name}} == nil {
		return errors.NotFound("", "{{titleSpace $name}} not found")
	}

	childResult, err := a.{{$name}}DAL.Query(ctx, schema.{{$name}}QueryParam{
		ParentPathPrefix: {{lowerCamel $name}}.ParentPath + {{lowerCamel $name}}.ID + utils.TreePathDelimiter,
		}, schema.{{$name}}QueryOptions{
		QueryOptions: utils.QueryOptions{
			SelectFields: []string{"id"},
		},
	})
	if err != nil {
		return err
	}

	return a.Trans.Exec(ctx, func(ctx context.Context) error {
		if err := a.{{$name}}DAL.Delete(ctx, id); err != nil {
			return err
		}

		for _, child := range childResult.Data {
			if err := a.{{$name}}DAL.Delete(ctx, child.ID); err != nil {
				return err
			}
		}

		return nil
	})
}
