# gin-admin-cli - [gin-admin](https://github.com/LyricTian/gin-admin)

> GinAdmin辅助工具，提供创建项目、快速生成功能模块的功能

## 下载并使用

```bash
$ go get -u -v github.com/LyricTian/gin-admin-cli
```

### 创建项目

```bash
USAGE:
   gin-admin-cli new [command options] [arguments...]

OPTIONS:
   --dir value, -d value     项目生成目录(默认GOPATH)
   --pkg value, -p value     项目包名
   --branch value, -b value  指定分支(默认master)
   --mirror, -m              使用国内镜像(gitee.com)
   --web, -w                 包含web项目
```

> 使用示例

```
$ gin-admin-cli new -p test-gin-admin -m
```

### 生成业务模块

#### 指定模块名称和说明生成模块

```bash
USAGE:
   gin-admin-cli generate [command options] [arguments...]

OPTIONS:
   --dir value, -d value      项目生成目录(默认GOPATH)
   --pkg value, -p value      项目包名
   --name value, -n value     业务模块名称(结构体名称)
   --comment value, -c value  业务模块注释(结构体注释)
   --file value, -f value     指定模板文件(.yaml，模板配置可参考说明)
   --module value, -m value   指定生成模块（默认生成全部模块，以逗号分隔，支持：schema,model,bll,api,mock,router）
   --storage value, -s value  指定model的实现存储（默认gorm，支持：mongo/gorm）
```

> 使用示例

```bash
$ gin-admin-cli g -p test-gin-admin -n Task -c '任务管理'
```

#### 指定配置文件生成模块

```bash
$ gin-admin-cli g -p 包名 -f 配置文件(yaml)
```

> 配置文件说明

```yaml
---
name: 结构体名称
comment: 结构体注释说明
fields:
- name: 结构体字段名称
  type: 结构体字段类型
  comment: 结构体字段注释
  required: 结构体字段是否是必选项
  binding_options: binding配置项（不包含required，required由required字段控制）
  gorm_options: gorm配置项
```

> 使用示例

> 创建`task.yaml`文件

``` yaml
name: Task
comment: 任务管理
fields:
  - name: Code
    type: string
    comment: 任务编号
    required: true
    binding_options: ""
    gorm_options: "size:50;index;"
  - name: Name
    type: string
    comment: 任务名称
    required: true
    binding_options: ""
    gorm_options: "size:50;index;"
  - name: Memo
    type: string
    comment: 任务备注
    required: false
    binding_options: ""
    gorm_options: "size:1024;"
```

```bash
$ gin-admin-cli g -p test-gin-admin -f task.yaml
```

## MIT License

    Copyright (c) 2020 Lyric