package new

// TplProjectStructure 项目结构
const TplProjectStructure = `
├── LICENSE
├── Makefile
├── README.md
├── cmd
│   └── gin-admin
│       └── main.go # 入口文件
├── configs # 配置文件
│   ├── config.toml
│   ├── menu.yaml
│   └── model.conf
├── docs
│   └── data_model.md
├── go.mod
├── go.sum
├── internal
│   └── app
│       ├── api
│       │   └── mock
│       ├── bll
│       ├── config
│       ├── contextx
│       ├── ginx
│       ├── middleware
│       ├── model
│       │   └── gormx
│       ├── module
│       │   └── adapter
│       ├── router
│       ├── schema
│       ├── swagger
│       │   ├── docs.go
│       │   ├── swagger.json
│       │   └── swagger.yaml
│       ├── test
├── pkg
│   ├── auth
│   │   └── jwtauth
│   ├── errors
│   ├── logger
│   │   ├── hook
│   │   │   ├── gorm
│   └── util
│       ├── hash
│       ├── json
│       ├── structure
│       ├── trace
│       ├── uuid
│       └── yaml
└── scripts
`
