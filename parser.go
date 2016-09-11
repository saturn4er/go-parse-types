package tparser

// Package main provides
import (
	"errors"
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
)

type packageParse struct {
	gopath       string
	path         string
	importPath   string
	imports      []string
	importsTypes map[string]map[string]*Type // Map of package => type name=> *Type
	Types        []*Type
	nameTypes    map[string]*Type
	fset         *token.FileSet
	pkg          *ast.Package
	doc          *doc.Package
}

func (pp *packageParse) GetType(name string) (*Type, error) {
	if v, ok := pp.nameTypes[name]; ok {
		return v, nil
	}
	return nil, errors.New("No such type")
}
func (pp *packageParse) parse() error {
	pp.nameTypes = make(map[string]*Type)
	pp.Types = []*Type{}

	pp.fset = token.NewFileSet() // positions are relative to fset
	d, err := parser.ParseDir(pp.fset, pp.path, nil, 0)
	if err != nil {
		return err
	}
	if len(d) > 1 {
		return errors.New("Multiple packages in folder")
	}
	for _, f := range d {
		pp.pkg = f
	}
	pp.doc = doc.New(pp.pkg, pp.path, 0)
	pp.imports = pp.doc.Imports
	pp.importsTypes = make(map[string]map[string]*Type)
	for _, i := range pp.imports {
		ip, err := New(filepath.Join(pp.gopath, "src", i))
		if err != nil {
			continue
		}
		pp.importsTypes[ip.pkg.Name] = make(map[string]*Type)
		for _, ty := range ip.Types {
			pp.importsTypes[ip.pkg.Name][ty.Name] = ty
		}
	}
	pp.parsePackage()
	return nil
}
func (pp *packageParse) getTypeByName(n string) *Type {
	if _, ok := pp.nameTypes[n]; !ok {
		ty := &Type{
			Name:    n,
			PkgName: pp.pkg.Name,
			PkgPath: pp.importPath,
			Kind:    Invalid,
		}
		pp.Types = append(pp.Types, ty)
		pp.nameTypes[n] = ty
	}
	return pp.nameTypes[n]
}
func (pp *packageParse) parsePackage() {
	ast.Inspect(pp.pkg, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			pp.parseFunction(x)
		case *ast.TypeSpec:
			pp.parseTypeSpec(x)
		}
		return true
	})
	for _, dt := range pp.doc.Types {
		ty := pp.getTypeByName(dt.Name)
		ty.Doc = dt.Doc
		for _, f := range dt.Methods {
			for _, m := range ty.Methods {
				if m.Name == f.Name {
					m.Doc = strings.TrimSpace(f.Doc)
				}
			}
		}
	}
}
func (pp *packageParse) parseFunction(f *ast.FuncDecl) {
	var name string
	var t *Type
	if f.Recv == nil || len(f.Recv.List) == 0 {
		t = pp.getTypeByName(f.Name.Name)
		t.Kind = Func
	} else {
		rt := pp.getTypeByName(getExprName(f.Recv.List[0].Type))
		t = &Type{}
		rt.Methods = append(rt.Methods, t)
	}
	if f.Name != nil {
		name = f.Name.Name
	}
	t.Name = name
	t.Kind = Func
	t.In = pp.parseFunctionParams(f.Type.Params)
	t.Out = pp.parseFunctionParams(f.Type.Results)
}
func (pp *packageParse) parseTypeSpec(t *ast.TypeSpec) {
	if t.Name == nil {
		return
	}
	ty := pp.getTypeByName(t.Name.Name)
	pp.parseType(t.Type, ty)
}
func (pp *packageParse) parseType(e ast.Expr, to *Type) *Type {
	result := &Type{}
	if to != nil {
		result = to
	}
	switch ty := e.(type) {
	case *ast.StarExpr:
		result.Kind = Ptr
		result.ElementT = pp.parseType(ty.X, nil)
	case *ast.SelectorExpr:
		packageName := getExprName(ty.X)
		if v, ok := pp.importsTypes[packageName]; ok {
			if vn, ok := v[ty.Sel.Name]; ok {
				return vn
			}
		}
	case *ast.ArrayType:
		result.ElementT = pp.parseType(ty.Elt, nil)
		result.Kind = Array
	case *ast.StructType:
		result.Kind = Struct
		result.Fields = pp.parseStructureFields(ty)
	case *ast.ChanType:
		result.Kind = Chan
	case *ast.FuncType:
		result.Kind = Func
	case *ast.InterfaceType:
		result.Kind = Interface
	case *ast.MapType:
		result.Kind = Map
		result.KeyT = pp.parseType(ty.Key, nil)
		result.ElementT = pp.parseType(ty.Value, nil)
	case *ast.Ident:
		switch ty.Name {
		case "bool":
			result.Kind = Bool
		case "int":
			result.Kind = Int
		case "int8":
			result.Kind = Int8
		case "int16":
			result.Kind = Int16
		case "int32":
			result.Kind = Int32
		case "int64":
			result.Kind = Int64
		case "uint":
			result.Kind = Uint
		case "uint8":
			result.Kind = Uint8
		case "uint16":
			result.Kind = Uint16
		case "uint32":
			result.Kind = Uint32
		case "uint64":
			result.Kind = Uint64
		case "uintptr":
			result.Kind = Uintptr
		case "float32":
			result.Kind = Float32
		case "float64":
			result.Kind = Float64
		case "complex64":
			result.Kind = Complex64
		case "complex128":
			result.Kind = Complex128
		default:
			return pp.getTypeByName(ty.Name)
		}
	default:
		fmt.Printf("Parsing type %T", e)
	}
	return result
}
func (pp *packageParse) parseFunctionParams(f *ast.FieldList) []*FunctionParameter {
	result := []*FunctionParameter{}
	if f == nil {
		return result
	}
	for _, arg := range f.List {
		ty := pp.parseType(arg.Type, nil)
		if len(arg.Names) == 0 {
			result = append(result, &FunctionParameter{
				Type: ty,
			})
			continue
		}
		for _, n := range arg.Names {
			result = append(result, &FunctionParameter{
				Type: ty,
				Name: n.Name,
			})
		}
	}
	return result
}
func (pp *packageParse) parseStructureFields(s *ast.StructType) []*Field {
	if s.Fields == nil {
		return nil
	}
	result := []*Field{}
	for _, f := range s.Fields.List {
		t := pp.parseType(f.Type, nil)
		if len(f.Names) == 0 {
			result = append(result, &Field{
				Type:      t,
				Anonymous: true,
			})
			continue
		}
		for _, n := range f.Names {
			if n == nil {
				continue
			}
			result = append(result, &Field{
				Type: t,
				Name: n.Name,
			})
		}
	}
	return result
}
func getExprName(e ast.Expr) string {
	switch ty := e.(type) {
	case *ast.StarExpr:
		return getExprName(ty.X)
	case *ast.Ident:
		return ty.Name
	case *ast.SelectorExpr:
		return getExprName(ty.X) + "." + ty.Sel.Name
	default:
		fmt.Printf("getExprName of unknown type %T", e)
	}
	return ""
}

func New(path string) (*packageParse, error) {
	gopath := GetGOPATH()
	if len(gopath) == 0 {
		return nil, errors.New("GOPATH is empty")
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	importPath, err := filepath.Rel(filepath.Join(gopath, "src"), absPath)
	if err != nil {
		return nil, err
	}
	p := &packageParse{
		importPath: importPath,
		gopath:     gopath,
		path:       path,
	}
	err = p.parse()
	return p, err
}
