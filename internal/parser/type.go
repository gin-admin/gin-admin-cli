package parser

type AstFlag int

func (f AstFlag) String() string {
	switch f {
	case AstFlagGen:
		return "G"
	case AstFlagRem:
		return "R"
	}
	return "?"
}

const (
	AstFlagGen AstFlag = 1 << iota
	AstFlagRem
)

type BasicArgs struct {
	Dir              string
	ModuleName       string
	ModulePath       string
	StructName       string
	GenPackages      []string
	Flag             AstFlag
	FillRouterPrefix bool
}
