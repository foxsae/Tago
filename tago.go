/*

 Tago "Emacs etags for Go"
 Author: Alex Combas
 Website: www.goplexian.com
 Email: alex.combas@gmail.com

 Version: 0.2
 © Alex Combas 2010
 Initial release: January 03 2010

 See README for usage, compiling, and other info.

*/

package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

// Get working directory and set it for savePath flag default
func whereAmI() string {
	var r string = ""
	if dir, err := os.Getwd(); err != nil {
		fmt.Printf("Error getting working directory: %s\n", err.Error())
	} else {
		r = dir + "/"
	}
	return r
}

// Setup flag variables
var saveDir = flag.String("d", whereAmI(), "Change save directory: -d=/path/to/my/tags/")
var tagsName = flag.String("n", "TAGS", "Change TAGS name: -n=MyTagsFile")
var appendMode = flag.Bool("a", false, "Append mode: -a")

type Tea struct {
	bag bytes.Buffer
}

func (t *Tea) String() string { return t.bag.String() }

func (t *Tea) Write(p []byte) (n int, err error) {
	t.bag.Write(p)
	return len(p), nil
}

// Writes a TAGS line to a Tea buffer
func (t *Tea) drink(leaf *ast.Ident, fileSet *token.FileSet) {
	if leaf.Obj == nil {
		return
	}
	position := fileSet.Position(leaf.Pos())
	s := scoop(position.Filename, position.Line)
	fmt.Fprintf(t, "%s%s%d,%d\n", s, leaf.Obj.Name, position.Line, position.Offset)
}

// TAGS file is either appended or created, not overwritten.
func (t *Tea) savor() {
	location := fmt.Sprintf("%s%s", *saveDir, *tagsName)
	if *appendMode {
		if file, err := os.OpenFile(location, os.O_APPEND|os.O_WRONLY, 0666); err != nil {
			fmt.Printf("Error appending file \"%s\": %s\n", location, err.Error())
		} else {
			b := t.bag.Len()
			file.WriteAt(t.bag.Bytes(), int64(b))
			file.Close()
		}
	} else {
		if file, err := os.OpenFile(location, os.O_CREATE|os.O_WRONLY, 0666); err != nil {
			fmt.Printf("Error writing file \"%s\": %s\n", location, err.Error())
		} else {
			file.WriteString(t.bag.String())
			file.Close()
		}
	}
}

// Returns the full line of source on which *ast.Ident.Name appears
func scoop(name string, n int) []byte {
	var newline byte = '\n'
	var line []byte // holds a line of source code
	if file, err := os.OpenFile(name, os.O_RDONLY, 0666); err != nil {
		fmt.Printf("Error opening file: %s\n", err.Error())
	} else {
		r := bufio.NewReader(file)

		// iterate until reaching line #n
		for i := 1; i <= n; i++ {
			if sought, err := r.ReadBytes(newline); err != nil {
				fmt.Printf("Error reading bytes: %s\n", err.Error())
			} else {
				line = sought[0:(len(sought) - 1)] //strip the newline
			}
		}
		file.Close()
	}
	return line
}

// Parses the source files given on the commandline, returns a TAGS chunk for each file
func brew() string {
	teaPot := new(Tea)
	fileSet := token.NewFileSet()
	for i := 0; i < len(flag.Args()); i++ {
		teaCup := new(Tea)
		if ptree, perr := parser.ParseFile(fileSet, flag.Arg(i), nil, 0); perr != nil {
			fmt.Println("Error parsing file: ", perr.Error())
			return ""
		} else {
			// if there were no parsing errors then process normally
			for _, l := range ptree.Decls {
				switch leaf := l.(type) {
				case *ast.FuncDecl:
					teaCup.drink(leaf.Name, fileSet)
				case *ast.GenDecl:
					for _, c := range leaf.Specs {
						switch cell := c.(type) {
						case *ast.TypeSpec:
							teaCup.drink(cell.Name, fileSet)
						case *ast.ValueSpec:
							for _, atom := range cell.Names {
								teaCup.drink(atom, fileSet)
							}
						}
					}
				}
			}
			totalBytes := teaCup.bag.Len()
			position := fileSet.Position(ptree.Pos())
			fmt.Fprintf(teaPot, "\f\n%s,%d\n%s", position.Filename, totalBytes, teaCup)
		}
	}
	return teaPot.String()
}

func main() {
	flag.Parse()
	tea := new(Tea)
	fmt.Fprint(tea, brew())

	// if the string is empty there were parsing errors, abort
	if tea.String() == "" {
		fmt.Println("Parsing errors experienced, aborting...")
	} else {
		tea.savor()
	}
}
