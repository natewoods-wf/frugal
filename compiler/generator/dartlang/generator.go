package dartlang

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"gopkg.in/yaml.v2"

	"github.com/Workiva/frugal/compiler/generator"
	"github.com/Workiva/frugal/compiler/globals"
	"github.com/Workiva/frugal/compiler/parser"
)

const (
	lang               = "dart"
	defaultOutputDir   = "gen-dart"
	minimumDartVersion = "1.12.0"
	tab                = "  "
	tabtab             = tab + tab
	tabtabtab          = tab + tab + tab
)

type Generator struct {
	*generator.BaseGenerator
}

func NewGenerator() generator.MultipleFileGenerator {
	return &Generator{&generator.BaseGenerator{}}
}

func (g *Generator) GetOutputDir(dir string, p *parser.Program) string {
	if pkg, ok := p.Namespaces[lang]; ok {
		path := generator.GetPackageComponents(pkg)
		dir = filepath.Join(append([]string{dir}, path...)...)
	} else {
		dir = filepath.Join(dir, p.Name)
	}
	return dir
}

func (g *Generator) DefaultOutputDir() string {
	return defaultOutputDir
}

func (g *Generator) GenerateDependencies(p *parser.Program, dir string) error {
	if err := g.addToPubspec(p, dir); err != nil {
		return err
	}
	if err := g.exportClasses(p, dir); err != nil {
		return err
	}
	return nil
}

type pubspec struct {
	Name         string `yaml:"name"`
	Version      string `yaml:"version"`
	Description  string `yaml:"description"`
	Environment  env    `yaml:"environment"`
	Dependencies deps   `yaml:"dependencies"`
}

type env struct {
	SDK string `yaml:"sdk"`
}

type deps struct {
	Thrift dep `yaml:"thrift"`
	Frugal dep `yaml:"frugal"`
}

type dep struct {
	Git gitDep `yaml:"git"`
}

type gitDep struct {
	URL string `yaml:"url"`
}

func (g *Generator) addToPubspec(p *parser.Program, dir string) error {
	pubFilePath := filepath.Join(dir, "pubspec.yaml")
	ps := &pubspec{
		Name:        strings.ToLower(p.Name),
		Version:     globals.Version,
		Description: "Autogenerated by the frugal compiler",
		Environment: env{
			SDK: "^" + minimumDartVersion,
		},
		Dependencies: deps{
			Thrift: dep{
				Git: gitDep{
					URL: "git@github.com:Workiva/thrift-dart.git",
				},
			},
			Frugal: dep{
				Git: gitDep{
					URL: "git@github.com:Workiva/frugal-dart.git",
				},
			},
		},
	}

	d, err := yaml.Marshal(&ps)
	if err != nil {
		return err
	}
	// create and write to new file
	newPubFile, err := os.Create(pubFilePath)
	defer newPubFile.Close()
	if err != nil {
		return nil
	}
	if _, err := newPubFile.Write(d); err != nil {
		return err
	}
	return nil
}

func (g *Generator) exportClasses(p *parser.Program, dir string) error {
	dartFile := fmt.Sprintf("%s.%s", strings.ToLower(p.Name), lang)
	mainFilePath := filepath.Join(dir, "lib", dartFile)
	mainFile, err := os.OpenFile(mainFilePath, syscall.O_RDWR, 0777)
	defer mainFile.Close()
	if err != nil {
		return err
	}

	exports := "\n"
	for _, scope := range p.Scopes {
		exports += fmt.Sprintf("export 'src/%s%s.%s' show %sPublisher, %sSubscriber;\n",
			generator.FilePrefix, strings.ToLower(scope.Name), lang, scope.Name, scope.Name)
	}
	stat, err := mainFile.Stat()
	if err != nil {
		return err
	}
	_, err = mainFile.WriteAt([]byte(exports), stat.Size())
	return err
}

func (g *Generator) CheckCompile(path string) error {
	// TODO: Add compile to js
	return nil
}

func (g *Generator) GenerateFile(name, outputDir string, fileType generator.FileType) (*os.File, error) {
	if fileType != generator.CombinedFile {
		return nil, fmt.Errorf("frugal: Bad file type for dartlang generator: %s", fileType)
	}
	outputDir = filepath.Join(outputDir, "lib")
	outputDir = filepath.Join(outputDir, "src")
	return g.CreateFile(name, outputDir, lang, true)
}

func (g *Generator) GenerateDocStringComment(file *os.File) error {
	comment := fmt.Sprintf(
		"// Autogenerated by Frugal Compiler (%s)\n"+
			"// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING",
		globals.Version)

	_, err := file.WriteString(comment)
	return err
}

func (g *Generator) GeneratePackage(file *os.File, p *parser.Program, scope *parser.Scope) error {
	// TODO: Figure out what this does
	pkg, ok := p.Namespaces[lang]
	if ok {
		components := generator.GetPackageComponents(pkg)
		pkg = components[len(components)-1]
	} else {
		pkg = p.Name
	}
	_, err := file.WriteString(fmt.Sprintf("library %s.src.%s%s;", pkg,
		generator.FilePrefix, strings.ToLower(scope.Name)))
	return err
}

func (g *Generator) GenerateImports(file *os.File, scope *parser.Scope) error {
	imports := "import 'dart:async';\n\n"
	imports += "import 'package:thrift/thrift.dart' as thrift;\n"
	imports += "import 'package:frugal/frugal.dart' as frugal;\n\n"
	params := map[string]bool{}
	for _, op := range scope.Operations {
		if _, ok := params[op.Param]; !ok {
			params[op.Param] = true
		}
	}
	for key, _ := range params {
		lowerKey := strings.ToLower(key)
		imports += fmt.Sprintf("import '%s.dart' as t_%s;\n", lowerKey, lowerKey)
	}
	_, err := file.WriteString(imports)
	return err
}

func (g *Generator) GenerateConstants(file *os.File, name string) error {
	constants := fmt.Sprintf("const String delimiter = '%s';", globals.TopicDelimiter)
	_, err := file.WriteString(constants)
	return err
}

func (g *Generator) GeneratePublisher(file *os.File, scope *parser.Scope) error {
	publishers := ""
	publishers += fmt.Sprintf("class %sPublisher {\n", scope.Name)
	publishers += tab + "frugal.Transport transport;\n"
	publishers += tab + "thrift.TProtocol protocol;\n"
	publishers += tab + "int seqId;\n\n"

	publishers += fmt.Sprintf(tab+"%sPublisher(frugal.Provider provider) {\n", scope.Name)
	publishers += tabtab + "var tp = provider.newTransportProtocol();\n"
	publishers += tabtab + "transport = tp.transport;\n"
	publishers += tabtab + "protocol = tp.protocol;\n"
	publishers += tabtab + "seqId = 0;\n"
	publishers += tab + "}\n\n"

	args := ""
	if len(scope.Prefix.Variables) > 0 {
		for _, variable := range scope.Prefix.Variables {
			args = fmt.Sprintf("%sString %s, ", args, variable)
		}
	}
	prefix := ""
	for _, op := range scope.Operations {
		publishers += prefix
		prefix = "\n\n"
		publishers += fmt.Sprintf(tab+"Future publish%s(%st_%s.%s req) {\n", op.Name, args,
			strings.ToLower(op.Param), op.Param)
		publishers += fmt.Sprintf(tabtab+"var op = \"%s\";\n", op.Name)
		publishers += fmt.Sprintf(tabtab+"var prefix = \"%s\";\n", generatePrefixStringTemplate(scope))
		publishers += tabtab + "var topic = \"${prefix}" + scope.Name + "${delimiter}${op}\";\n"
		publishers += tabtab + "transport.preparePublish(topic);\n"
		publishers += tabtab + "var oprot = protocol;\n"
		publishers += tabtab + "seqId++;\n"
		publishers += tabtab + "var msg = new thrift.TMessage(op, thrift.TMessageType.CALL, seqId);\n"
		publishers += tabtab + "oprot.writeMessageBegin(msg);\n"
		publishers += tabtab + "req.write(oprot);\n"
		publishers += tabtab + "oprot.writeMessageEnd();\n"
		publishers += tabtab + "return oprot.transport.flush();\n"
		publishers += tab + "}\n"
	}

	publishers += "}\n"

	_, err := file.WriteString(publishers)
	return err
}

func generatePrefixStringTemplate(scope *parser.Scope) string {
	if scope.Prefix.String == "" {
		return ""
	}
	template := ""
	template += scope.Prefix.Template()
	template += globals.TopicDelimiter
	if len(scope.Prefix.Variables) == 0 {
		return template
	}
	vars := make([]interface{}, len(scope.Prefix.Variables))
	for i, variable := range scope.Prefix.Variables {
		vars[i] = fmt.Sprintf("${%s}", variable)
	}
	template = fmt.Sprintf(template, vars...)
	return template
}

func (g *Generator) GenerateSubscriber(file *os.File, scope *parser.Scope) error {
	subscribers := ""
	subscribers += fmt.Sprintf("class %sSubscriber {\n", scope.Name)
	subscribers += tab + "frugal.Provider provider;\n\n"

	subscribers += fmt.Sprintf(tab+"%sSubscriber(frugal.Provider provider) {\n", scope.Name)
	subscribers += tabtab + "this.provider = provider;\n"
	subscribers += tab + "}\n\n"

	args := ""
	if len(scope.Prefix.Variables) > 0 {
		for _, variable := range scope.Prefix.Variables {
			args = fmt.Sprintf("%sString %s, ", args, variable)
		}
	}
	prefix := ""
	for _, op := range scope.Operations {
		paramLower := strings.ToLower(op.Param)
		subscribers += prefix
		prefix = "\n\n"
		subscribers += fmt.Sprintf(tab+"Future<frugal.Subscription> subscribe%s(%sdynamic on%s(t_%s.%s req)) async {\n",
			op.Name, args, op.Param, paramLower, op.Param)
		subscribers += fmt.Sprintf(tabtab+"var op = \"%s\";\n", op.Name)
		subscribers += fmt.Sprintf(tabtab+"var prefix = \"%s\";\n", generatePrefixStringTemplate(scope))
		subscribers += tabtab + "var topic = \"${prefix}" + scope.Name + "${delimiter}${op}\";\n"
		subscribers += tabtab + "var tp = provider.newTransportProtocol();\n"
		subscribers += tabtab + "await tp.transport.subscribe(topic);\n"
		subscribers += tabtab + "tp.transport.signalRead.listen((_) {\n"
		subscribers += fmt.Sprintf(tabtabtab+"on%s(_recv%s(op, tp.protocol));\n", op.Param, op.Name)
		subscribers += tabtab + "});\n"
		subscribers += tabtab + "var sub = new frugal.Subscription(topic, tp.transport);\n"
		subscribers += tabtab + "tp.transport.error.listen((Error e) {;\n"
		subscribers += tabtabtab + "sub.signal(e);\n"
		subscribers += tabtab + "});\n"
		subscribers += tabtab + "return sub;\n"
		subscribers += tab + "}\n\n"

		subscribers += fmt.Sprintf(tab+"t_%s.%s _recv%s(String op, thrift.TProtocol iprot) {\n",
			paramLower, op.Param, op.Name)
		subscribers += tabtab + "var tMsg = iprot.readMessageBegin();\n"
		subscribers += tabtab + "if (tMsg.name != op) {\n"
		subscribers += tabtabtab + "thrift.TProtocolUtil.skip(iprot, thrift.TType.STRUCT);\n"
		subscribers += tabtabtab + "iprot.readMessageEnd();\n"
		subscribers += tabtabtab + "throw new thrift.TApplicationError(\n"
		subscribers += tabtabtab + "thrift.TApplicationErrorType.UNKNOWN_METHOD, tMsg.name);\n"
		subscribers += tabtab + "}\n"
		subscribers += fmt.Sprintf(tabtab+"var req = new t_%s.%s();\n", paramLower, op.Param)
		subscribers += tabtab + "req.read(iprot);\n"
		subscribers += tabtab + "iprot.readMessageEnd();\n"
		subscribers += tabtab + "return req;\n"
		subscribers += tab + "}\n"
	}

	subscribers += "}\n"

	_, err := file.WriteString(subscribers)
	return err
}
