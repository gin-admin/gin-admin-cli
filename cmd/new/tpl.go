package new

// TplProjectStructure 项目结构
const TplProjectStructure = `

├── cmd
│   └── %s
│       └── main.go       # 入口文件
├── configs
│   ├── config.toml       # 配置文件
│   ├── menu.yaml         # 菜单初始化配置
│   └── model.conf        # casbin 策略配置
├── docs                  # 文档
├── internal
│   └── app
│       ├── api           # API 处理层
│       ├── config        # 配置文件映射
│       ├── contextx      # 统一上下文处理
│       ├── dao           # 数据访问层
│       ├── ginx          # gin 扩展模块
│       ├── middleware    # gin 中间件模块
│       ├── module        # 通用业务处理模块
│       ├── router        # 路由层
│       ├── schema        # 统一入参、出参对象映射
│       ├── service       # 业务逻辑层
│       ├── swagger       # swagger 生成文件
│       ├── test          # 模块单元测试
├── pkg
│   ├── auth              
│   │   └── jwtauth       # jwt 认证模块
│   ├── errors            # 错误处理模块
│   ├── gormx             # gorm 扩展模块
│   ├── logger            # 日志模块
│   │   ├── hook
│   └── util              # 工具包
│       ├── conv         
│       ├── hash         
│       ├── json
│       ├── snowflake
│       ├── structure
│       ├── trace
│       ├── uuid
│       └── yaml
└── scripts               # 统一处理脚本

`
