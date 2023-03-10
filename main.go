package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/slintes/sort-imports/files"
)

func main() {

	var overwriteFile = flag.Bool("w", false, "overwrite files with sorted imports")
	flag.Parse()

	if flag.NArg() != 1 {
		log.Fatal("usage: sort-imports [-w] <project_dir>")
	}

	root := flag.Arg(0)

	ownModule, err := getOwnModule(root)
	if err != nil {
		log.Fatalf("could no determine own project's package: %v", err)
	}

	hasDiff := false

	err = filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if info.Name() == "vendor" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		myFile := &files.MyFile{
			Path:      path,
			OwnModule: ownModule,
		}
		err = myFile.Parse()
		if err != nil {
			return err
		}

		thisHasDiff := myFile.Diff()

		// TODO is this thread safe?
		hasDiff = hasDiff || thisHasDiff

		if *overwriteFile {
			err = myFile.Write()
			if err != nil {
				return err
			}
			return nil
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	if hasDiff && !*overwriteFile {
		os.Exit(2)
	}
}

func getOwnModule(root string) (string, error) {
	goMod, err := os.Open(root + "/go.mod")
	if err != nil {
		return "", err
	}
	defer goMod.Close()

	scanner := bufio.NewScanner(goMod)
	for scanner.Scan() {
		t := scanner.Text()
		if strings.HasPrefix(t, "module") {
			parts := strings.Split(t, " ")
			if len(parts) != 2 {
				return "", fmt.Errorf("unexpected module definition in go.mod: %s", t)
			}
			return parts[1], nil
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", fmt.Errorf("no module definiton found in go.mod")
}
