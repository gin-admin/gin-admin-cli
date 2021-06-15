package generate

import (
	"context"
	"fmt"

	"github.com/gin-admin/gin-admin-cli/v4/util"
)

func getModelImplMongoFileName(dir, name string) string {
	fullname := fmt.Sprintf("%s/internal/app/model/impl/mongo/repo/%s.repo.go", dir, util.ToLowerUnderlinedNamer(name))
	return fullname
}

// 生成model实现文件
func genModelImplMongo(ctx context.Context, pkgName, dir, name, comment string) error {
	data := map[string]interface{}{
		"PkgName":    pkgName,
		"Name":       name,
		"PluralName": util.ToPlural(name),
		"Comment":    comment,
	}

	buf, err := execParseTpl(modelImplMongoTpl, data)
	if err != nil {
		return err
	}

	fullname := getModelImplMongoFileName(dir, name)
	err = createFile(ctx, fullname, buf)
	if err != nil {
		return err
	}

	fmt.Printf("文件[%s]写入成功\n", fullname)

	return execGoFmt(fullname)
}

const modelImplMongoTpl = `
package repo

import (
	"context"
	"time"

	"{{.PkgName}}/internal/app/model"
	"{{.PkgName}}/internal/app/model/impl/mongo/entity"
	"{{.PkgName}}/internal/app/schema"
	"{{.PkgName}}/pkg/errors"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ model.I{{.Name}} = (*{{.Name}})(nil)

// {{.Name}}Set 注入{{.Name}}
var {{.Name}}Set = wire.NewSet(wire.Struct(new({{.Name}}), "*"), wire.Bind(new(model.I{{.Name}}), new(*{{.Name}})))

// {{.Name}} {{.Comment}}存储
type {{.Name}} struct {
	Client *mongo.Client
}

func (a *{{.Name}}) getQueryOption(opts ...schema.{{.Name}}QueryOptions) schema.{{.Name}}QueryOptions {
	var opt schema.{{.Name}}QueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

// Query 查询数据
func (a *{{.Name}}) Query(ctx context.Context, params schema.{{.Name}}QueryParam, opts ...schema.{{.Name}}QueryOptions) (*schema.{{.Name}}QueryResult, error) {
	opt := a.getQueryOption(opts...)

	c := entity.Get{{.Name}}Collection(ctx, a.Client)
	filter := DefaultFilter(ctx)

	// TODO: 查询条件

	opt.OrderFields = append(opt.OrderFields, schema.NewOrderField("_id", schema.OrderByDESC))

	var list entity.{{.PluralName}}
	pr, err := WrapPageQuery(ctx, c, params.PaginationParam, filter, &list, options.Find().SetSort(ParseOrder(opt.OrderFields)))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.{{.Name}}QueryResult{
		PageResult: pr,
		Data:       list.ToSchema{{.PluralName}}(),
	}

	return qr, nil
}

// Get 查询指定数据
func (a *{{.Name}}) Get(ctx context.Context, id string, opts ...schema.{{.Name}}GetOptions) (*schema.{{.Name}}, error) {
	c := entity.Get{{.Name}}Collection(ctx, a.Client)
	filter := DefaultFilter(ctx, Filter("_id", id))
	var item entity.{{.Name}}
	ok, err := FindOne(ctx, c, filter, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item.ToSchema{{.Name}}(), nil
}

// Create 创建数据
func (a *{{.Name}}) Create(ctx context.Context, item schema.{{.Name}}) error {
	eitem := entity.Schema{{.Name}}(item).To{{.Name}}()
	eitem.CreatedAt = time.Now()
	eitem.UpdatedAt = time.Now()
	c := entity.Get{{.Name}}Collection(ctx, a.Client)
	err := Insert(ctx, c, eitem)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Update 更新数据
func (a *{{.Name}}) Update(ctx context.Context, id string, item schema.{{.Name}}) error {
	eitem := entity.Schema{{.Name}}(item).To{{.Name}}()
	eitem.UpdatedAt = time.Now()
	c := entity.Get{{.Name}}Collection(ctx, a.Client)
	err := Update(ctx, c, DefaultFilter(ctx, Filter("_id", id)), eitem)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Delete 删除数据
func (a *{{.Name}}) Delete(ctx context.Context, id string) error {
	c := entity.Get{{.Name}}Collection(ctx, a.Client)
	err := Delete(ctx, c, DefaultFilter(ctx, Filter("_id", id)))
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

`
