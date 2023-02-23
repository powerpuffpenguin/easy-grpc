package cnf

import (
	"strings"
	"time"
)

// 數據庫設定
type DB struct {
	Driver           string
	Source           []string
	ShowSQL          bool
	Cache            DBCache
	MaxOpen, MaxIdle int
}

// 數據庫緩存
type DBCache struct {
	// 存儲後端
	Redis Redis
	// 時間緩存
	Modtime Redis
	// 禁用緩存列表
	Direct []string
	// 表獨立緩存
	Special []struct {
		// 表名
		Name string
		// 存儲後端
		Redis Redis
	}
}
type Redis struct {
	Write string
	Read  string
	// 緩存過期時間
	Timeout time.Duration
}

func (r *Redis) IsValid() bool {
	return strings.HasPrefix(r.Write, "redis://")
}
func (r *Redis) Key() string {
	if strings.HasPrefix(r.Read, `redis://`) {
		return r.Timeout.String() + ` ` + r.Write + ` ` + r.Read
	} else {
		return r.Timeout.String() + ` ` + r.Write
	}
}
