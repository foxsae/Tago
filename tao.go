/*
* Tao version 0.1 "Emacs etags for Go"
* Author: Alex Combas
* Website: www.goplexian.com
* Email: alex.combas@gmail.com
*
*
* Copyright: Alex Combas 2010
* License: GNU GPL
* Initial release: January 03 2010
*

This program is free software; you can redistribute it and/or
modify it under the terms of the GNU General Public License
as published by the Free Software Foundation; either version 2
of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program; if not, write to the Free Software
Foundation, Inc., 59 Temple Place - Suite 330, Boston, MA  02111-1307, USA.

*
* COMPILING:
*
* $> cd tao-0.1
* $> make 
* $> cp tao /path/to/bin
* $> make clean
*
*
*
* USAGE: 
*
* tao *.go 
* tao fileX.go fileY.go fileZ.go
*
* Tao will write a TAGS file to your present working directory.
*
* WARNING: If a TAGS file exists in the pwd then it will be overwritten.
*
* To add the TAGS file to Emacs: M+x visit-tags-table RET /path/to/TAGS RET yes
*
*
*
* TODO:
*
* Add flag support: 
* -a append to TAGS file, 
* -f specify TAGS location, 
* -h print help, 
* etc
*
*
*/

package main

import (
	"go/parser"
	"go/ast"
	"bufio"
	"bytes"
	"fmt"
	"os"
)

type Tea struct { bag bytes.Buffer }

func (t *Tea) String() string { return t.bag.String() }

func (t *Tea) Write(p []byte) (n int, err os.Error) {
	t.bag.Write(p)
	return len(p), nil
}

// Writes a TAGS line to a Tea buffer
func (t *Tea) drink(leaf *ast.Ident) {
	s := spoon(leaf.Position.Filename, leaf.Position.Line)
	fmt.Fprintf(t, "%s%s%d,%d\n", s, leaf.Value, leaf.Position.Line, leaf.Position.Offset)
}

// TAGS file is written to the working directory, it is either created or overwritten
func (t *Tea) save() {
	var location string
	
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting working directory: ", err.String())
	} else {
		location = fmt.Sprintf("%s%s", wd, "/TAGS")
	}

	file, err := os.Open(location, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Error writing file: ", err.String())
	} else {
		file.WriteString(t.bag.String())
	}
	file.Close()
}

// Returns the full line of source on which *ast.Ident.Name appears
func spoon(name string, n int) []byte {
	var newline byte = '\n'
	var line []byte // holds a line of source code
	file, err := os.Open(name, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println("Error opening file: ", err.String())
	}
	r := bufio.NewReader(file)

	// iterate until reaching line #n
	for i := 1 ; i <= n; i++ { 
	sought, err := r.ReadBytes(newline)
		if err != nil {
			fmt.Println("Error reading bytes: ", err.String())
		}
		line = sought[0:(len(sought)-1)] //strip the newline
	}
	file.Close()
	return line
}

// Parses the source files given on the commandline, returns a TAGS chunk for each file
func brew() string {
	tea := new(Tea)
	for i := 1 ; i < len(os.Args) ; i++ {
		teaPot := new(Tea)
		ptree, perr := parser.ParseFile(os.Args[i], nil, 0)
		if perr != nil {
			fmt.Println("Error parsing file: ", perr.String())
		}
		
		for _, l := range ptree.Decls {
			switch leaf := l.(type) {
			case *ast.FuncDecl:
				teaPot.drink(leaf.Name)
			case *ast.GenDecl:
				for _, c := range leaf.Specs {
					switch cell := c.(type) {
					case *ast.TypeSpec:
						teaPot.drink(cell.Name)
					case *ast.ValueSpec:
						for _, atom := range cell.Names {
							teaPot.drink(atom)
						}
					}
				}
			}
		}

		totalBytes := teaPot.bag.Len()
		fmt.Fprintf(tea, "\n%s,%d\n%s", ptree.Position.Filename, totalBytes, teaPot)
	}
	return tea.String()
}

func main() {
	tea := new(Tea)
	fmt.Fprint(tea, brew())
	tea.save()
}	
