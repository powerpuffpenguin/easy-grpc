package cnf

import (
	"time"

	"google.golang.org/grpc/keepalive"
)

// http 服務定義
type HTTP struct {
	// 服務器監聽地址
	Addr string

	// tls 證書，如果爲空字符串則使用 http
	CertFile string
	KeyFile  string

	// 爲 api 提供 Swagger 服務 /serve/swagger/
	Swagger bool
	// 爲 靜態檔案提供 服務 /serve/assets/
	Assets bool
	// 爲 項目提供在線 說明 /serve/document/
	Document bool
	// 提供在線的 protobuf 定義 /serve/protobuf/
	Protobuf bool

	// 服務器一些細節定義
	Option ServerOption
}

func (h *HTTP) H2() bool {
	return h.CertFile != `` && h.KeyFile != ``
}

func (h *HTTP) H2C() bool {
	return h.CertFile == `` || h.KeyFile == ``
}

type ServerOption struct {
	WriteBufferSize, ReadBufferSize          int
	InitialWindowSize, InitialConnWindowSize int32
	MaxRecvMsgSize, MaxSendMsgSize           int
	MaxConcurrentStreams                     uint32
	ConnectionTimeout                        time.Duration
	Keepalive                                keepalive.ServerParameters
}
