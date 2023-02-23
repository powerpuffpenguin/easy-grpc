package configure

import (
	"github.com/powerpuffpenguin/easy-grpc/core/cnf"
)

var defaultConfigure Configure

// 返回設定單例
func Default() *Configure {
	return &defaultConfigure
}

// 設定檔案 內存映射
type Configure struct {
	// http 服務
	HTTP cnf.HTTP
	// 會話
	Session cnf.Session
	// 數據庫
	DB cnf.DB
	// 日誌
	Logger cnf.Logger
}

// 加載設定
func (c *Configure) Load(filename string) (e error) {
	e = cnf.Load(filename, c)
	if e != nil {
		return
	}
	return
}
