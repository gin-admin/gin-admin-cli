package actions

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"text/template"

	"github.com/gin-admin/gin-admin-cli/v10/internal/schema"
	"github.com/gin-admin/gin-admin-cli/v10/internal/tfs"
	"github.com/gin-admin/gin-admin-cli/v10/internal/utils"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

type GenAntdConfig struct {
	Dir string
}

func GenAntd(cfg GenAntdConfig) *GenAntdAction {
	return &GenAntdAction{
		logger: zap.S().Named("[GEN-ANTD]"),
		cfg:    &cfg,
	}
}

type GenAntdAction struct {
	logger *zap.SugaredLogger
	cfg    *GenAntdConfig
}

// Run generate command
func (a *GenAntdAction) RunWithConfig(ctx context.Context, cfgName string) error {
	var parseFile = func(name string) ([]*schema.S, error) {
		var data []*schema.S
		switch filepath.Ext(name) {
		case ".json":
			if err := utils.ParseJSONFile(name, &data); err != nil {
				return nil, err
			}
		case ".yaml", ".yml":
			if err := utils.ParseYAMLFile(name, &data); err != nil {
				return nil, err
			}
		default:
			a.logger.Warnf("Ignore file %s, only support json/yaml/yml", name)
		}
		if len(data) == 0 {
			return nil, nil
		}
		return data, nil
	}

	if utils.IsDir(cfgName) {
		var data []*schema.S
		err := filepath.WalkDir(cfgName, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			items, err := parseFile(path)
			if err != nil {
				return err
			}
			data = append(data, items...)
			return nil
		})
		if err != nil {
			return err
		}
		return a.run(ctx, data)
	}

	data, err := parseFile(cfgName)
	if err != nil {
		return err
	} else if len(data) == 0 {
		a.logger.Warnf("No data found in file %s", cfgName)
		return nil
	}
	return a.run(ctx, data)
}

func (a *GenAntdAction) run(ctx context.Context, data []*schema.S) error {
	if len(data) == 0 {
		return fmt.Errorf("no data found")
	}
	for _, s := range data {
		if err := a.runWithS(ctx, s); err != nil {
			return err
		}
	}
	return nil
}

func (a *GenAntdAction) runWithS(_ context.Context, s *schema.S) error {
	sort.Sort(s.Fields)
	s = s.Format()
	schemaFileName := cast.ToString(s.Extra["schema"])
	if schemaFileName == "" {
		return fmt.Errorf("Struct %s schema file name is empty", s.Name)
	}

	schemaFile := filepath.Join(a.cfg.Dir, schemaFileName)
	exists, err := utils.ExistsFile(schemaFile)
	if err != nil {
		return err
	} else if !s.ForceWrite && exists {
		a.logger.Warnf("File %s exists, ignore", schemaFile)
		return nil
	}

	_ = os.MkdirAll(filepath.Dir(schemaFile), os.ModePerm)

	tplData, err := tfs.Ins.ReadFile("antd/schema.tsx.tpl")
	if err != nil {
		return err
	}

	t, err := template.New("").Funcs(utils.FuncMap).Delims("{%", "%}").Parse(string(tplData))
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	if err := t.Execute(buf, s); err != nil {
		return err
	}

	a.logger.Info("Generate file:", schemaFile)
	return os.WriteFile(schemaFile, buf.Bytes(), os.ModePerm)
}
