package main

import (
	"flag"
	"fmt"
	"strings"
)

var (
	RootPath  = "."
	AutoPrune *bool
)

func __init() {
	AutoPrune = flag.Bool("auto-prune", false, "auto-prune: non interactive deletes all dependency folder")
	flag.Parse()
	path := flag.Args()
	if len(path) == 1 {
		RootPath = path[0]
	}
}

func main() {
	__init()
	paths := MapAllPaths(RootPath)
	if len(paths.Paths) < 1 {
		return
	}
	_continue := "n"
	for _, path := range paths.Paths {
		fmt.Println(path)
	}
	if !*AutoPrune {
		fmt.Printf("%s will be freed, continue? [y/N]: ", paths.Size)
		fmt.Scanln(&_continue)
	}
	if strings.ToLower(_continue) == "y" {
		DeleteAll(paths.Paths)
		fmt.Printf("Freed: %s\n", paths.Size)
	}
}
