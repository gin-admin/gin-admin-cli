# [gin-admin](https://github.com/LyricTian/gin-admin) efficiency assistant

> A gin-admin efficiency assistant that provides project initialization, code generation, greatly improves work efficiency, and quickly completes the development of business logic.

## Quick start

### Get and install

```bash
go install github.com/gin-admin/gin-admin-cli/v10@latest
```

### Create a new project

```bash
gin-admin-cli new -d ~/go/src --name testapp --desc 'A test API service based on golang.' --module 'github.com/xxx/testapp'
```

### Generate a new module

Using `Dictionary` as an example, the configuration file is as follows `dictionary.yaml`:

```yaml
- name: Dictionary
  comment: Dictionary management for system
  disable_pagination: true
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

### Remove a module

```bash
gin-admin-cli rm -d ~/go/src/testapp -m SYS -s Dictionary
```

## Command help

```text
NAME:
   gin-admin-cli new - Create a new project

USAGE:
   gin-admin-cli new [command options] [arguments...]

OPTIONS:
   --dir value, -d value  The directory to generate the project (default: current directory)
   --name value           The project name
   --desc value           The project description
   --version value        The project version (default: 1.0.0)
   --module value         The project module name (default: project name)
   --git-url value        Use git repository to initialize the project (default: https://github.com/LyricTian/gin-admin.git)
   --git-branch value     Use git branch to initialize the project
   --help, -h             show help
```

```text
NAME:
   gin-admin-cli generate - Generate structs to the specified module, support config file

USAGE:
   gin-admin-cli generate [command options] [arguments...]

OPTIONS:
   --dir value, -d value      The directory to generate the struct from
   --module value, -m value   The module to generate the struct from
   --tpl-type value           The template type to generate the struct from (default: crud)
   --module-path value        The module path to generate the struct from (default: internal/mods)
   --wire-path value          The wire generate path to generate the struct from (default: internal/wirex)
   --swag-path value          The swagger generate path to generate the struct from (default: internal/swagger)
   --config value, -c value   The config file or directory to generate the struct from (JSON/YAML)
   --structs value, -s value  The struct to generate (multiple structs can be separated by a comma)
   --structs-comment value    Specify the struct comment
   --structs-output value     Specify the packages to generate the struct (default: schema,dal,biz,api)
   --tpl-path value           The template path to generate the struct from (default use tpls)
   --help, -h                 show help
```

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
