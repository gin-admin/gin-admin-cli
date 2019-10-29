package generate

// TplItem 模板项
type TplItem struct {
	StructName string         `json:"struct_name"` // 结构体名称
	Comment    string         `json:"comment"`     // 注释
	Fields     []TplFieldItem `json:"fields"`      // 字段项
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
		})
	}
	return items
}

func (t TplItem) toEntityFields() []entityField {
	var items []entityField
	for _, f := range t.Fields {
		items = append(items, entityField{
			Name:        f.StructFieldName,
			Comment:     f.Comment,
			Type:        f.StructFieldType,
			GormOptions: f.GormOptions,
		})
	}
	return items
}

// TplFieldItem 模板字段项
type TplFieldItem struct {
	StructFieldName     string `json:"struct_field_name"`     // 结构体字段名称
	StructFieldRequired bool   `json:"struct_field_required"` // 结构字段必选项
	Comment             string `json:"comment"`               // 注释
	StructFieldType     string `json:"struct_field_type"`     // 结构体字段类型
	GormOptions         string `json:"gorm_options"`          // gorm配置项
	BindingOptions      string `json:"binding_options"`       // binding配置项
}

// {
// 	"struct_name": "Task",
// 	"comment": "任务管理",
// 	"fields": [
// 	  {
// 		"struct_field_name": "RecordID",
// 		"comment": "记录ID",
// 		"struct_field_required": false,
// 		"struct_field_type": "string",
// 		"gorm_options": "size:36;index;"
// 	  },
// 	  {
// 		"struct_field_name": "Name",
// 		"comment": "任务名称",
// 		"struct_field_required": true,
// 		"struct_field_type": "string",
// 		"gorm_options": "size:50;index;"
// 	  },
// 	  {
// 		"struct_field_name": "Memo",
// 		"comment": "备注",
// 		"struct_field_required": false,
// 		"struct_field_type": "string",
// 		"gorm_options": "size:500;"
// 	  },
// 	  {
// 		"struct_field_name": "Creator",
// 		"comment": "创建者",
// 		"struct_field_required": false,
// 		"struct_field_type": "string",
// 		"gorm_options": "size:36;index;"
// 	  }
// 	]
//   }
