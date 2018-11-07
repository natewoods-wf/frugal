/*
 * Copyright 2017 Workiva
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *     http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package gopherjs

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"unicode"

	"golang.org/x/tools/imports"

	"github.com/Workiva/frugal/compiler/generator"
	"github.com/Workiva/frugal/compiler/globals"
	"github.com/Workiva/frugal/compiler/parser"
)

const (
	lang                = "golang"
	fileSuffix          = "go"
	defaultOutputDir    = "gen-go"
	serviceSuffix       = "_service"
	scopeSuffix         = "_scope"
	packagePrefixOption = "package_prefix"
	frugalImportOption  = "frugal_import"
)

// Generator implements the LanguageGenerator interface for Go.
type Generator struct {
	*generator.BaseGenerator
	typesFile  *os.File
	structTPL  *template.Template
	enumTPL    *template.Template
	serviceTPL *template.Template
}

// NewGenerator creates a new Go LanguageGenerator.
func NewGenerator(options map[string]string) generator.LanguageGenerator {
	return &Generator{
		BaseGenerator: &generator.BaseGenerator{Options: options},
	}
}

var enumTemplate = `
// {{title .Name}} is an enum.
type {{title .Name}} int64

// {{title .Name}} values.
const (
	{{range .Values -}}
	{{title $.Name}}{{.Name}} {{title $.Name}} = {{.Value}}
	{{end -}}
)
`

const structTemplate = `
// {{title .Name}} is a frual serializable object.
type {{title .Name}} struct {
	{{range .Fields -}}
	{{title .Name}} {{fieldType .}}
	{{end}}
}

// New{{title .Name}} constructs a {{title .Name}}.
func New{{title .Name}}() *{{title .Name}} {
	return &{{title .Name}}{
		// TODO: default values
		{{/*
		func (g *Generator) generateConstructor(s *parser.Struct, sName string) string {
			contents := ""

			contents += fmt.Sprintf("func New%s() *%s {\n", sName, sName)
			contents += fmt.Sprintf("\treturn &%s{\n", sName)

			for _, field := range s.Fields {
				// Use the default if it exists and it's not a pointer field, otherwise the zero value is implicitly used
				if field.Default != nil && !g.isPointerField(field) {
					val := g.generateConstantValue(field.Type, field.Default)
					contents += fmt.Sprintf("\t\t%s: %s,\n", title(field.Name), val)
				}
			}

			contents += "\t}\n"
			contents += "}\n\n"
			return contents
		}
		*/}}
	}
}

// Unpack deserializes {{.Name}} objects.
func (p *{{title .Name}}) Unpack(prot frugal.Protocol) {
  prot.UnpackStructBegin("{{.Name}}")
  for typeID, id := prot.UnpackFieldBegin(); typeID != frugal.STOP; typeID, id = prot.UnpackFieldBegin() {
    switch id {
    {{range .Fields -}}
    case {{.ID}}:
      {{unpackField . -}}
    {{end -}}
    default:
      prot.Skip(typeID)
    }
    prot.UnpackFieldEnd()
  }
  prot.UnpackStructEnd()
}

// Pack serializes {{.Name}} objects.
func (p *{{title .Name}}) Pack(prot frugal.Protocol) {
	{{if eq "union" .Type.String -}}
	count := 0
	{{range .Fields -}}
		if p.{{title .Name}} != nil {
			count++
		}
	{{end -}}
	if count != 1 {
		prot.Set(errors.New("{{.Name}} invalid union state"))
		return
	}
	{{end -}}
  prot.PackStructBegin("{{.Name}}")
	{{range .Fields}}{{packField .}}{{end -}}
	prot.PackFieldStop()
  prot.PackStructEnd()
}

{{if eq "exception" .Type.String -}}
func (p *{{title .Name}}) Error() string {
	return "TODO: generate errorz"
}
{{end -}}
`

const serviceTemplate = `
{{define "args"}}{{range .}}, {{.Name}} {{go .Type}}{{end}}{{end}}
{{define "res"}}{{if .}}(r {{go .}}, err error){{else}}(err error){{end}}{{end}}
{{define "func"}}{{title .Name}}(ctx frugal.Context{{template "args" .Arguments}}) {{template "res" .ReturnType}}{{end}}

// {{title .Name}} is a service or a client.
type {{title .Name}} interface {
	{{range .Methods -}}
	{{template "func" .}}
	{{end -}}
}

// {{title .Name}}Client is the client.
type {{title .Name}}Client struct {
	call frugal.CallFunc
}

// New{{title .Name}}Client constructs a {{.Name}}Client.
func New{{title .Name}}Client(cf frugal.CallFunc) *{{.Name}}Client {
	return &{{.Name}}Client{
		call: cf,
	}
}

{{range .Methods -}}
// {{title .Name}} calls a server.
func (c *{{title $.Name}}Client) {{template "func" .}} {
	args := &{{$.Name}}{{title .Name}}Args{
		{{range .Arguments -}}
		{{title .Name}}: {{.Name}},
		{{end -}}
	}
	{{if .Oneway -}}
		return c.call(ctx, "{{lower $.Name}}", "{{lower .Name}}", args, nil)
	{{else -}}
		res := &{{$.Name}}{{title .Name}}Result{}
		err = c.call(ctx, "{{lower $.Name}}", "{{lower .Name}}", args, res)
		if err != nil {
			return
		}
		{{range .Exceptions -}}
		if err = res.{{title .Name}}; err != nil {
			return
		}
		{{end -}}
		{{if .ReturnType -}}
			return res.Success, nil
		{{else -}}
			return nil
		{{end -}}
	{{end -}}
}
{{end -}}

// {{title .Name}}Processor is the client.
type {{title .Name}}Processor struct {
	handler {{title .Name}}
}

// New{{title .Name}}Processor constructs a {{.Name}}Processor.
func New{{title .Name}}Processor(handler {{title .Name}}) *{{.Name}}Processor {
	return &{{.Name}}Processor{
		handler: handler,
	}
}

// Invoke handles internal processing of {{title $.Name}} invocations.
func (p *{{title .Name}}Processor) Invoke(ctx frugal.Context, method string, in frugal.Protocol) (frugal.Packer, error) {
	switch method {
	{{range .Methods -}}
	case "{{lower .Name}}":
		args := &{{$.Name}}{{title .Name}}Args{}
		args.Unpack(in)
		err := in.Err()
		if err != nil {
			return nil, err
		}
		res := &{{$.Name}}{{title .Name}}Result{}
		res.Success, err = p.handler.{{title .Name}}(ctx{{range .Arguments}}, args.{{title .Name}}{{end}})
		switch terr := err.(type) {
		{{range .Exceptions -}}
	case {{go .Type}}:
			res.{{title .Name}} = terr
			err = nil
		{{end -}}
		}
		return res, err
	{{end -}}
	default:
		return nil, errors.New("{{.Name}}: unsupported method " + method)
	}
}
`

// SetupGenerator initializes globals the generator needs, like the types file.
func (g *Generator) SetupGenerator(outputDir string) error {
	t, err := g.GenerateFile("", outputDir, generator.TypeFile)
	if err != nil {
		return err
	}
	g.typesFile = t
	if err = g.GenerateDocStringComment(g.typesFile); err != nil {
		return err
	}
	if err = g.GenerateNewline(g.typesFile, 2); err != nil {
		return err
	}
	if err = g.generatePackage(g.typesFile); err != nil {
		return err
	}
	if err = g.GenerateNewline(g.typesFile, 2); err != nil {
		return err
	}
	if err = g.GenerateTypesImports(g.typesFile); err != nil {
		return err
	}
	if err = g.GenerateNewline(g.typesFile, 2); err != nil {
		return err
	}

	// Build out templates!
	var funcMap = template.FuncMap{
		"title": title,
		"go":    g.getGoTypeFromThriftType,
		"lower": parser.LowercaseFirstLetter,
		"unpackField": func(field *parser.Field) string {
			return g.generateReadFieldRec(field, true)
		},
		"fieldType": func(field *parser.Field) string {
			return g.getGoTypeFromThriftTypePtr(field.Type, g.isPointerField(field))
		},
		"packField": func(field *parser.Field) string {
			return g.generateWriteFieldRec(field, "p.")
		},
	}
	g.structTPL, err = template.New("struct").Funcs(funcMap).Parse(structTemplate)
	if err != nil {
		return err
	}
	g.enumTPL, err = template.New("enum").Funcs(funcMap).Parse(enumTemplate)
	if err != nil {
		return err
	}
	g.serviceTPL, err = template.New("service").Funcs(funcMap).Parse(serviceTemplate)
	if err != nil {
		return err
	}
	return nil
}

// TeardownGenerator cleanups globals the generator needs, like the types file.
func (g *Generator) TeardownGenerator() error {
	if err := g.typesFile.Close(); err != nil { // write pending changes to disk
		return err
	}
	return g.PostProcess(g.typesFile)
}

// GetOutputDir returns the output directory for generated files.
func (g *Generator) GetOutputDir(dir string) string {
	if namespace := g.Frugal.Namespace(lang); namespace != nil {
		path := generator.GetPackageComponents(namespace.Value)
		dir = filepath.Join(append([]string{dir}, path...)...)
	} else {
		dir = filepath.Join(dir, g.Frugal.Name)
	}
	return dir
}

// DefaultOutputDir returns the default output directory for generated files.
func (g *Generator) DefaultOutputDir() string {
	return defaultOutputDir
}

// PostProcess file runs gofmt and goimports on the given file.
func (g *Generator) PostProcess(f *os.File) error {
	contents, err := ioutil.ReadFile(f.Name())
	if err != nil {
		return err
	}
	contents, err = imports.Process(f.Name(), contents, nil)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(f.Name(), contents, 0)
}

// GenerateDependencies is a no-op.
func (g *Generator) GenerateDependencies(dir string) error {
	return nil
}

// GenerateFile generates the given FileType.
func (g *Generator) GenerateFile(name, outputDir string, fileType generator.FileType) (*os.File, error) {
	switch fileType {
	case generator.CombinedServiceFile:
		return g.CreateFile(strings.ToLower(name)+serviceSuffix, outputDir, fileSuffix, true)
	case generator.CombinedScopeFile:
		return g.CreateFile(strings.ToLower(name)+scopeSuffix, outputDir, fileSuffix, true)
	case generator.TypeFile:
		return g.CreateFile("types", outputDir, fileSuffix, true)
	case generator.ServiceArgsResultsFile:
		return g.CreateFile(strings.ToLower(name), outputDir, fileSuffix, true)
	default:
		return nil, fmt.Errorf("Bad file type for gopherjs generator: %s", fileType)
	}
}

// GenerateDocStringComment generates the autogenerated notice.
func (g *Generator) GenerateDocStringComment(file *os.File) error {
	comment := fmt.Sprintf(
		"// Autogenerated by Frugal Compiler (%s)\n"+
			"// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING",
		globals.Version)

	_, err := file.WriteString(comment)
	return err
}

// GenerateServicePackage generates the package for the given service.
func (g *Generator) GenerateServicePackage(file *os.File, s *parser.Service) error {
	return g.generatePackage(file)
}

// GenerateScopePackage generates the package for the given scope.
func (g *Generator) GenerateScopePackage(file *os.File, s *parser.Scope) error {
	return g.generatePackage(file)
}

func (g *Generator) generatePackage(file *os.File) error {
	pkg := ""
	namespace := g.Frugal.Namespace(lang)
	if namespace != nil {
		components := generator.GetPackageComponents(namespace.Value)
		pkg = components[len(components)-1]
	} else {
		pkg = g.Frugal.Name
	}
	_, err := file.WriteString(fmt.Sprintf("package %s", pkg))
	return err
}

// GenerateConstantsContents generates constants.
func (g *Generator) GenerateConstantsContents(constants []*parser.Constant) error {
	// Use a const expression for basic types and an init function for complex
	// or typedef'd types
	contents := ""
	inits := []string{}

	for _, constant := range constants {
		if constant.Comment != nil {
			contents += g.GenerateInlineComment(constant.Comment, "")
		}

		cName := title(constant.Name)
		value := g.generateConstantValue(constant.Type, constant.Value)
		// Don't use underlying type so typedefs aren't consts
		if (constant.Type.IsPrimitive() || g.Frugal.IsEnum(constant.Type)) && constant.Type.Name != "binary" {
			contents += fmt.Sprintf("// %s is a constant.\n", cName)
			contents += fmt.Sprintf("const %s = %s\n\n", cName, value)
		} else {
			contents += fmt.Sprintf("var %s %s\n\n", cName, g.getGoTypeFromThriftType(constant.Type))
			inits = append(inits, fmt.Sprintf("\t%s = %s", cName, value))
		}
	}

	if len(inits) > 0 {
		contents += "func init() {\n" + strings.Join(inits, "\n") + "\n}\n\n"
	}
	g.typesFile.WriteString(contents)
	return nil
}

// generateConstantValue recursively generates the string representation of
// a, possibly complex, constant value.
func (g *Generator) generateConstantValue(t *parser.Type, value interface{}) string {
	// If the value being referenced is of type Identifier, it's referencing
	// another constant.
	identifier, ok := value.(parser.Identifier)
	if ok {
		idCtx := g.Frugal.ContextFromIdentifier(identifier)
		switch idCtx.Type {
		case parser.LocalConstant:
			return title(idCtx.Constant.Name)
		case parser.LocalEnum:
			return fmt.Sprintf("%s_%s", title(idCtx.Enum.Name), idCtx.EnumValue.Name)
		case parser.IncludeConstant:
			include := idCtx.Include.Name
			if namespace := g.Frugal.NamespaceForInclude(include, lang); namespace != nil {
				include = namespace.Value
			}
			return fmt.Sprintf("%s.%s", includeNameToReference(include), title(idCtx.Constant.Name))
		case parser.IncludeEnum:
			include := idCtx.Include.Name
			if namespace := g.Frugal.NamespaceForInclude(include, lang); namespace != nil {
				include = namespace.Value
			}
			return fmt.Sprintf("%s.%s_%s", includeNameToReference(include), title(idCtx.Enum.Name), idCtx.EnumValue.Name)
		default:
			panic(fmt.Sprintf("The Identifier %s has unexpected type %d", identifier, idCtx.Type))
		}
	}

	underlyingType := g.Frugal.UnderlyingType(t)
	if underlyingType.IsPrimitive() || underlyingType.IsContainer() {
		switch underlyingType.Name {
		case "bool", "i8", "byte", "i16", "i32", "i64", "double":
			return fmt.Sprintf("%v", value)
		case "string":
			return strconv.Quote(value.(string))
		case "binary":
			return fmt.Sprintf("[]byte(\"%s\")", value)
		case "list":
			contents := ""
			contents += fmt.Sprintf("%s{\n", g.getGoTypeFromThriftType(underlyingType))
			for _, v := range value.([]interface{}) {
				val := g.generateConstantValue(underlyingType.ValueType, v)
				contents += fmt.Sprintf("%s,\n", val)
			}
			contents += "}"
			return contents
		case "set":
			contents := ""
			contents += fmt.Sprintf("%s{\n", g.getGoTypeFromThriftType(underlyingType))
			for _, v := range value.([]interface{}) {
				val := g.generateConstantValue(underlyingType.ValueType, v)
				contents += fmt.Sprintf("%s: true,\n", val)
			}
			contents += "}"
			return contents
		case "map":
			contents := ""
			contents += fmt.Sprintf("%s{\n", g.getGoTypeFromThriftType(underlyingType))
			for _, pair := range value.([]parser.KeyValue) {
				key := g.generateConstantValue(underlyingType.KeyType, pair.Key)
				val := g.generateConstantValue(underlyingType.ValueType, pair.Value)
				contents += fmt.Sprintf("%s: %s,\n", key, val)
			}
			contents += "}"
			return contents
		}
	} else if g.Frugal.IsEnum(underlyingType) {
		return fmt.Sprintf("%d", value)
	} else if g.Frugal.IsStruct(underlyingType) {
		s := g.Frugal.FindStruct(underlyingType)
		if s == nil {
			panic("no struct for type " + underlyingType.Name)
		}

		contents := ""
		contents += fmt.Sprintf("&%s{\n", title(s.Name))

		for _, pair := range value.([]parser.KeyValue) {
			name := title(pair.KeyToString())
			for _, field := range s.Fields {
				if name == title(field.Name) {
					val := g.generateConstantValue(field.Type, pair.Value)
					contents += fmt.Sprintf("\t%s: %s,\n", name, val)
				}
			}
		}

		contents += "}"
		return contents
	}

	panic("no entry for type " + underlyingType.Name)
}

// GenerateTypeDef generates the given typedef.
func (g *Generator) GenerateTypeDef(typedef *parser.TypeDef) error {
	contents := fmt.Sprintf("// %s is a typeDef\n", title(typedef.Name))
	contents += fmt.Sprintf("type %s %s\n", title(typedef.Name), g.getGoTypeFromThriftType(typedef.Type))
	_, err := g.typesFile.WriteString(contents)
	return err
}

// GenerateEnum generates the given enum.
func (g *Generator) GenerateEnum(enum *parser.Enum) error {
	return g.enumTPL.Execute(g.typesFile, enum)
}

// GenerateStruct generates the given struct.
func (g *Generator) GenerateStruct(s *parser.Struct) error {
	return g.structTPL.Execute(g.typesFile, s)
}

// GenerateUnion generates the given union.
func (g *Generator) GenerateUnion(union *parser.Struct) error {
	return g.GenerateStruct(union)
}

// GenerateException generates the given exception.
func (g *Generator) GenerateException(exception *parser.Struct) error {
	return g.GenerateStruct(exception)
}

func (g *Generator) generateReadFieldRec(field *parser.Field, first bool) string {
	// first indicates if this is the first recursive call
	// first time calls assign to struct members instead of generating variables
	eq := ":="
	prefix := ""
	fName := field.Name
	if first {
		eq = "="
		prefix = "p."
		fName = title(field.Name)
	} else if strings.Contains(fName, "[i]") {
		eq = "="
	}
	contents := ""

	isPointerField := g.isPointerField(field)
	underlyingType := g.Frugal.UnderlyingType(field.Type)
	goOrigType := g.getGoTypeFromThriftTypePtr(field.Type, false)
	goUnderlyingType := g.getGoTypeFromThriftTypePtr(underlyingType, false)

	isEnum := g.Frugal.IsEnum(underlyingType)
	if underlyingType.IsPrimitive() || isEnum {
		thriftType := ""
		switch underlyingType.Name {
		case "bool":
			thriftType = "Bool"
		case "byte", "i8":
			thriftType = "Byte"
		case "i16":
			thriftType = "I16"
		case "i32":
			thriftType = "I32"
		case "i64":
			thriftType = "I64"
		case "double":
			thriftType = "Double"
		case "string":
			thriftType = "String"
		case "binary":
			thriftType = "Binary"
		default:
			if isEnum {
				thriftType = "I32"
			} else {
				panic("unknown thrift type: " + underlyingType.Name)
			}
		}

		cast := ""
		// enums and typedefs need to be cast
		if isEnum || goOrigType != goUnderlyingType {
			cast = goOrigType
		}

		if cast == "" {
			if isPointerField {
				contents += fmt.Sprintf("\t v := prot.Unpack%s()\n", thriftType)
				contents += fmt.Sprintf("\t\t%s%s = &v\n", prefix, fName)
			} else {
				contents += fmt.Sprintf("\t\t%s%s %s prot.Unpack%s()\n", prefix, fName, eq, thriftType)
			}
		} else if isPointerField {
			contents += fmt.Sprintf("\t\ttemp := %s(prot.Unpack%s())\n", cast, thriftType)
			contents += fmt.Sprintf("\t\t%s%s = &temp\n", prefix, fName)
		} else {
			contents += fmt.Sprintf("\t\t%s%s %s %s(prot.Unpack%s())\n", prefix, fName, eq, cast, thriftType)
		}

		// contents += "\t}\n"
	} else if g.Frugal.IsStruct(underlyingType) {
		// All structs types should start with a pointer
		// Need to extract the struct name from the package prefix
		// ie *base.APIException -> base.NewAPIException()
		lastInd := strings.LastIndex(goUnderlyingType, ".")
		if lastInd == -1 {
			lastInd = 0
		}
		initializer := fmt.Sprintf("%sNew%s()", goUnderlyingType[1:lastInd+1], goUnderlyingType[lastInd+1:])

		contents += fmt.Sprintf("\t%s%s %s %s\n", prefix, fName, eq, initializer)
		contents += fmt.Sprintf("\t%s%s.Unpack(prot)\n", prefix, fName)
	} else if underlyingType.IsContainer() {
		maybePointer := ""
		if isPointerField {
			maybePointer = "*"
		}
		// TODO 2.0 use this to get the value reading code, respecting the type,
		// instead of the current code for list and set
		switch underlyingType.Name {
		case "list":
			contents += "\tsize := prot.UnpackListBegin()\n"
			contents += "\tif size > 0 {"
			if !isPointerField {
				contents += fmt.Sprintf("\t%s%s %s make(%s, size)\n", prefix, fName, eq, goOrigType)
			} else {
				contents += fmt.Sprintf("\ttemp := make(%s, size)\n", goOrigType)
				contents += fmt.Sprintf("\t%s%s %s &temp\n", prefix, fName, eq)
			}
			contents += "\tfor i := 0; i < size; i++ {\n"
			valElem := g.GetElem()
			valField := parser.FieldFromType(underlyingType.ValueType, valElem)
			valField.Name = fmt.Sprintf("(%s%s%s)[i]", maybePointer, prefix, fName)
			contents += g.generateReadFieldRec(valField, false)
			contents += "\t}\n\t}\n"
			contents += "\tprot.UnpackListEnd()\n"
		case "set":
			contents += fmt.Sprintf("\t// TODO: sets! %s\n", fName)
		case "map":
			contents += "\tsize := prot.UnpackMapBegin()\n"
			if !isPointerField {
				contents += fmt.Sprintf("\t%s%s %s make(%s, size)\n", prefix, fName, eq, goOrigType)
			} else {
				contents += fmt.Sprintf("\ttemp := make(%s, size)\n", goOrigType)
				contents += fmt.Sprintf("\t%s%s %s &temp\n", prefix, fName, eq)
			}
			contents += "\tfor i := 0; i < size; i++ {\n"
			keyElem := g.GetElem()
			keyField := parser.FieldFromType(underlyingType.KeyType, keyElem)
			contents += g.generateReadFieldRec(keyField, false)
			// TODO 2.0 use the valContents for all the collections
			valElem := g.GetElem()
			valField := parser.FieldFromType(underlyingType.ValueType, valElem)
			contents += g.generateReadFieldRec(valField, false)
			contents += fmt.Sprintf("\t\t(%s%s%s)[%s] = %s\n", maybePointer, prefix, fName, keyElem, valElem)
			contents += "\t}\n"
			contents += "\tprot.UnpackMapEnd()\n"
		default:
			panic("unrecognized thrift type: " + underlyingType.Name)
		}
	}
	return contents
}

func (g *Generator) generateWriteFieldRec(field *parser.Field, prefix string) string {
	underlyingType := g.Frugal.UnderlyingType(field.Type)
	isPointerField := g.isPointerField(field)
	fName := title(field.Name)
	contents := ""

	var trailer string
	if field.Modifier == parser.Optional {
		trailer = "\t}\n"
		contents += "if " + prefix + fName + " != nil {\n"
	}

	isEnum := g.Frugal.IsEnum(underlyingType)
	if underlyingType.IsPrimitive() || isEnum {
		if isPointerField {
			prefix = "*" + prefix
		}

		write := "Pack"
		switch underlyingType.Name {
		// Just typecast everything to get around typedefs
		case "bool":
			write += "Bool(%q, %d, bool(%s))"
		case "byte", "i8":
			write += "Byte(%q, %d, int8(%s))"
		case "i16":
			write += "I16(%q, %d, int16(%s))"
		case "i32":
			write += "I32(%q, %d, int32(%s))"
		case "i64":
			write += "I64(%q, %d, int64(%s))"
		case "double":
			write += "Double(%q, %d, float64(%s))"
		case "string":
			write += "String(%q, %d, string(%s))"
		case "binary":
			write += "Binary(%q, %d, []byte(%s))"
		default:
			if isEnum {
				write += "I32(%q, %d, int32(%s))"
			} else {
				panic("unknown thrift type: " + underlyingType.Name)
			}
		}
		write = fmt.Sprintf(write, field.Name, field.ID, prefix+fName)
		contents += fmt.Sprintf("\tprot.%s\n", write)
	} else if g.Frugal.IsStruct(underlyingType) {
		contents += fmt.Sprintf("\tprot.PackFieldBegin(%q, frugal.STRUCT, %d)\n", field.Name, field.ID)
		contents += fmt.Sprintf("\t%s.Pack(prot)\n", prefix+fName)
		contents += fmt.Sprintf("\tprot.PackFieldEnd(%d)\n", field.ID)
	} else if underlyingType.IsContainer() {
		if isPointerField {
			prefix = "*" + prefix
		}
		valEnumType := g.getEnumFromThriftType(underlyingType.ValueType)
		valField := parser.FieldFromType(underlyingType.ValueType, "")
		valField.ID = -1

		switch underlyingType.Name {
		case "list":
			contents += fmt.Sprintf("\tprot.PackListBegin(%q, %d, %s, len(%s))\n", field.Name, field.ID, valEnumType, prefix+fName)
			contents += fmt.Sprintf("\tfor _, v := range %s {\n", prefix+fName)
			contents += g.generateWriteFieldRec(valField, "v")
			contents += "\t}\n"
			contents += fmt.Sprintf("\tprot.PackListEnd(%d)\n", field.ID)
		case "set":
			contents += fmt.Sprintf("\t// TODO: sets %s\n", prefix+fName)
		case "map":
			keyEnumType := g.getEnumFromThriftType(underlyingType.KeyType)
			contents += fmt.Sprintf("\tprot.PackMapBegin(%q, %d, %s, %s, len(%s))\n", field.Name, field.ID, keyEnumType, valEnumType, prefix+fName)
			contents += fmt.Sprintf("\tfor k, v := range %s {\n", prefix+fName)
			keyField := parser.FieldFromType(underlyingType.KeyType, "")
			keyField.ID = -1
			contents += g.generateWriteFieldRec(keyField, "k")
			contents += g.generateWriteFieldRec(valField, "v")
			contents += "\t}\n"
			contents += fmt.Sprintf("\tprot.PackMapEnd(%d)\n", field.ID)
		default:
			panic("unknow type: " + underlyingType.Name)
		}
	}

	return contents + trailer
}

// GenerateTypesImports generates the necessary Go types imports.
func (g *Generator) GenerateTypesImports(file *os.File) error {
	contents := "import (\n"
	if g.Options[frugalImportOption] != "" {
		contents += "\t\"" + g.Options[frugalImportOption] + "\"\n"
	} else {
		contents += "\t\"github.com/Workiva/frugal/lib/gopherjs/frugal\"\n"
	}

	pkgPrefix := g.Options[packagePrefixOption]
	for _, include := range g.Frugal.Includes {
		imp, err := g.generateIncludeImport(include, pkgPrefix)
		if err != nil {
			return err
		}
		contents += imp
	}

	contents += ")\n\n"
	_, err := file.WriteString(contents)
	return err
}

// GenerateServiceResultArgsImports generates the necessary imports for service
// args and result types.
func (g *Generator) GenerateServiceResultArgsImports(file *os.File) error {
	contents := "import (\n"

	pkgPrefix := g.Options[packagePrefixOption]
	for _, include := range g.Frugal.Includes {
		if imp, err := g.generateIncludeImport(include, pkgPrefix); err != nil {
			return err
		} else {
			contents += imp
		}
	}

	contents += ")\n\n"

	_, err := file.WriteString(contents)
	return err
}

// GenerateServiceImports generates necessary imports for the given service.
func (g *Generator) GenerateServiceImports(file *os.File, s *parser.Service) error {
	imports := "import (\n"
	if g.Options[frugalImportOption] != "" {
		imports += "\t\"" + g.Options[frugalImportOption] + "\"\n"
	} else {
		imports += "\t\"github.com/Workiva/frugal/lib/gopherjs/frugal\"\n"
	}

	pkgPrefix := g.Options[packagePrefixOption]
	includes, err := s.ReferencedIncludes()
	if err != nil {
		return err
	}
	for _, include := range includes {
		if imp, err := g.generateIncludeImport(include, pkgPrefix); err != nil {
			return err
		} else {
			imports += imp
		}
	}

	imports += ")\n\n"

	_, err = file.WriteString(imports)
	return err
}

// GenerateScopeImports generates necessary imports for the given scope.
func (g *Generator) GenerateScopeImports(file *os.File, s *parser.Scope) error {
	imports := "import (\n"
	if g.Options[frugalImportOption] != "" {
		imports += "\t\"" + g.Options[frugalImportOption] + "\"\n"
	} else {
		imports += "\t\"github.com/Workiva/frugal/lib/go\"\n"
	}

	pkgPrefix := g.Options[packagePrefixOption]
	scopeIncludes, err := g.Frugal.ReferencedScopeIncludes()
	if err != nil {
		return err
	}
	for _, include := range scopeIncludes {
		if imp, err := g.generateIncludeImport(include, pkgPrefix); err != nil {
			return err
		} else {
			imports += imp
		}
	}

	imports += ")"

	_, err = file.WriteString(imports)
	return err
}

func (g *Generator) generateIncludeImport(include *parser.Include, pkgPrefix string) (string, error) {
	includeName := filepath.Base(include.Name)
	importPath := fmt.Sprintf("%s%s", pkgPrefix, includeNameToImport(includeName))
	namespace := g.Frugal.NamespaceForInclude(includeName, lang)

	_, vendored := include.Annotations.Vendor()
	vendored = vendored && g.UseVendor()
	vendorPath := ""

	if namespace != nil {
		importPath = fmt.Sprintf("%s%s", pkgPrefix, includeNameToImport(namespace.Value))
		if nsVendorPath, ok := namespace.Annotations.Vendor(); ok {
			vendorPath = nsVendorPath
		}
	}

	// If -use-vendor is set and this include is vendored, honor the path
	// specified by the include's namespace vendor annotation.
	if vendored {
		if vendorPath == "" {
			return "", fmt.Errorf("Vendored include %s does not specify vendor path for go namespace",
				include.Name)
		}
		importPath = vendorPath
	}

	return fmt.Sprintf("\t\"%s\"\n", importPath), nil
}

// GenerateConstants generates any static constants.
func (g *Generator) GenerateConstants(file *os.File, name string) error {
	return nil
}

// GeneratePublisher generates the publisher for the given scope.
func (g *Generator) GeneratePublisher(file *os.File, scope *parser.Scope) error {
	return nil
	// var (
	// 	scopeLower = parser.LowercaseFirstLetter(scope.Name)
	// 	scopeCamel = snakeToCamel(scope.Name)
	// 	publisher  = ""
	// )
	//
	// if scope.Comment != nil {
	// 	publisher += g.GenerateInlineComment(scope.Comment, "")
	// }
	// args := ""
	// if len(scope.Prefix.Variables) > 0 {
	// 	prefix := ""
	// 	for _, variable := range scope.Prefix.Variables {
	// 		args += prefix + variable
	// 		prefix = ", "
	// 	}
	// 	args += " string, "
	// }
	//
	// publisher += fmt.Sprintf("type %sPublisher interface {\n", scopeCamel)
	// publisher += "\tOpen() error\n"
	// publisher += "\tClose() error\n"
	// for _, op := range scope.Operations {
	// 	publisher += fmt.Sprintf("\tPublish%s(ctx frugal.Context, %sreq %s) error\n", op.Name, args, g.getGoTypeFromThriftType(op.Type))
	// }
	// publisher += "}\n\n"
	//
	// publisher += fmt.Sprintf("type %sPublisher struct {\n", scopeLower)
	// publisher += "\ttransport frugal.FPublisherTransport\n"
	// publisher += "\tprotocolFactory *frugal.FProtocolFactory\n"
	// publisher += "\tmethods   map[string]*frugal.Method\n"
	// publisher += "}\n\n"
	//
	// publisher += fmt.Sprintf("func New%sPublisher(provider *frugal.FScopeProvider, middleware ...frugal.ServiceMiddleware) %sPublisher {\n",
	// 	scopeCamel, scopeCamel)
	// publisher += "\ttransport, protocolFactory := provider.NewPublisher()\n"
	// publisher += "\tmethods := make(map[string]*frugal.Method)\n"
	// publisher += fmt.Sprintf("\tpublisher := &%sPublisher{\n", scopeLower)
	// publisher += "\t\ttransport: transport,\n"
	// publisher += "\t\tprotocolFactory:  protocolFactory,\n"
	// publisher += "\t\tmethods:   methods,\n"
	// publisher += "\t}\n"
	// publisher += "\tmiddleware = append(middleware, provider.GetMiddleware()...)\n"
	// for _, op := range scope.Operations {
	// 	publisher += fmt.Sprintf("\tmethods[\"publish%s\"] = frugal.NewMethod(publisher, publisher.publish%s, \"publish%s\", middleware)\n",
	// 		op.Name, op.Name, op.Name)
	// }
	// publisher += "\treturn publisher\n"
	// publisher += "}\n\n"
	//
	// publisher += fmt.Sprintf("func (p *%sPublisher) Open() error {\n", scopeLower)
	//
	// publisher += "\treturn p.transport.Open()\n"
	// publisher += "}\n\n"
	//
	// publisher += fmt.Sprintf("func (p *%sPublisher) Close() error {\n", scopeLower)
	// publisher += "\treturn p.transport.Close()\n"
	// publisher += "}\n\n"
	//
	// prefix := ""
	// for _, op := range scope.Operations {
	// 	publisher += prefix
	// 	prefix = "\n\n"
	// 	publisher += g.generatePublishMethod(scope, op, args)
	// }
	//
	// _, err := file.WriteString(publisher)
	// return err
}

func (g *Generator) generatePublishMethod(scope *parser.Scope, op *parser.Operation, args string) string {
	var (
		// scopeLower = parser.LowercaseFirstLetter(scope.Name)
		publisher = ""
	)

	// 	if op.Comment != nil {
	// 		publisher += g.GenerateInlineComment(op.Comment, "")
	// 	}
	//
	// 	publisher += fmt.Sprintf("func (p *%sPublisher) Publish%s(ctx frugal.Context, %sreq %s) error {\n",
	// 		scopeLower, op.Name, args, g.getGoTypeFromThriftType(op.Type))
	// 	publisher += fmt.Sprintf("\tret := p.methods[\"publish%s\"].Invoke(%s)\n", op.Name, g.generateScopeArgs(scope))
	// 	publisher += "\tif ret[0] != nil {\n"
	// 	publisher += "\t\treturn ret[0].(error)\n"
	// 	publisher += "\t}\n"
	// 	publisher += "\treturn nil\n"
	// 	publisher += "}\n\n"
	//
	// 	publisher += g.generateInternalPublishMethod(scope, op, args)
	//
	// 	return publisher
	// }
	//
	// func (g *Generator) generateInternalPublishMethod(scope *parser.Scope, op *parser.Operation, args string) string {
	// 	var (
	// 		scopeLower = parser.LowercaseFirstLetter(scope.Name)
	// 		scopeTitle = strings.Title(scope.Name)
	// 		publisher  = ""
	// 	)
	//
	// 	publisher += fmt.Sprintf("func (p *%sPublisher) publish%s(ctx frugal.Context, %sreq %s) error {\n",
	// 		scopeLower, op.Name, args, g.getGoTypeFromThriftType(op.Type))
	//
	// 	// Inject the prefix variables into the FContext to send
	// 	for _, prefixVar := range scope.Prefix.Variables {
	// 		publisher += fmt.Sprintf("\tctx.AddRequestHeader(\"_topic_%s\", %s)\n", prefixVar, prefixVar)
	// 	}
	//
	// 	publisher += fmt.Sprintf("\top := \"%s\"\n", op.Name)
	// 	publisher += fmt.Sprintf("\tprefix := %s\n", generatePrefixStringTemplate(scope))
	// 	publisher += "\ttopic := fmt.Sprintf(\"%s" + scopeTitle + "%s%s\", prefix, delimiter, op)\n"
	// 	publisher += "\tbuffer := frugal.NewTMemoryOutputBuffer(p.transport.GetPublishSizeLimit())\n"
	// 	publisher += "\toprot := p.protocolFactory.GetProtocol(buffer)\n"
	// 	publisher += "\tif err := oprot.WriteRequestHeader(ctx); err != nil {\n"
	// 	publisher += "\t\treturn err\n"
	// 	publisher += "\t}\n"
	// 	publisher += "\tif err := oprot.WriteMessageBegin(op, thrift.CALL, 0); err != nil {\n"
	// 	publisher += "\t\treturn err\n"
	// 	publisher += "\t}\n"
	// 	publisher += g.generateWriteFieldRec(parser.FieldFromType(op.Type, ""), "req")
	// 	publisher += "\tif err := oprot.WriteMessageEnd(); err != nil {\n"
	// 	publisher += "\t\treturn err\n"
	// 	publisher += "\t}\n"
	// 	publisher += "\tif err := oprot.Flush(); err != nil {\n"
	// 	publisher += "\t\treturn err\n"
	// 	publisher += "\t}\n"
	// 	publisher += "\treturn p.transport.Publish(topic, buffer.Bytes())\n"
	// 	publisher += "}\n"
	return publisher
}

func generatePrefixStringTemplate(scope *parser.Scope) string {
	if len(scope.Prefix.Variables) == 0 {
		if scope.Prefix.String == "" {
			return `""`
		}
		return fmt.Sprintf(`"%s%s"`, scope.Prefix.String, globals.TopicDelimiter)
	}
	template := "fmt.Sprintf(\""
	template += scope.Prefix.Template("%s")
	template += globals.TopicDelimiter + "\", "
	prefix := ""
	for _, variable := range scope.Prefix.Variables {
		template += prefix + variable
		prefix = ", "
	}
	template += ")"
	return template
}

// GenerateSubscriber generates the subscriber for the given scope.
func (g *Generator) GenerateSubscriber(file *os.File, scope *parser.Scope) error {
	return nil
	// var (
	// 	scopeLower = parser.LowercaseFirstLetter(scope.Name)
	// 	scopeCamel = snakeToCamel(scope.Name)
	// 	subscriber = ""
	// )
	//
	// if scope.Comment != nil {
	// 	subscriber += g.GenerateInlineComment(scope.Comment, "")
	// }
	//
	// args := ""
	// argsWithoutTypes := ""
	// prefix := ""
	// if len(scope.Prefix.Variables) > 0 {
	// 	for _, variable := range scope.Prefix.Variables {
	// 		args += prefix + variable
	// 		prefix = ", "
	// 	}
	// 	argsWithoutTypes = args + ", "
	// 	args += " string, "
	// }
	//
	// subscriber += fmt.Sprintf("type %sSubscriber interface {\n", scopeCamel)
	// for _, op := range scope.Operations {
	// 	subscriber += fmt.Sprintf("\tSubscribe%s(%shandler func(frugal.Context, %s)) (*frugal.FSubscription, error)\n",
	// 		op.Name, args, g.getGoTypeFromThriftType(op.Type))
	// }
	// subscriber += "}\n\n"
	//
	// if scope.Comment != nil {
	// 	subscriber += g.GenerateInlineComment(scope.Comment, "")
	// }
	// subscriber += fmt.Sprintf("type %sErrorableSubscriber interface {\n", scopeCamel)
	// for _, op := range scope.Operations {
	// 	subscriber += fmt.Sprintf("\tSubscribe%sErrorable(%shandler func(frugal.Context, %s) error) (*frugal.FSubscription, error)\n",
	// 		op.Name, args, g.getGoTypeFromThriftType(op.Type))
	// }
	// subscriber += "}\n\n"
	//
	// subscriber += fmt.Sprintf("type %sSubscriber struct {\n", scopeLower)
	// subscriber += "\tprovider   *frugal.FScopeProvider\n"
	// subscriber += "\tmiddleware []frugal.ServiceMiddleware\n"
	// subscriber += "}\n\n"
	//
	// subscriber += fmt.Sprintf("func New%sSubscriber(provider *frugal.FScopeProvider, middleware ...frugal.ServiceMiddleware) %sSubscriber {\n",
	// 	scopeCamel, scopeCamel)
	// subscriber += "\tmiddleware = append(middleware, provider.GetMiddleware()...)\n"
	// subscriber += fmt.Sprintf("\treturn &%sSubscriber{provider: provider, middleware: middleware}\n", scopeLower)
	// subscriber += "}\n\n"
	//
	// subscriber += fmt.Sprintf("func New%sErrorableSubscriber(provider *frugal.FScopeProvider, middleware ...frugal.ServiceMiddleware) %sErrorableSubscriber {\n",
	// 	scopeCamel, scopeCamel)
	// subscriber += "\tmiddleware = append(middleware, provider.GetMiddleware()...)\n"
	// subscriber += fmt.Sprintf("\treturn &%sSubscriber{provider: provider, middleware: middleware}\n", scopeLower)
	// subscriber += "}\n\n"
	//
	// prefix = ""
	// for _, op := range scope.Operations {
	// 	subscriber += prefix
	// 	prefix = "\n\n"
	// 	subscriber += g.generateSubscribeMethod(scope, op, args, argsWithoutTypes)
	// }
	//
	// _, err := file.WriteString(subscriber)
	// return err
}

func (g *Generator) generateSubscribeMethod(scope *parser.Scope, op *parser.Operation, args, argsWithoutTypes string) string {
	var (
		// scopeLower = parser.LowercaseFirstLetter(scope.Name)
		// scopeTitle = strings.Title(scope.Name)
		subscriber = ""
	)
	// if op.Comment != nil {
	// 	subscriber += g.GenerateInlineComment(op.Comment, "")
	// }
	//
	// subscriber += fmt.Sprintf("func (l *%sSubscriber) Subscribe%s(%shandler func(frugal.Context, %s)) (*frugal.FSubscription, error) {\n",
	// 	scopeLower, op.Name, args, g.getGoTypeFromThriftType(op.Type))
	// subscriber += fmt.Sprintf("\treturn l.Subscribe%sErrorable(%sfunc(fctx frugal.Context, arg %s) error {\n",
	// 	op.Name, argsWithoutTypes, g.getGoTypeFromThriftType(op.Type))
	// subscriber += "\t\thandler(fctx, arg)\n"
	// subscriber += "\t\treturn nil\n"
	// subscriber += "\t})\n"
	// subscriber += "}\n\n"
	//
	// if op.Comment != nil {
	// 	subscriber += g.GenerateInlineComment(op.Comment, "")
	// }
	// subscriber += fmt.Sprintf("func (l *%sSubscriber) Subscribe%sErrorable(%shandler func(frugal.Context, %s) error) (*frugal.FSubscription, error) {\n",
	// 	scopeLower, op.Name, args, g.getGoTypeFromThriftType(op.Type))
	// subscriber += fmt.Sprintf("\top := \"%s\"\n", op.Name)
	// subscriber += fmt.Sprintf("\tprefix := %s\n", generatePrefixStringTemplate(scope))
	// subscriber += "\ttopic := fmt.Sprintf(\"%s" + scopeTitle + "%s%s\", prefix, delimiter, op)\n"
	// subscriber += "\ttransport, protocolFactory := l.provider.NewSubscriber()\n"
	// subscriber += fmt.Sprintf("\tcb := l.recv%s(op, protocolFactory, handler)\n", op.Name)
	// subscriber += "\tif err := transport.Subscribe(topic, cb); err != nil {\n"
	// subscriber += "\t\treturn nil, err\n"
	// subscriber += "\t}\n\n"
	//
	// subscriber += "\tsub := frugal.NewFSubscription(topic, transport)\n"
	// subscriber += "\treturn sub, nil\n"
	// subscriber += "}\n\n"
	//
	// subscriber += fmt.Sprintf("func (l *%sSubscriber) recv%s(op string, pf *frugal.FProtocolFactory, handler func(frugal.Context, %s) error) frugal.FAsyncCallback {\n",
	// 	scopeLower, op.Name, g.getGoTypeFromThriftType(op.Type))
	// subscriber += fmt.Sprintf("\tmethod := frugal.NewMethod(l, handler, \"Subscribe%s\", l.middleware)\n", op.Name)
	// subscriber += "\treturn func(transport thrift.TTransport) error {\n"
	// subscriber += "\t\tiprot := pf.GetProtocol(transport)\n"
	// subscriber += "\t\tctx, err := iprot.ReadRequestHeader()\n"
	// subscriber += "\t\tif err != nil {\n"
	// subscriber += "\t\t\treturn err\n"
	// subscriber += "\t\t}\n\n"
	// subscriber += "\t\tname, _, _, err := iprot.ReadMessageBegin()\n"
	// subscriber += "\t\tif err != nil {\n"
	// subscriber += "\t\t\treturn err\n"
	// subscriber += "\t\t}\n\n"
	// subscriber += "\t\tif name != op {\n"
	// subscriber += "\t\t\tiprot.Skip(thrift.STRUCT)\n"
	// subscriber += "\t\t\tiprot.ReadMessageEnd()\n"
	// subscriber += "\t\t\treturn thrift.NewTApplicationException(frugal.APPLICATION_EXCEPTION_UNKNOWN_METHOD, \"Unknown function\"+name)\n"
	// subscriber += "\t\t}\n"
	// subscriber += g.generateReadFieldRec(parser.FieldFromType(op.Type, "req"), false)
	// subscriber += "\t\tiprot.ReadMessageEnd()\n\n"
	// subscriber += "\t\treturn method.Invoke([]interface{}{ctx, req}).Error()\n"
	// subscriber += "\t}\n"
	// subscriber += "}"

	return subscriber
}

// GenerateService generates the given service.
func (g *Generator) GenerateService(file *os.File, s *parser.Service) error {
	if err := g.serviceTPL.Execute(file, s); err != nil {
		return err
	}
	for _, typ := range g.GetServiceMethodTypes(s) {
		typ.Name = title(s.Name) + title(typ.Name) // prepend service to type name
		if err := g.structTPL.Execute(file, typ); err != nil {
			return err
		}
	}
	return nil
}

func (g *Generator) getServiceExtendsNamespace(service *parser.Service) string {
	namespace := ""
	if service.ExtendsInclude() != "" {
		if ns := g.Frugal.NamespaceForInclude(service.ExtendsInclude(), lang); ns != nil {
			namespace = ns.Value
		} else {
			namespace = service.ExtendsInclude()
		}
		namespace = includeNameToReference(namespace)
		namespace += "."
	}
	return namespace
}

func (g *Generator) generateServer(service *parser.Service) string {
	contents := ""
	contents += g.generateProcessor(service)
	for _, method := range service.Methods {
		contents += g.generateMethodProcessor(service, method)
	}
	return contents
}

func (g *Generator) generateProcessor(service *parser.Service) string {
	var (
		servTitle = snakeToCamel(service.Name)
		servLower = strings.ToLower(service.Name)
		contents  = ""
	)

	contents += fmt.Sprintf("type F%sProcessor struct {\n", servTitle)
	if service.Extends == "" {
		contents += "\t*frugal.FBaseProcessor\n"
	} else {
		contents += fmt.Sprintf("\t*%sF%sProcessor\n",
			g.getServiceExtendsNamespace(service), service.ExtendsService())
	}
	contents += "}\n\n"

	contents += fmt.Sprintf("func NewF%sProcessor(handler F%s, middleware ...frugal.ServiceMiddleware) *F%sProcessor {\n",
		servTitle, servTitle, servTitle)
	if service.Extends != "" {
		contents += fmt.Sprintf("\tp := &F%sProcessor{%sNewF%sProcessor(handler, middleware...)}\n",
			servTitle, g.getServiceExtendsNamespace(service), service.ExtendsService())
	} else {
		contents += fmt.Sprintf("\tp := &F%sProcessor{frugal.NewFBaseProcessor()}\n", servTitle)
	}
	for _, method := range service.Methods {
		methodLower := parser.LowercaseFirstLetter(method.Name)
		contents += fmt.Sprintf(
			"\tp.AddToProcessorMap(\"%s\", &%sF%s{frugal.NewFBaseProcessorFunction(p.GetWriteMutex(), frugal.NewMethod(handler, handler.%s, \"%s\", middleware))})\n",
			methodLower, servLower, snakeToCamel(method.Name), snakeToCamel(method.Name), snakeToCamel(method.Name))
		if len(method.Annotations) > 0 {
			contents += fmt.Sprintf("\tp.AddToAnnotationsMap(\"%s\", map[string]string{\n", methodLower)
			for _, annotation := range method.Annotations {
				contents += fmt.Sprintf("\t\t\"%s\": %s,\n", annotation.Name, strconv.Quote(annotation.Value))
			}
			contents += "\t})\n"
		}
	}

	contents += "\treturn p\n"
	contents += "}\n\n"

	return contents
}

func (g *Generator) generateMethodProcessor(service *parser.Service, method *parser.Method) string {
	var (
		servTitle = snakeToCamel(service.Name)
		servLower = strings.ToLower(service.Name)
		nameTitle = snakeToCamel(method.Name)
		nameLower = parser.LowercaseFirstLetter(method.Name)
	)

	contents := fmt.Sprintf("type %sF%s struct {\n", servLower, nameTitle)
	contents += "\t*frugal.FBaseProcessorFunction\n"
	contents += "}\n\n"

	contents += fmt.Sprintf("func (p *%sF%s) Process(ctx frugal.Context, iprot, oprot *frugal.FProtocol) error {\n", servLower, nameTitle)

	if _, ok := method.Annotations.Deprecated(); ok {
		contents += fmt.Sprintf("\tlogrus.Warn(\"Deprecated function '%s.%s' was called by a client\")\n", service.Name, nameTitle)
	}

	contents += fmt.Sprintf("\targs := %s%sArgs{}\n", servTitle, nameTitle)
	contents += "\tvar err error\n"
	contents += "\tif err = args.Read(iprot); err != nil {\n"
	contents += "\t\tiprot.ReadMessageEnd()\n"
	if !method.Oneway {
		contents += "\t\tp.GetWriteMutex().Lock()\n"
		contents += fmt.Sprintf("\t\terr = %sWriteApplicationError(ctx, oprot, frugal.APPLICATION_EXCEPTION_PROTOCOL_ERROR, \"%s\", err.Error())\n", servLower, nameLower)
		contents += "\t\tp.GetWriteMutex().Unlock()\n"
	}
	contents += "\t\treturn err\n"
	contents += "\t}\n\n"

	contents += "\tiprot.ReadMessageEnd()\n"
	if !method.Oneway {
		contents += fmt.Sprintf("\tresult := %s%sResult{}\n", servTitle, nameTitle)
	}
	contents += "\tvar err2 error\n"
	if method.ReturnType != nil {
	}
	contents += fmt.Sprintf("\tret := p.InvokeMethod(%s)\n", g.generateHandlerArgs(method))
	numReturn := "2"
	if method.ReturnType == nil {
		numReturn = "1"
	}
	contents += fmt.Sprintf("\tif len(ret) != %s {\n", numReturn)
	contents += fmt.Sprintf("\t\tpanic(fmt.Sprintf(\"Middleware returned %%d arguments, expected %s\", len(ret)))\n", numReturn)
	contents += "\t}\n"
	if method.ReturnType != nil {
		contents += "\tif ret[1] != nil {\n"
		contents += "\t\terr2 = ret[1].(error)\n"
		contents += "\t}\n"
	} else {
		contents += "\tif ret[0] != nil {\n"
		contents += "\t\terr2 = ret[0].(error)\n"
		contents += "\t}\n"
	}
	contents += "\tif err2 != nil {\n"
	contents += "\t\tif err3, ok := err2.(thrift.TApplicationException); ok {\n"
	contents += "\t\t\tp.GetWriteMutex().Lock()\n"
	contents += "\t\t\toprot.WriteResponseHeader(ctx)\n"
	contents += fmt.Sprintf("\t\t\toprot.WriteMessageBegin(\"%s\", thrift.EXCEPTION, 0)\n", nameLower)
	contents += "\t\t\terr3.Write(oprot)\n"
	contents += "\t\t\toprot.WriteMessageEnd()\n"
	contents += "\t\t\toprot.Flush()\n"
	contents += "\t\t\tp.GetWriteMutex().Unlock()\n"
	contents += "\t\t\treturn nil\n"
	contents += "\t\t}\n"
	if len(method.Exceptions) > 0 {
		contents += "\t\tswitch v := err2.(type) {\n"
		for _, err := range method.Exceptions {
			contents += fmt.Sprintf("\t\tcase %s:\n", g.getGoTypeFromThriftType(err.Type))
			contents += fmt.Sprintf("\t\t\tresult.%s = v\n", snakeToCamel(err.Name))
		}
		contents += "\t\tdefault:\n"
		contents += g.generateMethodException("\t\t\t", service, method)
		contents += "\t\t}\n"
	} else {
		contents += g.generateMethodException("\t\t", service, method)
	}
	if method.ReturnType != nil {
		contents += "\t} else {\n"
		contents += fmt.Sprintf("\t\tvar retval %s = ret[0].(%s)\n",
			g.getGoTypeFromThriftType(method.ReturnType), g.getGoTypeFromThriftType(method.ReturnType))
		if g.isPrimitive(method.ReturnType) || g.Frugal.IsEnum(method.ReturnType) {
			contents += "\t\tresult.Success = &retval\n"
		} else {
			contents += "\t\tresult.Success = retval\n"
		}
	}
	contents += "\t}\n"

	if method.Oneway {
		contents += "\treturn err\n"
		contents += "}\n\n"
		return contents
	}

	contents += "\tp.GetWriteMutex().Lock()\n"
	contents += "\tdefer p.GetWriteMutex().Unlock()\n"
	contents += "\tif err2 = oprot.WriteResponseHeader(ctx); err2 != nil {\n"
	contents += g.generateErrTooLarge(service, method)
	contents += "\t}\n"
	contents += fmt.Sprintf("\tif err2 = oprot.WriteMessageBegin(\"%s\", "+
		"thrift.REPLY, 0); err2 != nil {\n", nameLower)
	contents += g.generateErrTooLarge(service, method)
	contents += "\t}\n"
	contents += "\tif err2 = result.Write(oprot); err == nil && err2 != nil {\n"
	contents += g.generateErrTooLarge(service, method)
	contents += "\t}\n"
	contents += "\tif err2 = oprot.WriteMessageEnd(); err == nil && err2 != nil {\n"
	contents += g.generateErrTooLarge(service, method)
	contents += "\t}\n"
	contents += "\tif err2 = oprot.Flush(); err == nil && err2 != nil {\n"
	contents += g.generateErrTooLarge(service, method)
	contents += "\t}\n"
	contents += "\treturn err\n"
	contents += "}\n\n"

	return contents
}

func (g *Generator) generateClientArgs(method *parser.Method) string {
	args := "[]interface{}{ctx"
	for _, arg := range method.Arguments {
		args += ", " + strings.ToLower(arg.Name)
	}
	args += "}"
	return args
}

func (g *Generator) generateScopeArgs(scope *parser.Scope) string {
	args := "[]interface{}{ctx"
	for _, v := range scope.Prefix.Variables {
		args += ", " + v
	}
	args += ", req"
	args += "}"
	return args
}

func (g *Generator) generateHandlerArgs(method *parser.Method) string {
	args := "[]interface{}{ctx"
	for _, arg := range method.Arguments {
		args += ", args." + title(arg.Name)
	}
	args += "}"
	return args
}

func (g *Generator) generateErrTooLarge(service *parser.Service, method *parser.Method) string {
	servLower := strings.ToLower(service.Name)
	nameLower := parser.LowercaseFirstLetter(method.Name)
	contents := "\t\tif frugal.IsErrTooLarge(err2) {\n"
	contents += fmt.Sprintf(
		"\t\t\t%sWriteApplicationError(ctx, oprot, frugal.APPLICATION_EXCEPTION_RESPONSE_TOO_LARGE, \"%s\", err2.Error())\n",
		servLower, nameLower)
	contents += "\t\t\treturn nil\n"
	contents += "\t\t}\n"
	contents += "\t\terr = err2"
	return contents
}

func (g *Generator) generateMethodException(prefix string, service *parser.Service, method *parser.Method) string {
	contents := ""
	servLower := strings.ToLower(service.Name)
	nameLower := parser.LowercaseFirstLetter(method.Name)
	if !method.Oneway {
		contents += prefix + "p.GetWriteMutex().Lock()\n"
		msg := fmt.Sprintf("\"Internal error processing %s: \"+err2.Error()", nameLower)
		contents += fmt.Sprintf(
			prefix+"err2 := %sWriteApplicationError(ctx, oprot, frugal.APPLICATION_EXCEPTION_INTERNAL_ERROR, \"%s\", %s)\n", servLower, nameLower, msg)
		contents += prefix + "p.GetWriteMutex().Unlock()\n"
	}
	contents += prefix + "return err2\n"
	return contents
}

func (g *Generator) getGoTypeFromThriftType(t *parser.Type) string {
	return g.getGoTypeFromThriftTypePtr(t, false)
}

func (g *Generator) getGoTypeFromThriftTypePtr(t *parser.Type, pointer bool) string {
	maybePointer := ""
	if pointer {
		maybePointer = "*"
	}
	switch t.Name {
	case "bool":
		return maybePointer + "bool"
	case "byte", "i8":
		return maybePointer + "int8"
	case "i16":
		return maybePointer + "int16"
	case "i32":
		return maybePointer + "int32"
	case "i64":
		return maybePointer + "int64"
	case "double":
		return maybePointer + "float64"
	case "string":
		return maybePointer + "string"
	case "binary":
		return maybePointer + "[]byte"
	case "list":
		return fmt.Sprintf("%s[]%s", maybePointer,
			g.getGoTypeFromThriftTypePtr(t.ValueType, false))
	case "set":
		return fmt.Sprintf("%smap[%s]bool", maybePointer,
			g.getGoTypeFromThriftTypePtr(t.ValueType, false))
	case "map":
		return fmt.Sprintf("%smap[%s]%s", maybePointer,
			g.getGoTypeFromThriftTypePtr(t.KeyType, false),
			g.getGoTypeFromThriftTypePtr(t.ValueType, false))
	default:
		// Custom type, either typedef or struct.
		name := g.qualifiedTypeName(t)
		if g.Frugal.IsStruct(t) {
			// This is a struct, return a pointer to it.
			return "*" + name
		}
		return maybePointer + name
	}
}

func (g *Generator) getGoTypeFromThriftTypeEnum(typ *parser.Type) string {
	switch typ.Name {
	// Just typecast everything to get around typedefs
	case "bool":
		return "bool"
	case "byte", "i8":
		return "int8"
	case "i16":
		return "int16"
	case "i32":
		return "int32"
	case "i64":
		return "int64"
	case "double":
		return "float64"
	case "string":
		return "string"
	case "binary":
		return "[]byte"
	default:
		if g.Frugal.IsEnum(typ) {
			return "int32"
		}
		panic("unknown thrift type: " + typ.Name)
	}
}

func (g *Generator) getEnumFromThriftType(t *parser.Type) string {
	underlyingType := g.Frugal.UnderlyingType(t)
	switch underlyingType.Name {
	case "bool":
		return "frugal.BOOL"
	case "byte", "i8":
		return "frugal.BYTE"
	case "i16":
		return "frugal.I16"
	case "i32":
		return "frugal.I32"
	case "i64":
		return "frugal.I64"
	case "double":
		return "frugal.DOUBLE"
	case "string", "binary":
		return "frugal.STRING"
	case "list":
		return "frugal.LIST"
	case "set":
		return "frugal.SET"
	case "map":
		return "frugal.MAP"
	default:
		if g.Frugal.IsEnum(underlyingType) {
			return "frugal.I32"
		} else if g.Frugal.IsStruct(underlyingType) {
			return "frugal.STRUCT"
		}
		panic("not a valid thrift type: " + underlyingType.Name)
	}
}

func (g *Generator) isPrimitive(t *parser.Type) bool {
	underlyingType := g.Frugal.UnderlyingType(t)
	switch underlyingType.Name {
	case "bool", "byte", "i8", "i16", "i32", "i64", "double", "string":
		return true
	default:
		return false
	}
}

func (g *Generator) isPointerField(field *parser.Field) bool {
	underlyingType := g.Frugal.UnderlyingType(field.Type)
	// Structs as fields are always pointers
	if g.Frugal.IsStruct(underlyingType) {
		return true
	}
	// If it's not optional, it's not a pointer
	if field.Modifier != parser.Optional {
		return false
	}

	hasDefault := field.Default != nil
	switch underlyingType.Name {
	case "binary":
		// According to thrift, these are always like this, not sure why
		return false
	case "bool", "byte", "i8", "i16", "i32", "i64", "double", "string":
		// If there's no default, needs to be a pointer to be nillable
		return !hasDefault
	case "list", "set", "map":
		// slices and maps are nillable by default, use a pointer
		// if there's a default to differentiate between the default and
		// not set
		return hasDefault
	default:
		// Custom type, either typedef or struct-like.
		if g.Frugal.IsEnum(underlyingType) {
			// Same case as nums
			return !hasDefault
		}
		return hasDefault
	}
}

func (g *Generator) qualifiedTypeName(t *parser.Type) string {
	param := snakeToCamel(t.ParamName())
	include := t.IncludeName()
	if include != "" {
		name := include
		if namespace := g.Frugal.NamespaceForInclude(include, lang); namespace != nil {
			name = namespace.Value
		}
		param = fmt.Sprintf("%s.%s", includeNameToReference(name), param)
	}

	// // The Thrift generator uses a convention of appending a suffix of '_'
	// // if the argument starts with 'New', ends with 'Result' or ends with 'Args'.
	// // This effort must be duplicated to correctly reference Thrift generated code.
	// if strings.HasPrefix(param, "New") || strings.HasSuffix(param, "Result") || strings.HasSuffix(param, "Args") {
	// 	param += "_"
	// }
	return param
}

// UseVendor determines if this editor supports use vendor.
func (g *Generator) UseVendor() bool {
	_, ok := g.Options["use_vendor"]
	return ok
}

func includeNameToImport(includeName string) string {
	return strings.Replace(includeName, ".", "/", -1)
}

func includeNameToReference(includeName string) string {
	split := strings.FieldsFunc(includeName, func(r rune) bool {
		return r == '.' || r == '/'
	})
	return split[len(split)-1]
}

// snakeToCamel returns a string converted from snake case to uppercase.
func snakeToCamel(s string) string {
	if len(s) == 0 {
		return s
	}

	var result string

	words := strings.Split(s, "_")

	for _, word := range words {
		if upper := strings.ToUpper(word); commonInitialisms[upper] {
			result += upper
			continue
		}

		w := []rune(word)
		w[0] = unicode.ToUpper(w[0])
		result += string(w)
	}

	return result
}

func title(name string) string {
	if len(name) == 0 {
		return name
	}

	// Keep screaming caps
	if name == strings.ToUpper(name) {
		return name
	}
	result := snakeToCamel(name)

	// if strings.HasPrefix(result, "New") || strings.HasSuffix(result, "Args") || strings.HasSuffix(result, "Result") {
	// 	result += "_"
	// }
	return result
}

// commonInitialisms, taken from
// https://github.com/golang/lint/blob/3d26dc39376c307203d3a221bada26816b3073cf/lint.go#L482
var commonInitialisms = map[string]bool{
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SSH":   true,
	"TLS":   true,
	"TTL":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
}
