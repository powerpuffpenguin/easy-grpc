package cnf

import "time"

// 定義會話
type Session struct {
	// 會話數據存儲路徑
	Path string

	// redis 寫入地址，如果設置則使用 redis 爲會話提供緩存
	Write string
	// redis 讀取地址，可以使用 redis 讀寫分離
	// 如果爲空字符串，則使用 Write 設定的值
	Read string
	// 訪問 token 有效時間
	Access time.Duration
	// 刷新 token 有效時間
	Refresh time.Duration
	// session 最長維持多久
	Deadline time.Duration
	// 支持的平臺，每個平臺同一時刻只允許相同用戶存在一個會話，但不同平臺可以存在同一用戶的多個會話
	Platform []string

	// session 簽名算法 HMD5 HS1 HS256 HS384 HS512
	Alg string
	// 簽名密鑰
	Key string
}
