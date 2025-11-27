package main

import "sync"

type FilesStatInfo struct {
	mu    sync.Mutex
	count int
}

func (fs *FilesStatInfo) IncreaseCount() {
	fs.mu.Lock()
	fs.count++
	fs.mu.Unlock()
}
