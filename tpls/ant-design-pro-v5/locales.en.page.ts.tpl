{{- $name := .Name}}
{{- $lowerCamelName := lowerCamel .Name}}
{{- $parentName := .Extra.ParentName}}
export default {
    'pages.{{$parentName}}.{{$lowerCamelName}}.add': 'Add {{$name}}',
    'pages.{{$parentName}}.{{$lowerCamelName}}.edit': 'Edit {{$name}}',
    'pages.{{$parentName}}.{{$lowerCamelName}}.delTip': 'Are you sure you want to delete this record?',
    {{- range .Fields}}
    'pages.{{$parentName}}.{{$lowerCamelName}}.form.{{lowerUnderline .Name}}': '{{.Name}}',
    'pages.{{$parentName}}.{{$lowerCamelName}}.form.{{lowerUnderline .Name}}.placeholder': 'Please enter the {{lowerSpace .Name}}',
    'pages.{{$parentName}}.{{$lowerCamelName}}.form.{{lowerUnderline .Name}}.required': '{{titleSpace .Name}} is required!',
    {{- end}}
};
