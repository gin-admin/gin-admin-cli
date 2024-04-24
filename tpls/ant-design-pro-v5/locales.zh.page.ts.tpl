{{- $name := .Name}}
{{- $lowerCamelName := lowerCamel .Name}}
{{- $parentName := .Extra.ParentName}}
export default {
    'pages.{{with $parentName}}{{.}}.{{end}}{{$lowerCamelName}}.add': 'Add {{$name}}',
    'pages.{{with $parentName}}{{.}}.{{end}}{{$lowerCamelName}}.edit': 'Edit {{$name}}',
    'pages.{{with $parentName}}{{.}}.{{end}}{{$lowerCamelName}}.delTip': 'Are you sure you want to delete this record?',
    {{- range .Fields}}
    'pages.{{with $parentName}}{{.}}.{{end}}{{$lowerCamelName}}.form.{{lowerUnderline .Name}}': '{{.Name}}',
    'pages.{{with $parentName}}{{.}}.{{end}}{{$lowerCamelName}}.form.{{lowerUnderline .Name}}.placeholder': 'Please enter the {{lowerSpace .Name}}',
    'pages.{{with $parentName}}{{.}}.{{end}}{{$lowerCamelName}}.form.{{lowerUnderline .Name}}.required': '{{titleSpace .Name}} is required!',
    {{- end}}
};
