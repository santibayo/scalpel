package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"syscall"
	"time"
)

type FileSecInfo struct {
	Name       string `json:"name"`
	Inode      uint64 `json:"inode"`
	Scan       string `json:"scan"`
	Fpath      string `json:"fpath"`
	Sha        string `json:"sha"`
	Timemod    int64  `json:"t_mod"`
	Timecreate int64  `json:"t_crt"`
	Mode       string `json:"mode"`
	TSize      int64  `json:"size"`
	Uid        uint32 `json:"uid"`
	Xdev       uint64 `json:"dev"`
	Gid        uint32 `json:"gid"`
	TimeAccess int64  `json:"t_acc"`
}

func calculateFileHash(filePath string) (string, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		os.Exit(4)
	}
	defer file.Close()

	// Create a new SHA-256 hash object
	hasher := sha256.New()

	// Copy the file content to the hash object
	_, err = io.Copy(hasher, file)
	if err != nil {
		os.Exit(5)
	}

	// Get the final hash as a byte slice
	hashBytes := hasher.Sum(nil)

	// Convert the hash to a hexadecimal string
	hashString := hex.EncodeToString(hashBytes)

	return hashString, nil
}

func main() {

	//outputFile, err := os.Create("output_and_errors.txt")
	//if err != nil {
	//	os.Exit(9)
	//	
	//}
	//defer outputFile.Close()
	//os.Stderr = outputFile

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

	// Check if the file system supports syscall.Stat_t information
	if stat, ok := file.Sys().(*syscall.Stat_t); ok {
		inode = stat.Ino
	} else {
		os.Exit(1)

	}

	shaFile, err := calculateFileHash(path)
	if err != nil {
		os.Exit(2)
	}
	owner := file.Sys().(*syscall.Stat_t).Uid
	gid := file.Sys().(*syscall.Stat_t).Gid
	dev := file.Sys().(*syscall.Stat_t).Dev
	creationTime := file.Sys().(*syscall.Stat_t).Ctim
	modificationTime := file.ModTime().UnixNano() / 1000
	unix_creation := creationTime.Nano() / 1000
	lastAccessTime := time.Unix(file.Sys().(*syscall.Stat_t).Atim.Unix())
	secinfo := FileSecInfo{
		Fpath:      path,
		Name:       file.Name(),
		Scan:       scan,
		Inode:      inode,
		Sha:        shaFile,
		Timemod:    modificationTime / 1000,
		Timecreate: unix_creation / 1000,
		Mode:       octalMode,
		TSize:      file.Size(),
		Uid:        owner,
		Gid:        gid,
		Xdev:       dev,
		TimeAccess: lastAccessTime.UnixMilli(),
	}

	jsonString, err := json.Marshal(secinfo)
	if err != nil {
		os.Exit(3)
	}

	fmt.Println(string(jsonString))
}
