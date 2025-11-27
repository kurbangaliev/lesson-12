package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type DirectoryScanner struct{}

func (ds *DirectoryScanner) Scan(filePath string, wg *sync.WaitGroup, folderWriter *FileInfoWriter, filesWriter *FileInfoWriter) {
	defer wg.Done()

	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("File '%s' does not exist.\n", filePath)
		} else {
			fmt.Printf("Error getting file info for '%s': %v\n", filePath, err)
		}
		return
	}

	entries, err := os.ReadDir(filePath)
	if err != nil {
		log.Fatalf("Error reading directory: %v", err)
	}

	fmt.Printf("Contents of directory '%s':\n", filePath)
	for _, entry := range entries {
		if entry.IsDir() {
			folderWriter.PutFileName(entry.Name())
			wg.Add(1)
			go ds.Scan(filepath.Join(filePath, entry.Name()), wg, folderWriter, filesWriter)
		} else {
			filesWriter.PutFileName(entry.Name())
		}
	}
}
