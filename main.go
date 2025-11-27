package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
)

/*
Реализовать на Go модуль, осуществляющий обход директории, переданной программе в качестве аргумента.
Во время обхода директории Все файлы записываются в файл filse.txt -> директории -> folders.txt | имена файлов
1) Go рутины для самых уверенных
2) аргумент можно принять os.Args[1]// main.exe c:\users\user\Desktop\
3) Для обхода воспользуйтесь filepath.WalkDir()
4) file.Stat() -> IsDir()?
*/

type FoldersStatInfo struct {
	mu           sync.Mutex
	filesCount   int
	foldersCount int
}

func directoryScanner(filePath string, wg *sync.WaitGroup, foldersChan chan string, filesChan chan string) {
	fmt.Println("Scanning", filePath)
	//	wg.Add(1)
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
			fmt.Printf("  [DIR] %s\n", entry.Name())
			foldersChan <- entry.Name()
			wg.Add(1)
			go directoryScanner(filepath.Join(filePath, entry.Name()), wg, foldersChan, filesChan)
		} else {
			filesChan <- entry.Name()
			fmt.Printf("  [FILE] %s\n", entry.Name())
		}
	}
}

func writeFiles(filesChan chan string, foldersStatInfo *FoldersStatInfo) {
	file, err := os.OpenFile("files.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()
	for fileName := range filesChan {
		_, err := file.WriteString(fileName + "\n")
		if err != nil {
			return
		}
		foldersStatInfo.mu.Lock()
		foldersStatInfo.filesCount++
		foldersStatInfo.mu.Unlock()
	}
}

func writeFolders(foldersChan chan string, foldersStatInfo *FoldersStatInfo) {
	file, err := os.OpenFile("folders.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()
	for folderName := range foldersChan {
		_, err := file.WriteString(folderName + "\n")
		if err != nil {
			return
		}
		foldersStatInfo.mu.Lock()
		foldersStatInfo.foldersCount++
		foldersStatInfo.mu.Unlock()
	}
}

func main() {
	var foldersStatInfo FoldersStatInfo
	scanPath, err := getDirectoryScan()
	if err != nil {
		fmt.Println(err)
		return
	}

	folders := make(chan string)
	files := make(chan string)

	defer close(files)
	defer close(folders)

	go writeFiles(files, &foldersStatInfo)
	go writeFolders(folders, &foldersStatInfo)

	var wg sync.WaitGroup
	wg.Add(1)
	go directoryScanner(scanPath, &wg, folders, files)
	wg.Wait()

	fmt.Printf("Total files: %d\n", foldersStatInfo.filesCount)
	fmt.Printf("Total folders: %d\n", foldersStatInfo.foldersCount)
}

func getDirectoryScan() (string, error) {
	if len(os.Args) > 1 {
		firstArg := os.Args[1]
		return firstArg, nil
	} else {
		return "", errors.New("No argument. Please input path of scan directory.")
	}
}
