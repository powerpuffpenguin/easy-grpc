package logger

import (
	"os"
	"sync"

	"go.uber.org/zap/zapcore"
)

var monitor = &writerMonitor{
	writer:    os.Stdout,
	listeners: make(map[Listener]bool),
}

// 增加一個監聽器用於獲取 輸出到 stdout 的日誌
func AddListener(l Listener) {
	monitor.Lock()
	monitor.listeners[l] = true
	monitor.Unlock()
}

// 移除監聽器
func RemoveListener(l Listener) {
	monitor.Lock()
	delete(monitor.listeners, l)
	monitor.Unlock()
}

// 這是一個監聽器，當有日誌輸出到 stdout 時它會被回調
type Listener interface {
	OnChanged([]byte)
}
type writerMonitor struct {
	writer    zapcore.WriteSyncer
	listeners map[Listener]bool
	sync.Mutex
}

func (w *writerMonitor) Write(data []byte) (count int, e error) {
	count, e = w.writer.Write(data)
	if count > 0 {
		w.Lock()
		for listener := range w.listeners {
			listener.OnChanged(data[:count])
		}
		w.Unlock()
	}
	return
}
func (w *writerMonitor) Sync() error {
	return w.writer.Sync()
}

// 這是一個 stdout 的日誌監聽器，如果你不即使讀取它，它將只保留最後的兩個日誌
type SnapshotListener struct {
	done <-chan struct{}
	ch   chan []byte
}

// 創建監聽器
func NewSnapshotListener(done <-chan struct{}) *SnapshotListener {
	return &SnapshotListener{
		done: done,
		ch:   make(chan []byte, 2),
	}
}

// 讀取收到的日誌
func (l *SnapshotListener) Channel() <-chan []byte {
	return l.ch
}

// 首先 Listener 接口
func (l *SnapshotListener) OnChanged(b []byte) {
	data := make([]byte, len(b))
	copy(data, b)

	select {
	case <-l.done:
		return
	case l.ch <- data:
		return
	default:
	}

	// 緩存已滿嘗試丟棄數據
	for {
		select {
		case <-l.done:
			return
		case <-l.ch: // 丟棄最早的數據
		case l.ch <- data: // 寫入數據
		}
	}
}
