package generate

import (
	"context"
	"fmt"
	"strings"

	"github.com/gin-admin/gin-admin-cli/v5/util"
)

func getWebServiceFileName(dir, name string) string {
	fullname := fmt.Sprintf("%s/src/services/%s.js", dir, name)
	return fullname
}

func genWebService(ctx context.Context, cmd *Command, item TplItem) error {
	name := strings.ToLower(item.StructName)
	data := map[string]interface{}{
		"StructName": util.ToPlural(name),
	}

	buf, err := execParseTpl(webServiceTpl, data)
	if err != nil {
		return err
	}

	fullname := getWebServiceFileName(cmd.cfg.React, name)
	err = createFile(ctx, fullname, buf)
	if err != nil {
		return err
	}

	fmt.Printf("File write success: %s\n", fullname)

	return execGoFmt(fullname)
}

const webServiceTpl = `import { C, D, R, U } from '@/services/base';

const api = '/api/v1/{{.StructName}}';
const query = async (params, options) => {
  return R(api, params, options);
};

const update = async (params, options) => {
  return U(api, params, options);
};

const create = async (params, options) => {
  return C(api, params, options);
};

const remove = async (params, options) => {
  return D(api, params, options);
};

export default {
  query,
  update,
  create,
  remove,
};
`
