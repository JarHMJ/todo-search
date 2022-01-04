package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var SkipFiles = map[string]struct{}{
	"vendor":  {},
	".idea":   {},
	".git":    {},
	".github": {},
}

func isScanFile(name string) bool {
	return strings.HasSuffix(name, ".go") || strings.HasSuffix(name, ".py")
}

func ScanFile(file string) error {
	fmt.Printf("start scan file: %+v \n", file)
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(f)
	lineNum := 1
	reg := regexp.MustCompile(`(?i)^(//|#).*\b(FIXME|TODO)\b.*`)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if reg.MatchString(line) {
			fmt.Printf("file:%+v line:%d \n%s \n", file, lineNum, line)
		}
		lineNum++
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func TraverseDir(dir string) error {
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		fileName := d.Name()

		if d.IsDir() {
			if _, ok := SkipFiles[fileName]; ok {
				return filepath.SkipDir
			}
		} else {
			if isScanFile(fileName) {
				return ScanFile(path)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func main() {
	err := TraverseDir(".")
	if err != nil {
		fmt.Println(err.Error())
	}
}
