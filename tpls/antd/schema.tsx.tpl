export default {
  apiName: '{%lowerHyphensPlural .Name%}',
  title: '{%.Comment%}',
  search: {
    type: 'object',
    layoutAuto: true,
    properties: {
      {%- range .Fields%}{%$fieldName := .Name%}{%$fieldType := .Type%}{%$fieldComment := .Comment%}{%$fieldExtra := .Extra%}
      {%- if $fieldExtra.search%}
      {%- if eq $fieldName "Status"%}
      status: {
        title: '{%$fieldComment%}',
        type: 'string',
        widget: 'select',
        props: {
          options: [
            {
              label: '启用',
              value: 'disabled',
            },
            {
              label: '禁用',
              value: 'enabled',
            },
          ],
        },
      },
      {%- else%}
      {%lowerCamel $fieldName%}: { title: '{%$fieldComment%}', type: '{%convGoTypeToTsType $fieldType%}', widget: '{%with $fieldExtra.widget%}{%.%}{%else%}Input{%end%}' },
      {%- end%}
      {%- end%}
      {%- end%}
    },
  },
  columns: [
    {%- range .Fields%}{%$fieldName := .Name%}{%$fieldComment := .Comment%}{%$fieldExtra := .Extra%}
    {%- if $fieldExtra.column%}
    {%- if eq $fieldName "Status"%}
    {
      title: '{%$fieldComment%}',
      dataIndex: 'status',
      valueType: 'tag',
      valueTypeProps: (value: any) => ({
        color: value === 'disabled' ? 'red' : 'blue',
      }),
      enum: {
        enabled: '启用',
        disabled: '禁用',
      },
    },
    {%- else%}
    { title: '{%$fieldComment%}', dataIndex: '{%lowerUnderline $fieldName%}' },
    {%- end%}
    {%- end%}
    {%- end%}
    { title: '创建时间', dataIndex: 'created_at', valueType: 'dateTime' },
    { title: '更新时间', dataIndex: 'updated_at', valueType: 'dateTime' },
  ],
  form: {
    type: 'object',
    displayType: 'row',
    fieldCol: 22,
    labelCol: 6,
    column: 2,
    properties: {
    {%- range .Fields%}{%$fieldName := .Name%}{%$fieldType := .Type%}{%$fieldComment := .Comment%}{%$fieldExtra := .Extra%}
    {%- if .Form%}
    {%- if eq $fieldName "Status"%}
      status: {
        title: '{%$fieldComment%}',
        type: 'string',
        widget: 'radio',
        default: 'enabled',
        required: true,
        props: {
          options: [
            {
              label: '启用',
              value: 'enabled',
            },
            {
              label: '禁用',
              value: 'disabled',
            },
          ],
        },
      },
    {%- else%}
     {%lowerUnderline $fieldName%}: {
        title: '{%$fieldComment%}',
        type: '{%convGoTypeToTsType $fieldType%}',
        widget: '{%with $fieldExtra.widget%}{%.%}{%else%}Input{%end%}',
        required: {%with $fieldExtra.required%}{%.%}{%else%}false{%end%},
        {%- with $fieldExtra.cellSpan%}
        cellSpan: 2,
        labelCol: 3,
        {%- end%}
        {%- if $fieldExtra.default%}
        default: {%if eq $fieldType "int"%}{%$fieldExtra.default%}{%else%}'{%$fieldExtra.default%}'{%end%}
        {%- end%}
        {%- with $fieldExtra.format%}
        format: '{%.%}'
        {%- end%}
        {%- if eq $fieldExtra.widget "XUpload"%}
        props: {
          data: {
            bu_type: '{%$fieldExtra.buType%}',
          },
        },
        {%- end%}
     },
    {%- end%}
    {%- end%}
    {%- end%}
    },
  },
};
