package generate

type TplItem struct {
	StructName string         `yaml:"name"`
	Comment    string         `yaml:"comment"`
	Fields     []TplFieldItem `yaml:"fields"`
}

func (t TplItem) toSchemaFields() []schemaField {
	var items []schemaField
	for _, f := range t.Fields {
		items = append(items, schemaField{
			Name:           f.StructFieldName,
			Comment:        f.Comment,
			Type:           f.StructFieldType,
			IsRequired:     f.StructFieldRequired,
			BindingOptions: f.BindingOptions,
			GormOptions:    f.GormOptions,
		})
	}
	return items
}

type TplFieldItem struct {
	StructFieldName     string `yaml:"name"`
	StructFieldRequired bool   `yaml:"required"`
	Comment             string `yaml:"comment"`
	StructFieldType     string `yaml:"type"`
	GormOptions         string `yaml:"gorm_options"`
	BindingOptions      string `yaml:"binding_options"`
}
