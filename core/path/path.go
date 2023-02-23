package path

import (
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

var rw sync.RWMutex
var basePath string

// 返回二進制程序所在的路徑
func BasePath() string {
	rw.RLock()
	val := basePath
	rw.RUnlock()
	if val != `` {
		return val
	}
	rw.Lock()
	val = basePath
	if val != `` {
		rw.Unlock()
		return val
	}
	defer rw.Unlock()

	filename, e := exec.LookPath(os.Args[0])
	if e != nil {
		panic(e)
	}
	filename, e = filepath.Abs(filename)
	if e != nil {
		panic(e)
	}
	val = filepath.Dir(filename)
	basePath = val
	return val
}

// 如果 path 是相對路徑，則返回它相對 basePath 的路徑
func Abs(bashPath, path string) string {
	if filepath.IsAbs(path) {
		path = filepath.Clean(path)
	} else {
		path = filepath.Clean(filepath.Join(bashPath, path))
	}
	return path
}
