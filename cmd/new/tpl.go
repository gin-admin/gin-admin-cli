package new

// TplProjectStructure 项目结构
const TplProjectStructure = `
├── cmd # 主服务（程序入口）
├── configs # 配置文件目录(包含运行配置参数及casbin模型配置)
├── docs # 文档目录
├── internal # 内部代码模块
│   └── app # 内部应用模块入口
│       ├── api # API控制器模块
│       │   └── mock # API Mock模块(包括swagger的注释描述)
│       ├── bll # 业务逻辑模块接口
│       │   └── impl
│       │       └── bll # 业务逻辑模块接口的实现
│       ├── config # 配置参数(与config.toml一一映射)
│       ├── context # 统一上下文模块
│       ├── ginplus # gin的扩展模块
│       ├── initialize # 初始化模块（提供依赖模块的初始化函数及依赖注入的初始化）
│       ├── middleware # gin中间件模块
│       ├── model # 存储层模块接口
│       │   └── impl
│       │       └── gorm
│       │           ├── entity # 与数据库表及字段的映射实体
│       │           └── model # 存储层模块接口的gorm实现
│       │       └── mongo
│       │           ├── entity # 与数据库表及字段的映射实体
│       │           └── model # 存储层模块接口的mongo实现
│       ├── module # 内部模块间依赖的公共模块
│       ├── router # gin的路由模块
│       ├── schema # 提供Request/Response的对象模块
│       ├── swagger # swagger配置及自动生成的文件
│       └── test # API的单元测试
├── pkg # 公共模块
│   ├── auth # JWT认证模块
│   ├── errors # 统一错误处理模块
│   ├── logger # 日志模块
│   |── unique # 唯一ID
│   └── util # 工具库模块
├── scripts # 脚本目录
`
