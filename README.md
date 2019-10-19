# gin-admin-cli - [gin-admin](https://github.com/LyricTian/gin-admin)

> GinAdmin辅助工具

## 下载并使用

```bash
$ go get -u -v github.com/LyricTian/gin-admin-cli
```

### 创建项目

```bash
USAGE:
   gin-admin-cli new [command options] [arguments...]

OPTIONS:
   --dir value, -d value  项目生成目录
   --pkg value, -p value  项目包名
   --mirror, -m           使用国内镜像(gitee.com)
   --web, -w              包含web项目(下载gin-admin-react项目到web目录)
```

> 使用示例

```
$ gin-admin-cli new -m -d ~/go/src/test-gin-admin -p test-gin-admin
```

### 生成业务模块

#### 指定模块名称和说明生成模块

```bash
USAGE:
   gin-admin-cli generate [command options] [arguments...]

OPTIONS:
   --dir value, -d value      项目生成目录
   --pkg value, -p value      项目包名
   --name value, -n value     模块名称(结构体名称)
   --comment value, -c value  模块注释
   --file value, -f value     指定模板文件(.json，模板配置可参考说明)
```

> 使用示例

```bash
$ gin-admin-cli g -d ./test-gin-admin -p test-gin-admin -n Task -c '任务管理'
```

#### 指定配置文件生成模块

```bash
$ gin-admin-cli g -d 项目目录 -p 包名 -f 配置文件(json)
```

> 配置文件说明

```json
{
  "struct_name": "结构体名称",
  "comment": "结构体注释说明",
  "fields": [
    {
      "struct_field_name": "结构体字段名称",
      "comment": "结构体字段注释",
      "struct_field_required": "结构体字段是否是必选项",
      "struct_field_type": "结构体字段类型",
      "gorm_options": "gorm配置项"
    }
  ]
}
```

> 使用示例

> 创建`task.json`文件

```json
{
  "struct_name": "Task",
  "comment": "任务管理",
  "fields": [
    {
      "struct_field_name": "RecordID",
      "comment": "记录ID",
      "struct_field_required": false,
      "struct_field_type": "string",
      "gorm_options": "size:36;index;"
    },
    {
      "struct_field_name": "Name",
      "comment": "任务名称",
      "struct_field_required": true,
      "struct_field_type": "string",
      "gorm_options": "size:50;index;"
    },
    {
      "struct_field_name": "Memo",
      "comment": "备注",
      "struct_field_required": false,
      "struct_field_type": "string",
      "gorm_options": "size:500;"
    },
    {
      "struct_field_name": "Creator",
      "comment": "创建者",
      "struct_field_required": false,
      "struct_field_type": "string",
      "gorm_options": "size:36;index;"
    }
  ]
}
```

```bash
$ gin-admin-cli g -d ./test-gin-admin -p test-gin-admin -f task.json
```

## MIT License

    Copyright (c) 2019 Lyric