# [gin-admin](https://github.com/LyricTian/gin-admin) efficiency assistant

> A gin-admin efficiency assistant that provides project initialization, code generation, greatly improves work efficiency, and quickly completes the development of business logic.

## Get and usage

```bash
go install github.com/gin-admin/gin-admin-cli/v10@latest

gin-admin-cli help gen
```

```
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
