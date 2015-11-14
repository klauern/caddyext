package directives

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
)

var (
	ErrInvalidDirectiveFile         = errors.New("Invalid directive file")
	ErrImportInvalidDirectiveFormat = errors.New("Directive file is importing using an invalid format")
	ErrDirectiveAlreadyImported     = errors.New("Directive already imported")
	ErrDirectiveNotFound            = errors.New("Directive not found")
	ErrDirectiveInvalidCore         = errors.New("Invalid directive (from the core)")
	ErrStackInvalidIndex            = errors.New("Invalid index")
)

// Load ...
func NewFrom(file string) (Manager, error) {
	if _, err := os.Stat(file); err != nil {
		return nil, err
	}

	dir := &DirectiveList{
		file: file,
		list: make([]*Directive, 0),
	}
	if err := dir.LoadList(); err != nil {
		return nil, err
	}

	return dir, nil
}

// Manager ...
type Manager interface {
	LoadList() error
	List() []Directive
	AddDirective(string, string) error
	RemoveDirective(string) error
	EnableDirective(string) error
	DisableDirective(string) error
	MoveDirective(string, int) error
	Save() error
}

// DirectiveList ...
type DirectiveList struct {
	file string
	list []*Directive
}

// Directive ...
type Directive struct {
	Name       string
	Setup      string
	ImportPath string
	Active     bool
	Core       bool
	Removed    bool
}

func (d *DirectiveList) List() []Directive {
	list := make([]Directive, len(d.list))
	for i, d := range d.list {
		list[i] = *d
	}
	return list
}

func (d *DirectiveList) AddDirective(name string, source string) error {
	for _, directive := range d.list {
		if directive.Name == name {
			return ErrDirectiveAlreadyImported
		}
	}

	dir := &Directive{
		Name:       name,
		ImportPath: source,
		Active:     true,
		Core:       false,
		Setup:      fmt.Sprintf("%s.Setup", name),
		Removed:    false,
	}
	d.list = append(d.list, dir)
	return nil
}

func (d *DirectiveList) RemoveDirective(name string) error {
	for _, directive := range d.list {
		if directive.Name == name {
			if directive.Core {
				return ErrDirectiveInvalidCore
			}
			directive.Removed = true
			return nil
		}
	}
	return ErrDirectiveNotFound
}

func (d *DirectiveList) EnableDirective(name string) error {
	for _, directive := range d.list {
		if directive.Name == name {
			directive.Active = true
			return nil
		}
	}
	return ErrDirectiveNotFound
}

func (d *DirectiveList) DisableDirective(name string) error {
	for _, directive := range d.list {
		if directive.Name == name {
			directive.Active = false
			return nil
		}
	}
	return ErrDirectiveNotFound
}

func (d *DirectiveList) MoveDirective(name string, index int) error {
	if index < 0 || index >= len(d.list) {
		return ErrStackInvalidIndex
	}

	actual := -1
	for i, dir := range d.list {
		if dir.Name == name {
			actual = i
		}
	}

	if actual == -1 {
		return ErrDirectiveNotFound
	}

	dir := d.list[actual]
	d.list = append(d.list[:actual], d.list[actual+1:]...)
	d.list = append(d.list[:index], append([]*Directive{dir}, d.list[index:]...)...)

	return nil
}

var regComment *regexp.Regexp = regexp.MustCompile("//@caddyext ")
var regCommentLine *regexp.Regexp = regexp.MustCompile(`//@caddyext[\s{"|a-zA-Z]+`)

func (d *DirectiveList) LoadList() error {
	data, err := ioutil.ReadFile(d.file)
	if err != nil {
		return err
	}

	fset := token.NewFileSet()

	extMatches := regCommentLine.FindAll(data, -1)
	data = regComment.ReplaceAll(data, []byte(""))

	f, err := parser.ParseFile(token.NewFileSet(), "", data, 0)
	impPaths := astutil.Imports(fset, f)
	node, ok := f.Scope.Lookup("directiveOrder").Decl.(ast.Node)
	if !ok {
		return ErrInvalidDirectiveFile
	}

	c := node.(*ast.ValueSpec).Values[0].(*ast.CompositeLit)
	for _, m := range c.Elts {
		var setup *ast.SelectorExpr
		var directive *ast.BasicLit
		token := m.(*ast.CompositeLit).Elts
		active := true

		if v, ok := token[0].(*ast.BasicLit); ok {
			directive = v
		} else {
			return ErrImportInvalidDirectiveFormat
		}

		if v, ok := m.(*ast.CompositeLit).Elts[1].(*ast.SelectorExpr); ok {
			setup = v
		} else if _, ok := m.(*ast.CompositeLit).Elts[1].(*ast.FuncLit); ok {
			active = false
		} else {
			return ErrImportInvalidDirectiveFormat
		}

		name, err := strconv.Unquote(directive.Value)
		if err != nil {
			return ErrImportInvalidDirectiveFormat
		}
		dir := &Directive{
			Name:    name,
			Active:  active,
			Core:    true,
			Removed: false,
		}

		if active {
			dir.Setup = fmt.Sprintf("%s.%s", setup.X, setup.Sel)
		} else {
			dir.Name = strings.TrimLeft(dir.Name, "!")
		}

		for _, imp := range impPaths[0] {
			if imp.Name.String() == dir.Name {
				dir.Core = false
				dir.ImportPath, _ = strconv.Unquote(imp.Path.Value)
				break
			}
		}

		for _, commented := range extMatches {
			if strings.Contains(string(commented), `{"`+dir.Name) {
				dir.Active = false
				break
			}
		}

		d.list = append(d.list, dir)
	}

	return nil
}

var dirOrderStart string = "directiveOrder = []directive{\n"
var dirOrderEnd string = "}\n"

func (d *DirectiveList) Save() error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, d.file, nil, 0)
	if err != nil {
		return ErrInvalidDirectiveFile
	}

	imps := astutil.Imports(fset, f)
	for _, imp := range imps[0] {
		path, _ := strconv.Unquote(imp.Path.Value)
		astutil.DeleteImport(fset, f, path)
	}

	astutil.AddImport(fset, f, "github.com/mholt/caddy/caddy/parse")
	astutil.AddImport(fset, f, "github.com/mholt/caddy/caddy/setup")
	astutil.AddImport(fset, f, "github.com/mholt/caddy/middleware")

	for _, dir := range d.list {
		if dir.Removed {
			continue
		}
		if dir.Core == false && len(dir.ImportPath) != 0 {
			astutil.AddImport(fset, f, fmt.Sprintf("{{import-%s}}", dir.Name))
		}
	}

	var buf bytes.Buffer
	err = printer.Fprint(&buf, fset, f)
	if err != nil {
		return err
	}

	out := buf.String()

	f, err = parser.ParseFile(token.NewFileSet(), "", out, 0)
	node, ok := f.Scope.Lookup("directiveOrder").Decl.(ast.Node)
	if !ok {
		return ErrInvalidDirectiveFile
	}

	begin := out[0 : node.Pos()-1]
	end := out[node.End():len(out)]

	dirOrder := ""
	for _, dir := range d.list {
		if dir.Removed {
			continue
		}

		comment := ""
		if dir.Active == false {
			comment = "//@caddyext "
		}

		begin = strings.Replace(begin, fmt.Sprintf(`"{{import-%s}}"`, dir.Name), fmt.Sprintf(`%s%s "%s"`, comment, dir.Name, dir.ImportPath), -1)

		if dir.Core == true {
			dirOrder = dirOrder + fmt.Sprintf(`	%s{"%s", %s},`+"\n", comment, dir.Name, dir.Setup)
		} else {
			dirOrder = dirOrder + fmt.Sprintf(`	%s{"%s", %s.Setup},`+"\n", comment, dir.Name, dir.Name)
		}
	}

	out = begin + dirOrderStart + dirOrder + dirOrderEnd + end

	return ioutil.WriteFile(d.file, []byte(out), os.FileMode(0660))
}
