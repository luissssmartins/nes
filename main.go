package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/fogleman/nes/ui"
)

func main() {
	log.SetFlags(0)
	paths := getPaths()
	if len(paths) == 0 {
		log.Fatalln("no rom files specified or found")
	}
	ui.Run(paths)
}

func getPaths() []string {
	args := os.Args[1:]
	if len(args) == 0 {
		wd, _ := os.Getwd()
		args = []string{wd}
	}
	var result []string
	for _, arg := range args {
		if isURL(arg) {
			p, err := downloadROM(arg)
			if err != nil {
				log.Println(err)
				continue
			}
			result = append(result, p)
		} else {
			result = append(result, localPaths(arg)...)
		}
	}
	return result
}

func localPaths(arg string) []string {
	info, err := os.Stat(arg)
	if err != nil {
		return nil
	}
	if info.IsDir() {
		infos, err := ioutil.ReadDir(arg)
		if err != nil {
			return nil
		}
		var result []string
		for _, info := range infos {
			name := info.Name()
			if !strings.HasSuffix(name, ".nes") {
				continue
			}
			result = append(result, path.Join(arg, name))
		}
		return result
	} else {
		return []string{arg}
	}
}
