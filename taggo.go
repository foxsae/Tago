/*
 
 Taggo version 0.1 "Emacs etags for Go"
 Author: Alex Combas
 Website: www.goplexian.com
 Email: alex.combas@gmail.com
 
 Copyright: Alex Combas 2010
 Initial release: January 03 2010
 LICENSE: GNU GPL

 See README for usage, compiling, and other info.
 
 */

package main

import (
	"go/parser"
	"go/ast"
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
)

// Get working directory and set it for savePath flag default
func whereAmI() string {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting working directory: %s\n", err.String())
	} else {
		dir += "/"
	}
	return dir
}

// Setup flag variables
var saveDir = flag.String("d", whereAmI(), "Change save directory: -d=/path/to/my/tags/")
var tagsName = flag.String("n", "TAGS", "Change TAGS name: -n=MyTagsFile")
var appendMode = flag.Bool("a", false, "Append mode: -a")

type Tea struct { bag bytes.Buffer }

func (t *Tea) String() string { return t.bag.String() }

func (t *Tea) Write(p []byte) (n int, err os.Error) {
	t.bag.Write(p)
	return len(p), nil
}

// Writes a TAGS line to a Tea buffer
func (t *Tea) drink(leaf *ast.Ident) {
	s := scoop(leaf.Position.Filename, leaf.Position.Line)
	fmt.Fprintf(t, "%s%s%d,%d\n", s, leaf.Value, leaf.Position.Line, leaf.Position.Offset)
}

// TAGS file is either appended or created, not overwritten.
func (t *Tea) savor() {
	location := fmt.Sprintf("%s%s", *saveDir, *tagsName)
	if *appendMode {
		file, err := os.Open(location, os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			fmt.Printf("Error appending file \"%s\": %s", location, err.String())
		} else {
			b := t.bag.Len()
			file.WriteAt(t.bag.Bytes(), int64(b))
			file.Close()
		}
	} else {
		
		file, err := os.Open(location, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0666)
		if err != nil {
			fmt.Printf("Error writing file \"%s\": %s\n",location, err.String())
			fmt.Println("Hint: taggo will not overwrite an existing tagsfile, only create or append.")
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
	file, err := os.Open(name, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Printf("Error opening file: %s\n", err.String())
	}
	r := bufio.NewReader(file)
	
	// iterate until reaching line #n
	for i := 1 ; i <= n; i++ { 
		sought, err := r.ReadBytes(newline)
		if err != nil {
			fmt.Printf("Error reading bytes: %s\n", err.String())
		}
		line = sought[0:(len(sought)-1)] //strip the newline
	}
	file.Close()
	return line
}

// Parses the source files given on the commandline, returns a TAGS chunk for each file
func brew() string {
	teaPot := new(Tea)
	for i := 0 ; i < len(flag.Args()) ; i++ {
		teaCup := new(Tea)
		ptree, perr := parser.ParseFile(flag.Arg(i), nil, 0)

		// return an empty string if there are any parsing errors.
		if perr != nil {
			fmt.Println("Error parsing file: ", perr.String())
			return ""
		}
		
		// if there were no parsing errors then process normally
		for _, l := range ptree.Decls {
			switch leaf := l.(type) {
			case *ast.FuncDecl:
				teaCup.drink(leaf.Name)
			case *ast.GenDecl:
				for _, c := range leaf.Specs {
					switch cell := c.(type) {
					case *ast.TypeSpec:
						teaCup.drink(cell.Name)
					case *ast.ValueSpec:
						for _, atom := range cell.Names {
							teaCup.drink(atom)
						}
					}
				}
			}
		}
		totalBytes := teaCup.bag.Len()
		fmt.Fprintf(teaPot, "\n%s,%d\n%s", ptree.Position.Filename, totalBytes, teaCup)
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
