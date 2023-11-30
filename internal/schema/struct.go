package schema

import (
	"fmt"
	"strings"

	"github.com/gin-admin/gin-admin-cli/v10/internal/utils"
)

type S struct {
	RootImportPath   string `yaml:"-" json:"-"`
	ModuleName       string `yaml:"-" json:"-"`
	ModuleImportPath string `yaml:"-" json:"-"`
	UtilImportPath   string `yaml:"-" json:"-"`
	Include          struct {
		ID        bool
		Status    bool
		CreatedAt bool
		UpdatedAt bool
		Sequence  bool
	} `yaml:"-" json:"-"`
	Name                 string            `yaml:"name,omitempty" json:"name,omitempty"`
	TableName            string            `yaml:"table_name,omitempty" json:"table_name,omitempty"`
	Comment              string            `yaml:"comment,omitempty" json:"comment,omitempty"`
	Outputs              []string          `yaml:"outputs,omitempty" json:"outputs,omitempty"`
	TplType              string            `yaml:"tpl_type,omitempty" json:"tpl_type,omitempty"` // crud/tree
	DisablePagination    bool              `yaml:"disable_pagination,omitempty" json:"disable_pagination,omitempty"`
	DisableDefaultFields bool              `yaml:"disable_default_fields,omitempty" json:"disable_default_fields,omitempty"`
	FillGormCommit       bool              `yaml:"fill_gorm_commit,omitempty" json:"fill_gorm_commit,omitempty"`
	Fields               []*Field          `yaml:"fields,omitempty" json:"fields,omitempty"`
	GenerateFE           bool              `yaml:"generate_fe,omitempty" json:"generate_fe,omitempty"`
	FEMapping            map[string]string `yaml:"fe_mapping,omitempty" json:"fe_mapping,omitempty"` // tpl -> file
	Extra                map[string]string `yaml:"extra,omitempty" json:"extra,omitempty"`
}

func (a *S) Format() *S {
	if a.TableName == "" {
		a.TableName = utils.ToLowerUnderlinedNamer(a.Name)
	}
	if a.TplType != "" {
		a.TplType = strings.ToLower(a.TplType)
	}

	if !a.DisableDefaultFields {
		var fields []*Field
		fields = append(fields, &Field{
			Name:    "ID",
			Type:    "string",
			GormTag: "size:20;primarykey;",
			Comment: "Unique ID",
		})
		fields = append(fields, a.Fields...)

		if a.TplType == "tree" {
			fields = append(fields, &Field{
				Name:    "ParentID",
				Type:    "string",
				GormTag: "size:20;index;",
				Comment: fmt.Sprintf("Parent ID (From %s.ID)", a.Name),
				Query:   &FieldQuery{},
				Form:    &FieldForm{},
			})
			fields = append(fields, &Field{
				Name:    "ParentPath",
				Type:    "string",
				GormTag: "size:255;index;",
				Comment: "Parent path (split by .)",
				Query: &FieldQuery{
					Name: "ParentPathPrefix",
					OP:   "LIKE",
					Args: `v + "%"`,
				},
			})
			fields = append(fields, &Field{
				Name:    "Children",
				Type:    fmt.Sprintf("*%s", utils.ToPlural(a.Name)),
				GormTag: "-",
				Comment: "Child nodes",
			})
		}

		fields = append(fields, &Field{
			Name:    "CreatedAt",
			Type:    "time.Time",
			GormTag: "index;",
			Comment: "Create time",
			Order:   "DESC",
		})
		fields = append(fields, &Field{
			Name:    "UpdatedAt",
			Type:    "time.Time",
			GormTag: "index;",
			Comment: "Update time",
		})
		a.Fields = fields
	}

	for i, item := range a.Fields {
		switch item.Name {
		case "ID":
			a.Include.ID = true
		case "Status":
			a.Include.Status = true
		case "CreatedAt":
			a.Include.CreatedAt = true
		case "UpdatedAt":
			a.Include.UpdatedAt = true
		case "Sequence":
			a.Include.Sequence = true
		}
		if a.FillGormCommit && item.Comment != "" {
			if len([]byte(item.GormTag)) > 0 && !strings.HasSuffix(item.GormTag, ";") {
				item.GormTag += ";"
			}
			item.GormTag += fmt.Sprintf("comment:%s;", item.Comment)
		}
		a.Fields[i] = item.Format()
	}

	return a
}

type Field struct {
	Name      string            `yaml:"name,omitempty" json:"name,omitempty"`
	Type      string            `yaml:"type,omitempty" json:"type,omitempty"`
	GormTag   string            `yaml:"gorm_tag,omitempty" json:"gorm_tag,omitempty"`
	JSONTag   string            `yaml:"json_tag,omitempty" json:"json_tag,omitempty"`
	CustomTag string            `yaml:"custom_tag,omitempty" json:"custom_tag,omitempty"`
	Comment   string            `yaml:"comment,omitempty" json:"comment,omitempty"`
	Query     *FieldQuery       `yaml:"query,omitempty" json:"query,omitempty"`
	Order     string            `yaml:"order,omitempty" json:"order,omitempty"`
	Form      *FieldForm        `yaml:"form,omitempty" json:"form,omitempty"`
	Extra     map[string]string `yaml:"extra,omitempty" json:"extra,omitempty"`
}

func (a *Field) Format() *Field {
	if a.JSONTag != "" {
		if vv := strings.Split(a.JSONTag, ","); len(vv) > 1 {
			if vv[0] == "" {
				vv[0] = utils.ToLowerUnderlinedNamer(a.Name)
				a.JSONTag = strings.Join(vv, ",")
			}
		}
	} else {
		a.JSONTag = utils.ToLowerUnderlinedNamer(a.Name)
	}

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
		if a.Form.Name == "" {
			a.Form.Name = a.Name
		}
		if a.Form.JSONTag != "" {
			if vv := strings.Split(a.Form.JSONTag, ","); len(vv) > 1 {
				if vv[0] == "" {
					vv[0] = utils.ToLowerUnderlinedNamer(a.Name)
					a.Form.JSONTag = strings.Join(vv, ",")
				}
			}
		} else {
			a.Form.JSONTag = utils.ToLowerUnderlinedNamer(a.Name)
		}
		if a.Form.Comment == "" {
			a.Form.Comment = a.Comment
		}
	}
	return a
}

type FieldQuery struct {
	Name       string `yaml:"name,omitempty" json:"name,omitempty"`
	InQuery    bool   `yaml:"in_query,omitempty" json:"in_query,omitempty"`
	FormTag    string `yaml:"form_tag,omitempty" json:"form_tag,omitempty"`
	BindingTag string `yaml:"binding_tag,omitempty" json:"binding_tag,omitempty"`
	CustomTag  string `yaml:"custom_tag,omitempty" json:"custom_tag,omitempty"`
	Comment    string `yaml:"comment,omitempty" json:"comment,omitempty"`
	IfCond     string `yaml:"cond,omitempty" json:"cond,omitempty"`
	OP         string `yaml:"op,omitempty" json:"op,omitempty"`     // LIKE/=/</>/<=/>=/<>
	Args       string `yaml:"args,omitempty" json:"args,omitempty"` // v + "%"
}

type FieldForm struct {
	Name       string `yaml:"name,omitempty" json:"name,omitempty"`
	JSONTag    string `yaml:"json_tag,omitempty" json:"json_tag,omitempty"`
	BindingTag string `yaml:"binding_tag,omitempty" json:"binding_tag,omitempty"`
	CustomTag  string `yaml:"custom_tag,omitempty" json:"custom_tag,omitempty"`
	Comment    string `yaml:"comment,omitempty" json:"comment,omitempty"`
}
