package files

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/slintes/sort-imports/imports"
)

type MyFile struct {
	Path            string
	OwnModule       string
	Imports         []*imports.MyImport
	SortedImports   string
	UnsortedImports string
	NewFile         string
}

func (f *MyFile) Parse() error {
	//fmt.Printf("handling %s\n", myFile.Path)

	f.Imports = make([]*imports.MyImport, 0)

	file, err := os.Open(f.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	inImports := false
	inMultiLineComment := false
	var lastImport *imports.MyImport
	lastComment := ""

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		t := scanner.Text()

		//fmt.Println(t)

		if strings.HasPrefix(t, "import (") {
			//log.Println("start imports")
			inImports = true
			f.NewFile += t + "\n"
			continue
		}

		if inImports {
			if strings.HasPrefix(t, ")") {
				//log.Println("stop imports")
				inImports = false
				if len(lastComment) > 0 && lastImport != nil {
					lastImport.AfterComment += lastComment
					lastComment = ""
				}
				f.SortedImports = imports.SortImports(f.Imports)
				f.NewFile += f.SortedImports
				f.NewFile += t + "\n"
				continue
			}

			f.UnsortedImports += t + "\n"

			trimmed := strings.TrimSpace(t)

			// comments...
			// when directly above an import, keep it there
			// when before an empty line, keep after last import
			commentHandled := false
			if strings.HasPrefix(trimmed, "//") && !inMultiLineComment {
				lastComment += t + "\n"
				commentHandled = true
			}
			if strings.HasPrefix(trimmed, "/*") {
				inMultiLineComment = true
				commentHandled = true
			}
			if inMultiLineComment {
				lastComment += t + "\n"
				commentHandled = true
			}
			if strings.HasSuffix(trimmed, "*/") {
				inMultiLineComment = false
				commentHandled = true
			}
			if commentHandled {
				continue
			}

			// handle empty line for comment placing
			if len(trimmed) == 0 {
				if len(lastComment) > 0 && lastImport != nil {
					// attach comment to last import
					lastImport.AfterComment += lastComment
					lastComment = ""
				}
				continue
			}

			// everything else is an import
			//log.Println("add import")
			lastImport = imports.ParseImport(t, f.OwnModule)
			// add comment
			if len(lastComment) > 0 {
				lastImport.BeforeComment += lastComment
				lastComment = ""
			}
			f.Imports = append(f.Imports, lastImport)
			continue
		}

		f.NewFile += t + "\n"
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func (f *MyFile) Diff() bool {
	if f.SortedImports != f.UnsortedImports {
		fmt.Printf("FILE: %s\nUNSORTED:\n%s\nSORTED:\n%s\n\n", f.Path, f.UnsortedImports, f.SortedImports)
		return true
	}
	return false
}

func (f *MyFile) Write() error {
	file, err := os.Create(f.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	_, err = w.WriteString(f.NewFile)
	if err != nil {
		return err
	}
	err = w.Flush()
	if err != nil {
		return err
	}

	return nil
}
