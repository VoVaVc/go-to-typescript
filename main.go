package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	goParser "github.com/zpatrick/go-parser"
	"golang.org/x/exp/slices"
)

func find(root, ext string) []string {
	var a []string
	filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if filepath.Ext(d.Name()) == ext {
			a = append(a, s)
		}
		return nil
	})
	return a
}

type Params struct {
	Imports     []string
	Structs     []string
	StructNames []string
	Envs        EnvVars
}

type EnvVars struct {
	Output            string
	Prefix            string
	Suffix            string
	Indent            string
	CreateFromMethod  bool
	CreateConstructor bool
	BackupDir         string
	DontExport        bool
	CreateInterface   bool
}

func env(name string) string {
	str := os.Getenv(name)
	if len(str) < 1 {
		return "\"\""
	} else {
		return fmt.Sprintf("\"%s\"", str)
	}
}

func main() {
	var params Params
	ident := os.Getenv("IDENT")
	if len(ident) < 1 {
		// set default identation to tab
		ident = "\"\t\""
	}

	envs := EnvVars{
		Output:            os.Getenv("OUTPUT_FILE"),
		Prefix:            env("PREFIX"),
		Suffix:            env("SUFFIX"),
		Indent:            ident,
		CreateFromMethod:  os.Getenv("CREATE_FROM_METHOD") == "true",
		CreateConstructor: os.Getenv("CREATE_CONSTRUCTOR") == "true",
		BackupDir:         env("BACKUP_DIR"),
		DontExport:        os.Getenv("DONT_EXPORT") == "true",
		CreateInterface:   os.Getenv("CREATE_INTERFACE") == "true",
	}
	files := find(os.Getenv("INPUT_FOLDER"), ".go")
	goFiles, err := goParser.ParseFiles(files)

	if err != nil {
		log.Fatal(err)
	}

	var Structs []string
	var StructNames []string
	var Imports []string

	for _, goFile := range goFiles {
		for _, goImport := range goFile.Imports {
			if !slices.Contains(Imports, goImport.Path) {
				Imports = append(Imports, goImport.Path)
			}
		}

		for _, goStruct := range goFile.Structs {
			StructNames = append(StructNames, goStruct.Name)
			typeStruct := fmt.Sprintf("type %s struct { \n", goStruct.Name)

			for _, goField := range goStruct.Fields {
				typeStruct += fmt.Sprintf("%s %s %v \n", goField.Name, goField.Type, goField.Tag.Value)
			}

			typeStruct += "} \n\n"
			Structs = append(Structs, typeStruct)
		}
	}

	params.Imports = Imports
	params.Structs = Structs
	params.StructNames = StructNames
	params.Envs = envs

	const generateFunc = `
		package main

		import (
			"fmt"
			"github.com/tkrajina/typescriptify-golang-structs/typescriptify"
			{{ range .Imports }}	{{ . }}
			{{ end }}
		)

		{{ range .Structs }}	{{ . }}
		{{ end }}

		func main() {
			t := typescriptify.New()
			t.Suffix = {{.Envs.Suffix}}
			t.Prefix = {{.Envs.Prefix}}
			t.Indent = {{.Envs.Indent}}
			t.CreateFromMethod = {{.Envs.CreateFromMethod}}
			t.CreateConstructor = {{.Envs.CreateConstructor}}
			t.DontExport = {{.Envs.DontExport}}
			t.CreateInterface = {{.Envs.CreateInterface}}

			{{ range .StructNames }}	t.Add({{ . }}{})
			{{ end }}
			err := t.ConvertToFile("{{.Envs.Output}}.ts")
			if err != nil {
				panic(err.Error())
			}
			fmt.Println("OK")
		}
	`

	populatedFunc := template.Must(template.New("").Parse(generateFunc))
	var tpl bytes.Buffer
	if err := populatedFunc.Execute(&tpl, params); err != nil {
		log.Fatal(err)
	}

	const tempFileName = "generate-types.go"

	file, err := os.Create(tempFileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	l, err := file.WriteString(tpl.String())
	if err != nil {
		fmt.Println(err)
		file.Close()
		return
	}

	fmt.Println(l, "bytes written successfully")
	err = file.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	cmd := exec.Command("go", "run", tempFileName)
	fmt.Println(strings.Join(cmd.Args, " "))
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		log.Panic(err)
	}
	fmt.Println(string(output))
}
