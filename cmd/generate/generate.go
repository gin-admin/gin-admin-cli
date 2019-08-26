package generate

// Config 配置参数
type Config struct {
	Dir     string
	PkgName string
	Name    string
	Comment string
	File    string
}

// Command 生成命令
type Command struct {
	cfg *Config
}

// Exec 执行命令
func (a *Command) Exec() error {

	return nil
}
