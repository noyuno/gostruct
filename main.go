package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

var (
	namereg   = regexp.MustCompile(`^type\s+([A-Za-z][A-Za-z0-9]*)\s+struct\s+{$`)
	memberreg = regexp.MustCompile(`^\s*([A-Za-z][A-Za-z0-9]*)\s+([\[\]\*]+|)([A-Za-z][A-Za-z0-9]*)`)
	embedreg  = regexp.MustCompile(`^\s*(\*?)([A-Za-z][A-Za-z0-9]*)`)

	builtinTypes = []string{
		"bool", "byte", "complex64", "complex128", "error", "float32", "float64",
		"int", "int8", "int16", "int32", "int64", "rune", "string",
		"uint", "uint8", "uint16", "uint32", "uint64", "uintptr", "map",
	}
)

type Member struct {
	Name      string
	Type      string
	Embed     bool
	Attribute string
	Visited   bool
}

func IsBuiltinType(t string) bool {
	for i := range builtinTypes {
		if builtinTypes[i] == t {
			return true
		}
	}
	return false
}

func Analyze(dirn string, builtin bool) (map[string][]Member, error) {
	dir, err := ioutil.ReadDir(dirn)
	if err != nil {
		return nil, errors.New("could not open directory\n")
	}
	list := map[string][]Member{}
	for _, file := range dir {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".go") {
			log.Debug(file.Name())
			var f *os.File
			f, err = os.Open(dirn + "/" + file.Name())
			if err != nil {
				log.Warningf("could not open file: %v\n", file.Name())
				continue
			}
			scanner := bufio.NewScanner(f)
			name := ""
			for scanner.Scan() {
				t := scanner.Text()
				if name == "" {
					if namereg.MatchString(t) {
						a := namereg.FindStringSubmatch(t)
						name = a[1]
						list[name] = []Member{}
					}
				} else {
					if memberreg.MatchString(t) {
						a := memberreg.FindStringSubmatch(t)
						//log.WithFields(log.Fields{"nobuiltin": nobuiltin}).Debug()
						if !(!builtin && IsBuiltinType(a[3])) {
							list[name] = append(list[name], Member{
								Name:      a[1],
								Attribute: a[2],
								Type:      a[3],
								Embed:     false,
								Visited:   false})
						}
					} else if embedreg.MatchString(t) {
						a := embedreg.FindStringSubmatch(t)
						if !(!builtin && IsBuiltinType(a[2])) {
							list[name] = append(list[name], Member{
								Name:      "",
								Attribute: "",
								Type:      a[2],
								Embed:     true,
								Visited:   false})
						}
					}
					if len(t) > 0 && t[0:1] == "}" {
						name = ""
					}
				}
			}
		}
	}
	return list, nil
}

func GenerateEdge(list map[string][]Member, name string, label bool) string {
	ret := ""
	if v, o := list[name]; o {
		for i, w := range v {
			if w.Visited {
				continue
			}
			list[name][i].Visited = true
			e := ""
			if w.Embed {
				e = "*"
			}
			//fmt.Println(w.Type)
			ret += "\t\"" + name + "\"->\"" + w.Type + "\""
			if label {
				ret += " [label = \"" + w.Name + " " + w.Attribute + e + w.Type + "\"]"
			}
			ret += ";\n"
			ret += GenerateEdge(list, w.Type, label)
		}
	}
	return ret
}

func Generate(list map[string][]Member, name string, label bool) string {
	ret := `
digraph G {
	rankdir = "LR";
	graph [fontname = "Inconsolata"];
	node [fontname = "Inconsolata"];
	edge [fontname = "Inconsolata"];
`
	ret += GenerateEdge(list, name, label)
	ret += `
}
`
	return ret
}

func main() {
	log.SetOutput(os.Stderr)
	log.SetFormatter(&log.TextFormatter{DisableSorting: false, QuoteEmptyFields: true})
	var directory string
	var debug bool
	var builtin bool
	var label bool
	flag.StringVar(&directory, "d", ".", "target directory")
	flag.BoolVar(&debug, "debug", false, "outputs debug text")
	flag.BoolVar(&builtin, "b", false, "show builtin types")
	flag.BoolVar(&builtin, "builtin", false, "show builtin types")
	flag.BoolVar(&label, "l", true, "show label")
	flag.BoolVar(&label, "label", true, "show label")
	flag.Parse()
	if debug {
		log.SetLevel(log.DebugLevel)
	}
	if len(flag.Args()) == 0 {
		os.Stderr.WriteString("argument missing. requires top struct name\n")
		os.Exit(1)
	}
	list, err := Analyze(directory, builtin)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}
	graph := Generate(list, flag.Args()[0], label)

	fmt.Println(graph)
}
