package imports

import (
	"sort"
	"strings"
)

type MyImport struct {
	// complete line with alias
	Line string

	// actual import only
	Value string

	Priority      int
	BeforeComment string
	AfterComment  string
}

func ParseImport(line string, ownModule string) *MyImport {

	parts := strings.Split(strings.TrimSpace(line), " ")
	value := parts[len(parts)-1]

	myImport := &MyImport{
		Line:  line,
		Value: value,
	}

	if !strings.Contains(value, ".") { // internal packages
		myImport.Priority = 1
	} else if strings.Contains(value, ownModule) {
		myImport.Priority = 5
	} else if strings.Contains(value, "github.com/openshift") {
		myImport.Priority = 4
	} else if strings.Contains(value, "k8s.io") {
		myImport.Priority = 3
	} else {
		myImport.Priority = 2
	}

	return myImport
}

func SortImports(myImports []*MyImport) string {
	sort.Slice(myImports, func(i, j int) bool {
		if myImports[i].Priority == myImports[j].Priority {
			return myImports[i].Value < myImports[j].Value
		}
		return myImports[i].Priority < myImports[j].Priority
	})

	lastType := 0
	sortedImports := ""
	for i, myImport := range myImports {
		// add empty line for separation
		if i > 0 && int(myImport.Priority) > lastType {
			sortedImports += "\n"
		}
		lastType = int(myImport.Priority)

		if myImport.BeforeComment != "" {
			sortedImports += "\n"
			sortedImports += myImport.BeforeComment
		}
		sortedImports += myImport.Line + "\n"
		sortedImports += myImport.AfterComment
	}
	return sortedImports
}
