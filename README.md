# [GIN-Admin](https://github.com/LyricTian/gin-admin) efficiency assistant

> A gin-admin efficiency assistant that provides project initialization, code generation, greatly improves work efficiency, and quickly completes the development of business logic.

## Dependencies

- [Go](https://golang.org/) 1.19+
- [Wire](github.com/google/wire) `go install github.com/google/wire/cmd/wire@latest`
- [Swag](github.com/swaggo/swag) `go install github.com/swaggo/swag/cmd/swag@latest`

## Quick start

### Get and install

```bash
go install github.com/gin-admin/gin-admin-cli/v10@latest
```

### Create a new project

```bash
gin-admin-cli new -d ~/go/src --name testapp --desc 'A test API service based on golang.' --pkg 'github.com/xxx/testapp'
```

### Quick generate a struct

```bash
gin-admin-cli gen -d ~/go/src/testapp -m SYS --structs Dictionary --structs-comment "Dictionaries management" --structs-router-prefix
```

### Use config file to generate struct

> More examples can be found in the [examples directory](https://github.com/gin-admin/gin-admin-cli/tree/master/examples)

Using `Dictionary` as an example, the configuration file is as follows `dictionary.yaml`:

```yaml
- name: Dictionary
  comment: Dictionaries management
  disable_pagination: true
  fill_gorm_commit: true
  fill_router_prefix: true
  rewrite:
    schema: true
    dal: true
    biz: true
    api: true
  tpl_type: "tree"
  fields:
    - name: Code
      type: string
      comment: Code of dictionary (unique for same parent)
      gorm_tag: "size:32;"
      form:
        binding_tag: "required,max=32"
    - name: Name
      type: string
      comment: Display name of dictionary
      gorm_tag: "size:128;index"
      query:
        name: LikeName
        in_query: true
        form_tag: name
        op: LIKE
      form:
        binding_tag: "required,max=128"
    - name: Description
      type: string
      comment: Details about dictionary
      gorm_tag: "size:1024"
      form: {}
    - name: Sequence
      type: int
      comment: Sequence for sorting
      gorm_tag: "index;"
      order: DESC
      form: {}
    - name: Status
      type: string
      comment: Status of dictionary (disabled, enabled)
      gorm_tag: "size:20;index"
      query: {}
      form:
        binding_tag: "required,oneof=disabled enabled"
```

```bash
./gin-admin-cli gen -d ~/go/src/testapp -m SYS -c dictionary.yaml
```

### Use `ant-design-pro-v5` template to generate struct and UI

```yaml
- name: Parameter
  comment: System parameter management
  generate_fe: true
  fe_tpl: "ant-design-pro-v5"
  extra:
    ParentName: system
    IndexAntdImport: "Tag"
    IndexProComponentsImport: ""
    FormAntdImport: ""
    FormProComponentsImport: "ProFormText, ProFormTextArea, ProFormSwitch"
    IncludeCreatedAt: true
    IncludeUpdatedAt: true
  fe_mapping:
    services.typings.d.ts.tpl: "src/services/system/parameter/typings.d.ts"
    services.index.ts.tpl: "src/services/system/parameter/index.ts"
    pages.components.form.tsx.tpl: "src/pages/system/Parameter/components/SaveForm.tsx"
    pages.index.tsx.tpl: "src/pages/system/Parameter/index.tsx"
    locales.en.page.ts.tpl: "src/locales/en-US/pages/system.parameter.ts"
  fields:
    - name: Name
      type: string
      comment: Name of parameter
      gorm_tag: "size:128;index"
      unique: true
      query:
        name: LikeName
        in_query: true
        form_tag: name
        op: LIKE
      form:
        binding_tag: "required,max=128"
      extra:
        ColumnComponent: >
          {
            title: intl.formatMessage({ id: 'pages.system.parameter.form.name' }),
            dataIndex: 'name',
            width: 160,
            key: 'name', // Query field name
          }
        FormComponent: >
          <ProFormText
           name="name"
           label={intl.formatMessage({ id: 'pages.system.parameter.form.name' })}
           placeholder={intl.formatMessage({ id: 'pages.system.parameter.form.name.placeholder' })}
           rules={[
             {
               required: true,
               message: intl.formatMessage({ id: 'pages.system.parameter.form.name.required' }),
             },
           ]}
           colProps={{ span: 24 }}
           labelCol={{ span: 2 }}
          />
    - name: Value
      type: string
      comment: Value of parameter
      gorm_tag: "size:1024;"
      form:
        binding_tag: "max=1024"
      extra:
        ColumnComponent: >
          {
            title: intl.formatMessage({ id: 'pages.system.parameter.form.value' }),
            dataIndex: 'value',
            width: 200,
          }
        FormComponent: >
          <ProFormTextArea
            name="value"
            label={intl.formatMessage({ id: 'pages.system.parameter.form.value' })}
            fieldProps={{ rows: 3 }}
            colProps={{ span: 24 }}
            labelCol={{ span: 2 }}
          />
    - name: Remark
      type: string
      comment: Remark of parameter
      gorm_tag: "size:255;"
      form: {}
      extra:
        ColumnComponent: >
          {
            title: intl.formatMessage({ id: 'pages.system.parameter.form.remark' }),
            dataIndex: 'remark',
            width: 180,
          }
        FormComponent: >
          <ProFormText
            name="remark"
            label={intl.formatMessage({ id: 'pages.system.parameter.form.remark' })}
            colProps={{ span: 24 }}
            labelCol={{ span: 2 }}
          />
    - name: Status
      type: string
      comment: Status of parameter (enabled, disabled)
      gorm_tag: "size:20;index"
      query:
        in_query: true
      form:
        binding_tag: "required,oneof=enabled disabled"
      extra:
        ColumnComponent: >
          {
            title: intl.formatMessage({ id: 'pages.system.parameter.form.status' }),
            dataIndex: 'status',
            width: 130,
            search: false,
            render: (status) => {
              return (
                <Tag color={status === 'enabled' ? 'success' : 'error'}>
                  {status === 'enabled'
                    ? intl.formatMessage({ id: 'pages.system.parameter.form.status.enabled', defaultMessage: 'Enabled' })
                    : intl.formatMessage({ id: 'pages.system.parameter.form.status.disabled', defaultMessage: 'Disabled' })}
                </Tag>
              );
            },
          }
        FormComponent: >
          <ProFormSwitch
            name="statusChecked"
            label={intl.formatMessage({ id: 'pages.system.parameter.form.status' })}
            fieldProps={{
              checkedChildren: intl.formatMessage({ id: 'pages.system.parameter.form.status.enabled', defaultMessage: 'Enabled' }),
              unCheckedChildren: intl.formatMessage({ id: 'pages.system.parameter.form.status.disabled', defaultMessage: 'Disabled' }),
            }}
            colProps={{ span: 24 }}
            labelCol={{ span: 2 }}
          />
```

```bash
./gin-admin-cli gen -d ~/go/src/testapp -m SYS -c parameter.yaml
```

### Remove a struct from the module

```bash
gin-admin-cli rm -d ~/go/src/testapp -m SYS -s Dictionary
```

## Command help

### New command

```text
NAME:
   gin-admin-cli new - Create a new project

USAGE:
   gin-admin-cli new [command options] [arguments...]

OPTIONS:
   --dir value, -d value  The directory to generate the project (default: current directory)
   --name value           The project name
   --app-name value       The application name (default: project name)
   --desc value           The project description
   --version value        The project version (default: 1.0.0)
   --pkg value            The project package name (default: project name)
   --git-url value        Use git repository to initialize the project (default: https://github.com/LyricTian/gin-admin.git)
   --git-branch value     Use git branch to initialize the project (default: main)
   --fe-dir value         The frontend directory to generate the project (if empty, the frontend project will not be generated)
   --fe-name value        The frontend project name (default: frontend)
   --fe-git-url value     Use git repository to initialize the frontend project (default: https://github.com/gin-admin/gin-admin-frontend.git)
   --fe-git-branch value  Use git branch to initialize the frontend project (default: main)
   --help, -h             show help
```

### Generate command

```text
NAME:
   gin-admin-cli generate - Generate structs to the specified module, support config file

USAGE:
   gin-admin-cli generate [command options] [arguments...]

OPTIONS:
   --dir value, -d value     The project directory to generate the struct
   --module value, -m value  The module to generate the struct from (like: RBAC)
   --module-path value       The module path to generate the struct from (default: internal/mods)
   --wire-path value         The wire generate path to generate the struct from (default: internal/wirex)
   --swag-path value         The swagger generate path to generate the struct from (default: internal/swagger)
   --config value, -c value  The config file or directory to generate the struct from (JSON/YAML)
   --structs value           The struct name to generate
   --structs-comment value   Specify the struct comment
   --structs-router-prefix   Use module name as router prefix (default: false)
   --structs-output value    Specify the packages to generate the struct (default: schema,dal,biz,api)
   --tpl-path value          The template path to generate the struct from (default use tpls)
   --tpl-type value          The template type to generate the struct from (default: default)
   --fe-dir value            The frontend project directory to generate the UI
   --help, -h                show help
```

### Remove command

```text
NAME:
   gin-admin-cli remove - Remove structs from the module

USAGE:
   gin-admin-cli remove [command options] [arguments...]

OPTIONS:
   --dir value, -d value      The directory to remove the struct from
   --module value, -m value   The module to remove the struct from
   --module-path value        The module path to remove the struct from (default: internal/mods)
   --structs value, -s value  The struct to remove (multiple structs can be separated by a comma)
   --config value, -c value   The config file to generate the struct from (JSON/YAML)
   --wire-path value          The wire generate path to remove the struct from (default: internal/library/wirex)
   --swag-path value          The swagger generate path to remove the struct from (default: internal/swagger)
   --help, -h                 show help
```

## MIT License

```text
Copyright (c) 2023 Lyric
```
