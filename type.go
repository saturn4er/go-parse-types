package tparser

import (
	"errors"
	"fmt"
)

type Types []Type

// Type represent golang type
type Type struct {
	PkgPath  string
	PkgName  string
	Name     string
	KeyT     *Type // KeyType
	ElementT *Type // ElementType
	Kind     Kind
	Doc      string
	Fields   []*Field
	Methods  []*Type
	In       []*FunctionParameter
	Out      []*FunctionParameter
	Length   int
}

func (t *Type) String() string {
	return t.StringO("")
}
func (t *Type) StringO(prefix string) string {
	if t == nil {
		return "nil\n"
	}
	result := ""
	if t.Name != "" {
		result += fmt.Sprintf(prefix+"Name: %s\n", t.Name)
	}
	if t.PkgName != "" {
		result += fmt.Sprintf(prefix+"Package: %s (%s)\n", t.PkgName, t.PkgPath)
	}
	result += fmt.Sprintf(prefix+"Type: %s\n", t.Kind)
	switch t.Kind {
	case Array:
		result += fmt.Sprintf(prefix+"Element kind: %s \n", t.ElementT)
	case Func:
		result += prefix + "Args:\n"
		for _, a := range t.In {
			result += prefix + "   Name: " + a.Name + "\n"
			result += prefix + "   Type:\n" + a.Type.StringO(prefix+"      ")
		}
	case Ptr:
		result += t.ElementT.StringO(prefix + "   ")
	case Struct:
		result += prefix + "Fields:\n"
		for _, f := range t.Fields {
			result += prefix + "   Name: " + f.Name + "\n"
			result += prefix + "   Type:\n" + f.Type.StringO(prefix+"      ")
		}
		result += prefix + "Methods:\n"
		for _, m := range t.Methods {
			result += prefix + "   Name: " + m.Name + "\n"
			result += prefix + "   Type:\n" + m.StringO(prefix+"      ")
		}
	}
	return result
}
func (t *Type) FindFieldWithType(packageName string, packagePath string, typeName string) (bool, string, error) {
	if t.Fields == nil {
		return false, "", errors.New("Not a struct type")
	}
	for _, value := range t.Fields {
		if value.Type.PkgName == packageName && value.Type.PkgPath == packagePath && typeName == value.Type.Name {
			return true, value.Path(), nil
		}
		if value.Type.Kind == Struct {
			b, p, err := value.Type.FindFieldWithType(packageName, packagePath, typeName)
			if err != nil {
				continue
			}
			return b, value.Path()+ p, nil
		}
		if value.Type.Kind == Ptr {
			pt := value.Type.ElementT
			if pt.PkgName == packageName && pt.PkgPath == packagePath && typeName == pt.Name {
				return true, "." + value.Path(), nil

			}
		}
	}
	return false, "", errors.New("Can't find field with such type")
}

type FunctionParameter struct {
	Name string
	Type *Type
}

type Field struct {
	Name      string
	Type      *Type
	Anonymous bool
}

func (f *Field) Path() string {
	if f.Anonymous {
		return "." + f.Type.Name
	}
	return "." + f.Name
}
