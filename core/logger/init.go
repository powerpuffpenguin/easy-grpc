package logger

import (
	"path/filepath"

	"github.com/powerpuffpenguin/easy-grpc/core/cnf"
)

// 日誌單例
var Logger Wrap

// 初始化日誌
func Init(basePath string, options *cnf.Logger) {
	if options.Filename == `` {
		options.Filename = filepath.Clean(filepath.Join(basePath, `var`, `logs`, `server.log`))
	} else {
		if filepath.IsAbs(options.Filename) {
			options.Filename = filepath.Clean(options.Filename)
		} else {
			options.Filename = filepath.Clean(filepath.Join(basePath, options.Filename))
		}
	}
	Logger.Attach(New(options))
}
