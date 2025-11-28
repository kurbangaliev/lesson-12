package main

import (
	"errors"
	"fmt"
	"os"
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

func getDirectoryScan() (string, error) {
	if len(os.Args) > 1 {
		firstArg := os.Args[1]
		return firstArg, nil
	} else {
		return "", errors.New("No argument. Please input path of scan directory.")
	}
}

func main() {
	var fileInfoWriter = FileInfoWriter{fileName: "files.txt"}
	var folderInfoWriter = FileInfoWriter{fileName: "folders.txt"}
	var scanner DirectoryScanner

	fileInfoWriter.Create()
	folderInfoWriter.Create()
	defer fileInfoWriter.Close()
	defer folderInfoWriter.Close()

	scanPath, err := getDirectoryScan()
	if err != nil {
		fmt.Println(err)
		return
	}

	go fileInfoWriter.WaitInfo()
	go folderInfoWriter.WaitInfo()

	var wg sync.WaitGroup
	wg.Add(1)
	go scanner.Scan(scanPath, &wg, &folderInfoWriter, &fileInfoWriter)
	wg.Wait()

	fmt.Printf("Total files: %d\n", fileInfoWriter.fileStatInfo.count)
	fmt.Printf("Total folders: %d\n", folderInfoWriter.fileStatInfo.count)
}
