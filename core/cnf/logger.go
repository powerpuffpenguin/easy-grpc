package cnf

// 日誌定義
type Logger struct {
	// 日誌存儲路徑
	Filename string
	// 單個檔案最大尺寸 MB
	MaxSize int
	// 最多保存多少個檔案
	MaxBackups int
	// 最多保存多少天日誌
	MaxDays int
	// 日誌寫入緩存大小
	BufferSize int
	// 如果爲 true 在日誌中輸出代碼所在行
	Caller bool
	// 輸出到檔案的日誌等級 [debug info warn error dpanic panic fatal]
	FileLevel string
	// 輸出到控制檯的日誌等級 [debug info warn error dpanic panic fatal]
	ConsoleLevel string
}
