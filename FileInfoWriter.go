package main

import (
	"log"
	"os"
)

type FileInfoWriter struct {
	fileName  string
	filesChan chan string
}

func (fw *FileInfoWriter) Create() {
	fw.filesChan = make(chan string)
}

func (fw *FileInfoWriter) Close() {
	close(fw.filesChan)
}

func (fw *FileInfoWriter) PutFileName(fileName string) {
	fw.filesChan <- fileName
}

func (fw *FileInfoWriter) WriteInfo(fsi *FilesStatInfo) {
	file, err := os.Create(fw.fileName)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	for folderName := range fw.filesChan {
		_, err := file.WriteString(folderName + "\n")
		if err != nil {
			return
		}
		fsi.IncreaseCount()
	}
}
