package main

import (
	"encoding/json"
	"fmt"
	"os"
	"syscall"
)

func main() {

	if len(os.Args) < 2 {
		os.Exit(4)
	}

	path := os.Args[1]
	scan := os.Getenv("SCAN")
	if scan == "" {
		scan = "DEFAULT"
	}

	file, err := os.Stat(path)
	if err != nil {
		os.Exit(8)
	}
	mode := file.Mode().Perm()
	octalMode := fmt.Sprintf("%o", mode)
	var inode uint64

	file.Size()

	if stat, ok := file.Sys().(*syscall.Stat_t); ok {
		inode = stat.Ino
	} else {
		os.Exit(1)

	}

	shaFile, err := calculateFileHash(path)
	if err != nil {
		os.Exit(2)
	}
	secinfo := populateFileStruct(file, path, scan, inode, shaFile, octalMode)
	jsonString, err := json.Marshal(secinfo)
	if err != nil {
		os.Exit(3)
	}

	fmt.Println(string(jsonString))
}
