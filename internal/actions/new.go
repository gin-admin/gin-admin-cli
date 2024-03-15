package actions

import (
	"bufio"
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-admin/gin-admin-cli/v10/internal/utils"
	"go.uber.org/zap"
)

const (
	defaultGitURL = "https://github.com/LyricTian/gin-admin.git"
)

func New(cfg NewConfig) *NewAction {
	if cfg.AppName == "" {
		cfg.AppName = strings.ToLower(strings.ReplaceAll(cfg.Name, "-", ""))
	}
	if cfg.GitURL == "" {
		cfg.GitURL = defaultGitURL
	}
	if cfg.PkgName == "" {
		cfg.PkgName = cfg.Name
	}
	if cfg.Version == "" {
		cfg.Version = "v1.0.0"
	}
	if cfg.Description == "" {
		name := []byte(cfg.AppName)
		name[0] = name[0] - 32
		cfg.Description = fmt.Sprintf("%s API service", name)
	}

	return &NewAction{
		logger: zap.S().Named("[NEW]"),
		cfg:    &cfg,
	}
}

type NewConfig struct {
	Dir         string
	Name        string
	AppName     string
	PkgName     string
	Description string
	Version     string
	GitURL      string
	GitBranch   string
	FeDir       string
	FeName      string
	FeGitURL    string
	FeGitBranch string
}

type NewAction struct {
	logger *zap.SugaredLogger
	cfg    *NewConfig
}

func (a *NewAction) Run(ctx context.Context) error {
	a.logger.Infof("Create project %s in %s", a.cfg.Name, a.cfg.Dir)
	projectDir := filepath.Join(a.cfg.Dir, a.cfg.Name)
	if exists, err := utils.ExistsFile(projectDir); err != nil {
		return err
	} else if exists {
		a.logger.Warnf("Project %s already exists", a.cfg.Name)
		return nil
	}
	_ = os.MkdirAll(a.cfg.Dir, os.ModePerm)

	if err := utils.ExecGitClone(a.cfg.Dir, a.cfg.GitURL, a.cfg.GitBranch, a.cfg.Name); err != nil {
		return err
	}

	cleanFiles := []string{".git", "CHANGELOG.md", "LICENSE", "README.md", "README_CN.md", "internal/swagger/v3", "internal/wirex/wire_gen.go", "swagger.jpeg"}
	for _, f := range cleanFiles {
		if err := os.RemoveAll(filepath.Join(projectDir, f)); err != nil {
			return err
		}
	}

	a.logger.Infof("Update project info...")
	oldModuleName, err := a.getModuleName(projectDir)
	if err != nil {
		return err
	}
	oldProjectInfo, err := a.getProjectInfo(projectDir)
	if err != nil {
		return err
	}

	appName := a.cfg.AppName
	err = filepath.WalkDir(projectDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		name := d.Name()
		if name == "main.go" || name == "config.go" ||
			name == "Makefile" || name == "Dockerfile" || name == ".gitignore" {
			f, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			f = []byte(strings.ReplaceAll(string(f), oldProjectInfo.AppName, appName))
			f = []byte(strings.ReplaceAll(string(f), oldProjectInfo.Version, a.cfg.Version))
			f = []byte(strings.ReplaceAll(string(f), oldProjectInfo.Description, a.cfg.Description))
			f = []byte(strings.ReplaceAll(string(f), oldModuleName, a.cfg.PkgName))
			return os.WriteFile(path, f, info.Mode())
		}

		if name == "go.mod" || strings.HasSuffix(name, ".go") {
			return utils.ReplaceFileContent(path, []byte(oldModuleName), []byte(a.cfg.PkgName), info.Mode())
		}

		if strings.HasSuffix(name, ".toml") {
			return utils.ReplaceFileContent(path, []byte(oldProjectInfo.AppName), []byte(appName), info.Mode())
		}

		return nil
	})
	if err != nil {
		return err
	}

	err = utils.WriteFile(filepath.Join(projectDir, "README.md"), []byte(a.getReadme()))
	if err != nil {
		return err
	}

	a.logger.Infof("Generate wire and swagger files...")
	_ = utils.ExecGoModTidy(projectDir)
	_ = utils.ExecSwagGen(projectDir, "./main.go", "./internal/swagger")
	_ = utils.ExecWireGen(projectDir, "internal/wirex")
	_ = utils.ExecGitInit(projectDir)

	fmt.Println("🎉  Congratulations, your project has been created successfully.")
	fmt.Println("------------------------------------------------------------")
	fmt.Println(utils.GetDefaultProjectTree())
	fmt.Println("------------------------------------------------------------")

	fmt.Println("🚀  You can execute the following commands to start the project:")
	fmt.Println("------------------------------------------------------------")
	fmt.Printf("cd %s\n", projectDir)
	fmt.Println("make start")
	fmt.Println("------------------------------------------------------------")

	if err := a.generateFE(ctx); err != nil {
		return err
	}

	return nil
}

func (a *NewAction) getModuleName(projectDir string) (string, error) {
	f, err := os.Open(filepath.Join(projectDir, "go.mod"))
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "module") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module")), nil
		}
	}
	return "", nil
}

type projectInfo struct {
	AppName     string
	Description string
	Version     string
}

func (a *NewAction) getProjectInfo(projectDir string) (*projectInfo, error) {
	var info projectInfo
	f, err := os.Open(filepath.Join(projectDir, "main.go"))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		switch {
		case strings.HasPrefix(line, "// @title"):
			info.AppName = strings.TrimSpace(strings.TrimPrefix(line, "// @title"))
		case strings.HasPrefix(line, "// @version"):
			info.Version = strings.TrimSpace(strings.TrimPrefix(line, "// @version"))
		case strings.HasPrefix(line, "// @description"):
			info.Description = strings.TrimSpace(strings.TrimPrefix(line, "// @description"))
		}
	}
	return &info, nil
}

func (a *NewAction) getReadme() string {
	var sb strings.Builder
	sb.WriteString("# " + a.cfg.Name + "\n\n")
	sb.WriteString("> " + a.cfg.Description + "\n\n")

	sb.WriteString("## Quick Start\n\n")
	sb.WriteString("```bash\n")
	sb.WriteString("make start\n")
	sb.WriteString("```\n\n")

	sb.WriteString("## Build\n\n")
	sb.WriteString("```bash\n")
	sb.WriteString("make build\n")
	sb.WriteString("```\n\n")

	sb.WriteString("## Generate wire inject files\n\n")
	sb.WriteString("```bash\n")
	sb.WriteString("make wire\n")
	sb.WriteString("```\n\n")

	sb.WriteString("## Generate swagger documents\n\n")
	sb.WriteString("```bash\n")
	sb.WriteString("make swagger\n")
	sb.WriteString("```\n\n")

	return sb.String()
}

func (a *NewAction) generateFE(_ context.Context) error {
	if a.cfg.FeDir == "" {
		return nil
	}

	a.logger.Infof("Create frontend project %s in %s", a.cfg.FeName, a.cfg.FeDir)
	feDir, err := filepath.Abs(filepath.Join(a.cfg.FeDir, a.cfg.FeName))
	if err != nil {
		return err
	}

	if exists, err := utils.ExistsFile(feDir); err != nil {
		return err
	} else if exists {
		a.logger.Warnf("Frontend project %s already exists", a.cfg.FeName)
		return nil
	}

	_ = os.MkdirAll(a.cfg.FeDir, os.ModePerm)
	err = utils.ExecGitClone(a.cfg.FeDir, a.cfg.FeGitURL, a.cfg.FeGitBranch, a.cfg.FeName)
	if err != nil {
		return err
	}

	cleanFiles := []string{".git", "LICENSE", "README.md", "demo.png"}
	for _, file := range cleanFiles {
		if err := os.RemoveAll(filepath.Join(feDir, file)); err != nil {
			return err
		}
	}

	err = utils.WriteFile(filepath.Join(feDir, "README.md"), []byte(a.getFeReadme()))
	if err != nil {
		return err
	}

	a.logger.Infof("🎉  Frontend project %s has been created successfully", a.cfg.FeName)
	fmt.Println("------------------------------------------------------------")
	fmt.Printf("Git repository: %s\n", a.cfg.FeGitURL)
	fmt.Printf("Branch: %s\n", a.cfg.FeGitBranch)
	fmt.Printf("Directory: %s\n", feDir)
	fmt.Println("------------------------------------------------------------")

	return nil
}

func (a *NewAction) getFeReadme() string {
	var sb strings.Builder
	sb.WriteString("# " + a.cfg.FeName + "\n\n")
	sb.WriteString("> " + a.cfg.Description + "\n\n")

	sb.WriteString("## Environment Prepare\n\n")
	sb.WriteString("> You can use [nvm](https://github.com/nvm-sh/nvm) to manage node version.\n\n")
	sb.WriteString("- Node.js v16.20.2\n\n")

	sb.WriteString("## Quick Start\n\n")

	sb.WriteString("### Install dependencies\n\n")
	sb.WriteString("```bash\n")
	sb.WriteString("npm install\n")
	sb.WriteString("```\n\n")

	sb.WriteString("### Start project\n\n")
	sb.WriteString("```bash\n")
	sb.WriteString("npm start\n")
	sb.WriteString("```\n\n")

	sb.WriteString("### Build project\n\n")
	sb.WriteString("```bash\n")
	sb.WriteString("npm run build\n")
	sb.WriteString("```\n\n")

	return sb.String()
}
