package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"strings"
)

type astModuleMainVisitor struct {
	fset *token.FileSet
	args BasicArgs
}

func (v *astModuleMainVisitor) Visit(node ast.Node) ast.Visitor {
	switch x := node.(type) {
	case *ast.GenDecl:
		if x.Tok == token.IMPORT && len(x.Specs) > 0 {
			v.modifyModuleImport(x)
		}
	case *ast.TypeSpec:
		if x.Name.Name == v.args.ModuleName {
			if xst, ok := x.Type.(*ast.StructType); ok {
				v.modifyStructField(xst)
			}
		}
	case *ast.FuncDecl:
		if x.Name.Name == "AutoMigrate" {
			v.modifyAutoMigrate(x)
		} else if x.Name.Name == "RegisterV1Routers" {
			v.modifyRegisterV1Routers(x)
		}
	}
	return v
}

func (v *astModuleMainVisitor) modifyModuleImport(x *ast.GenDecl) {
	if v.args.Flag&AstFlagGen != 0 {
		for _, pkgName := range []string{StructPackageAPI, StructPackageSchema} {
			findIndex := -1
			modulePath := GetModuleImportPath(v.args.Dir, v.args.ModulePath, v.args.ModuleName) + "/" + pkgName
			for i, spec := range x.Specs {
				if is, ok := spec.(*ast.ImportSpec); ok &&
					is.Path.Value == fmt.Sprintf("\"%s\"", modulePath) {
					findIndex = i
					break
				}
			}

			if findIndex == -1 {
				x.Specs = append(x.Specs, &ast.ImportSpec{
					Path: &ast.BasicLit{
						Kind:  token.STRING,
						Value: fmt.Sprintf("\"%s\"", modulePath),
					},
				})
			}
		}
	}
}

func (v *astModuleMainVisitor) modifyStructField(xst *ast.StructType) {
	findIndex := -1
	for i, field := range xst.Fields.List {
		starType, ok := field.Type.(*ast.StarExpr)
		if !ok {
			continue
		}
		selector, ok := starType.X.(*ast.SelectorExpr)
		if !ok {
			continue
		}
		if selector.Sel.Name == v.args.StructName && selector.X.(*ast.Ident).Name == StructPackageAPI {
			findIndex = i
			break
		}
	}

	if v.args.Flag&AstFlagGen != 0 {
		if findIndex != -1 {
			return
		}

		existsAPI := false
		for _, gpkg := range v.args.GenPackages {
			if gpkg == StructPackageAPI {
				existsAPI = true
				break
			}
		}

		if existsAPI {
			xst.Fields.List = append(xst.Fields.List, &ast.Field{
				Names: []*ast.Ident{
					{Name: GetStructAPIName(v.args.StructName)},
				},
				Type: &ast.StarExpr{
					X: &ast.SelectorExpr{
						X:   ast.NewIdent(StructPackageAPI),
						Sel: ast.NewIdent(v.args.StructName),
					},
				},
			})
		}
	} else if v.args.Flag&AstFlagRem != 0 {
		if findIndex != -1 {
			xst.Fields.List = append(xst.Fields.List[:findIndex], xst.Fields.List[findIndex+1:]...)
		}
	}
}

func (v *astModuleMainVisitor) modifyAutoMigrate(x *ast.FuncDecl) {
	if len(x.Body.List) == 0 {
		return
	}

	switch xst := x.Body.List[0].(type) {
	case *ast.ReturnStmt:
		if len(xst.Results) == 0 {
			return
		}

		result := xst.Results[0].(*ast.CallExpr)
		args := result.Args
		findIndex := -1
		for i, arg := range args {
			selector, ok := arg.(*ast.CallExpr).Args[0].(*ast.SelectorExpr)
			if !ok {
				continue
			}

			if selector.Sel.Name == v.args.StructName && selector.X.(*ast.Ident).Name == StructPackageSchema {
				findIndex = i
				break
			}
		}

		if v.args.Flag&AstFlagGen != 0 {
			if findIndex != -1 {
				return
			}
			args = append(args, &ast.CallExpr{
				Fun: ast.NewIdent("new"),
				Args: []ast.Expr{
					&ast.SelectorExpr{
						X:   ast.NewIdent(StructPackageSchema),
						Sel: ast.NewIdent(v.args.StructName),
					},
				},
			})
			result.Args = args
		} else if v.args.Flag&AstFlagRem != 0 {
			if findIndex == -1 {
				return
			}
			args = append(args[:findIndex], args[findIndex+1:]...)
			result.Args = args
		}
	}
}

func (v *astModuleMainVisitor) modifyRegisterV1Routers(x *ast.FuncDecl) {
	if len(x.Body.List) == 0 {
		return
	}

	structRouterVarName := GetStructRouterVarName(v.args.StructName)
	findIndex := -1
	for i, list := range x.Body.List {
		if lt, ok := list.(*ast.AssignStmt); ok {
			if len(lt.Lhs) == 0 {
				continue
			}

			if lt.Lhs[0].(*ast.Ident).Name == structRouterVarName {
				findIndex = i
				break
			}
		}
	}

	if v.args.Flag&AstFlagGen != 0 {
		if findIndex != -1 {
			return
		}

		existsAPI := false
		for _, gpkg := range v.args.GenPackages {
			if gpkg == StructPackageAPI {
				existsAPI = true
				break
			}
		}
		if !existsAPI {
			return
		}

		assignStmt := &ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.Ident{
					Name: structRouterVarName,
					Obj:  ast.NewObj(ast.Var, structRouterVarName),
				},
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("v1"),
						Sel: ast.NewIdent("Group"),
					},
					Args: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf("\"/%s\"", GetStructRouterGroupName(v.args.StructName)),
						},
					},
				},
			},
		}

		routes := [][]string{
			{"GET", "\"\"", "Query"},
			{"GET", "\"/:id\"", "Get"},
			{"POST", "\"\"", "Create"},
			{"PUT", "\"\"", "Update"},
			{"DELETE", "\"/:id\"", "Delete"},
		}

		var blockList []ast.Stmt
		for _, r := range routes {
			blockList = append(blockList, &ast.ExprStmt{
				X: &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent(structRouterVarName),
						Sel: ast.NewIdent(r[0]),
					},
					Args: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.STRING,
							Value: r[1],
						},
						&ast.SelectorExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("a"),
								Sel: ast.NewIdent(GetStructAPIName(v.args.StructName)),
							},
							Sel: ast.NewIdent(r[2]),
						},
					},
				},
			})
		}

		lastEle := x.Body.List[len(x.Body.List)-1]
		x.Body.List = append(x.Body.List[:len(x.Body.List)-1],
			assignStmt,
			&ast.BlockStmt{List: blockList},
			lastEle,
		)
	} else if v.args.Flag&AstFlagRem != 0 {
		if findIndex == -1 {
			return
		}
		x.Body.List = append(x.Body.List[:findIndex], x.Body.List[findIndex+2:]...)
	}
}

type astModuleWireVisitor struct {
	fset *token.FileSet
	args BasicArgs
}

func (v *astModuleWireVisitor) Visit(node ast.Node) ast.Visitor {
	switch x := node.(type) {
	case *ast.GenDecl:
		if x.Tok == token.IMPORT && len(x.Specs) > 0 {
			v.modifyModuleImport(x)
		}

		if x.Tok == token.VAR && len(x.Specs) > 0 {
			v.modifyNewSet(x)
		}
	}
	return v
}

func (v *astModuleWireVisitor) modifyModuleImport(x *ast.GenDecl) {
	if v.args.Flag&AstFlagGen != 0 {
		for _, pkgName := range []string{StructPackageAPI, StructPackageBIZ, StructPackageDAL} {
			findIndex := -1
			modulePath := GetModuleImportPath(v.args.Dir, v.args.ModulePath, v.args.ModuleName) + "/" + pkgName
			for i, spec := range x.Specs {
				if is, ok := spec.(*ast.ImportSpec); ok &&
					is.Path.Value == fmt.Sprintf("\"%s\"", modulePath) {
					findIndex = i
					break
				}
			}

			if findIndex == -1 {
				x.Specs = append(x.Specs, &ast.ImportSpec{
					Path: &ast.BasicLit{
						Kind:  token.STRING,
						Value: fmt.Sprintf("\"%s\"", modulePath),
					},
				})
			}
		}
	}
}

func (v *astModuleWireVisitor) modifyNewSet(x *ast.GenDecl) {
	vspec, ok := x.Specs[0].(*ast.ValueSpec)
	if !ok || len(vspec.Names) == 0 ||
		len(vspec.Values) == 0 || vspec.Names[0].Name != "Set" {
		return
	}

	args := vspec.Values[0].(*ast.CallExpr).Args

	if v.args.Flag&AstFlagGen != 0 {
		genPackagesMap := make(map[string]bool)
		for _, p := range v.args.GenPackages {
			if p == StructPackageSchema {
				continue
			}
			genPackagesMap[p] = true
		}

		for _, arg := range args {
			if wireS, ok := arg.(*ast.CallExpr); ok && len(wireS.Args) > 0 {
				if newS, ok := wireS.Args[0].(*ast.CallExpr); ok && len(newS.Args) > 0 {
					if s, ok := newS.Args[0].(*ast.SelectorExpr); ok {
						if s.Sel.Name == v.args.StructName {
							name := s.X.(*ast.Ident).Name
							if _, ok := genPackagesMap[name]; ok {
								delete(genPackagesMap, name)
								continue
							}
						}
					}
				}
			}
		}

		for _, p := range v.args.GenPackages {
			if _, ok := genPackagesMap[p]; !ok {
				continue
			}

			arg := &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent("wire"),
					Sel: ast.NewIdent("Struct"),
				},
				Args: []ast.Expr{
					&ast.CallExpr{
						Fun: ast.NewIdent("new"),
						Args: []ast.Expr{
							&ast.SelectorExpr{
								X:   ast.NewIdent(p),
								Sel: ast.NewIdent(v.args.StructName),
							},
						},
					},
					&ast.BasicLit{
						Kind:  token.STRING,
						Value: "\"*\"",
					},
				},
			}
			args = append(args, arg)
		}
		vspec.Values[0].(*ast.CallExpr).Args = args
	} else if v.args.Flag&AstFlagRem != 0 {
		var newArgs []ast.Expr
		for _, arg := range args {
			if wireS, ok := arg.(*ast.CallExpr); ok && len(wireS.Args) > 0 {
				if newS, ok := wireS.Args[0].(*ast.CallExpr); ok && len(newS.Args) > 0 {
					if s, ok := newS.Args[0].(*ast.SelectorExpr); ok {
						if s.Sel.Name == v.args.StructName {
							continue
						}
					}
				}
			}
			newArgs = append(newArgs, arg)
		}
		vspec.Values[0].(*ast.CallExpr).Args = newArgs
	}
}

type astModsVisitor struct {
	fset *token.FileSet
	args BasicArgs
}

func (v *astModsVisitor) Visit(node ast.Node) ast.Visitor {
	switch x := node.(type) {
	case *ast.GenDecl:
		if x.Tok == token.IMPORT && len(x.Specs) > 0 {
			v.modifyModuleImport(x)
		}

		if x.Tok == token.VAR && len(x.Specs) > 0 {
			v.modifyWireSet(x)
		}
	case *ast.TypeSpec:
		if x.Name.Name == "Mods" {
			if xst, ok := x.Type.(*ast.StructType); ok {
				v.modifyStructField(xst)
			}
		}
	case *ast.FuncDecl:
		if x.Name.Name == "Init" {
			v.modifyFuncInit(x)
		} else if x.Name.Name == "RegisterRouters" {
			v.modifyFuncRegisterRouters(x)
		}
	}
	return v
}

func (v *astModsVisitor) modifyModuleImport(x *ast.GenDecl) {
	findIndex := -1
	modulePath := GetModuleImportPath(v.args.Dir, v.args.ModulePath, v.args.ModuleName)
	for i, spec := range x.Specs {
		if is, ok := spec.(*ast.ImportSpec); ok &&
			is.Path.Value == fmt.Sprintf("\"%s\"", modulePath) {
			findIndex = i
			break
		}
	}

	if v.args.Flag&AstFlagGen != 0 {
		if findIndex == -1 {
			x.Specs = append(x.Specs, &ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf("\"%s\"", modulePath),
				},
			})
		}
	} else if v.args.Flag&AstFlagRem != 0 {
		if findIndex != -1 {
			x.Specs = append(x.Specs[:findIndex], x.Specs[findIndex+1:]...)
		}
	}
}

func (v *astModsVisitor) modifyWireSet(x *ast.GenDecl) {
	vspec, ok := x.Specs[0].(*ast.ValueSpec)
	if !ok || len(vspec.Names) == 0 ||
		len(vspec.Values) == 0 || vspec.Names[0].Name != "Set" {
		return
	}

	args := vspec.Values[0].(*ast.CallExpr).Args
	findIndex := -1
	for i, arg := range args {
		if sel, ok := arg.(*ast.SelectorExpr); ok && sel.X.(*ast.Ident).Name == GetModuleImportName(v.args.ModuleName) {
			findIndex = i
			break
		}
	}

	if v.args.Flag&AstFlagGen != 0 {
		if findIndex != -1 {
			return
		}
		arg := &ast.SelectorExpr{
			X:   ast.NewIdent(GetModuleImportName(v.args.ModuleName)),
			Sel: ast.NewIdent("Set"),
		}
		args = append(args, arg)
		vspec.Values[0].(*ast.CallExpr).Args = args
	} else if v.args.Flag&AstFlagRem != 0 {
		if findIndex == -1 {
			return
		}
		args = append(args[:findIndex], args[findIndex+1:]...)
		vspec.Values[0].(*ast.CallExpr).Args = args
	}
}

func (v *astModsVisitor) modifyStructField(xst *ast.StructType) {
	findIndex := -1
	for i, field := range xst.Fields.List {
		starType, ok := field.Type.(*ast.StarExpr)
		if !ok {
			continue
		}
		selector, ok := starType.X.(*ast.SelectorExpr)
		if !ok {
			continue
		}
		if selector.Sel.Name == v.args.ModuleName {
			findIndex = i
			break
		}
	}

	if v.args.Flag&AstFlagGen != 0 {
		if findIndex != -1 {
			return
		}
		xst.Fields.List = append(xst.Fields.List, &ast.Field{
			Names: []*ast.Ident{
				{Name: v.args.ModuleName},
			},
			Type: &ast.StarExpr{
				X: &ast.SelectorExpr{
					X:   ast.NewIdent(GetModuleImportName(v.args.ModuleName)),
					Sel: ast.NewIdent(v.args.ModuleName),
				},
			},
		})
	} else if v.args.Flag&AstFlagRem != 0 {
		if findIndex != -1 {
			xst.Fields.List = append(xst.Fields.List[:findIndex], xst.Fields.List[findIndex+1:]...)
		}
	}
}

func (v *astModsVisitor) modifyFuncInit(x *ast.FuncDecl) {
	findIndex := -1
	list := x.Body.List
	for i, stmt := range list {
		if s, ok := stmt.(*ast.IfStmt); ok {
			var sb strings.Builder
			_ = printer.Fprint(&sb, v.fset, s.Init)
			if strings.Contains(sb.String(), fmt.Sprintf("%s.Init", v.args.ModuleName)) {
				findIndex = i
				break
			}
		}
	}
	if v.args.Flag&AstFlagGen != 0 {
		if findIndex == -1 {
			e, err := parser.ParseExpr(fmt.Sprintf("a.%s.Init(ctx)", v.args.ModuleName))
			if err == nil {
				list = append(list[:len(list)-1], append([]ast.Stmt{&ast.IfStmt{
					Init: &ast.AssignStmt{
						Lhs: []ast.Expr{
							ast.NewIdent("err"),
						},
						Tok: token.DEFINE,
						Rhs: []ast.Expr{
							e,
						},
					},
					Cond: &ast.BinaryExpr{
						X:  ast.NewIdent("err"),
						Op: token.NEQ,
						Y:  ast.NewIdent("nil"),
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ReturnStmt{
								Results: []ast.Expr{
									ast.NewIdent("err"),
								},
							},
						},
					},
				}}, list[len(list)-1])...)
				x.Body.List = list
			}
		}
	} else if v.args.Flag&AstFlagRem != 0 {
		if findIndex != -1 {
			list = append(list[:findIndex], list[findIndex+1:]...)
			x.Body.List = list
		}
	}
}

func (v *astModsVisitor) modifyFuncRegisterRouters(x *ast.FuncDecl) {
	findIndex := -1
	list := x.Body.List
	for i, stmt := range list {
		if s, ok := stmt.(*ast.IfStmt); ok {
			var sb strings.Builder
			printer.Fprint(&sb, v.fset, s.Init)
			if strings.Contains(sb.String(), fmt.Sprintf("%s.RegisterV1Routers", v.args.ModuleName)) {
				findIndex = i
				break
			}
		}
	}
	if v.args.Flag&AstFlagGen != 0 {
		if findIndex == -1 {
			e, err := parser.ParseExpr(fmt.Sprintf("a.%s.RegisterV1Routers(ctx, v1)", v.args.ModuleName))
			if err == nil {
				list = append(list[:len(list)-1], append([]ast.Stmt{&ast.IfStmt{
					Init: &ast.AssignStmt{
						Lhs: []ast.Expr{
							ast.NewIdent("err"),
						},
						Tok: token.DEFINE,
						Rhs: []ast.Expr{
							e,
						},
					},
					Cond: &ast.BinaryExpr{
						X:  ast.NewIdent("err"),
						Op: token.NEQ,
						Y:  ast.NewIdent("nil"),
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ReturnStmt{
								Results: []ast.Expr{
									ast.NewIdent("err"),
								},
							},
						},
					},
				}}, list[len(list)-1])...)
				x.Body.List = list
			}
		}
	} else if v.args.Flag&AstFlagRem != 0 {
		if findIndex != -1 {
			list = append(list[:findIndex], list[findIndex+1:]...)
			x.Body.List = list
		}
	}
}
