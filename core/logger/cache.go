package logger

import (
	"runtime"

	"go.uber.org/zap/zapcore"
)

// 這是一個日誌緩存，使用它來提高日誌寫入效率。
// 它會儘量合併來自不同 goroutine 的日誌一併寫入，但是它不會保證日誌寫入成功也不會回報寫入錯誤
// 回報寫入錯誤是可行的，但它會相當影響效率，故本喵認爲沒有不要)
type loggerCache struct {
	ch     chan []byte
	writer zapcore.WriteSyncer
}

func newLoggerCache(writer zapcore.WriteSyncer, bufferSize int) *loggerCache {
	l := &loggerCache{
		ch:     make(chan []byte, runtime.GOMAXPROCS(0)),
		writer: writer,
	}
	go l.run(bufferSize)
	return l
}
func (l *loggerCache) Sync() error {
	return l.writer.Sync()
}
func (l *loggerCache) Write(b []byte) (count int, e error) {
	count = len(b)
	if count > 0 {
		data := make([]byte, count)
		copy(data, b)
		l.ch <- data
	}
	return
}
func (l *loggerCache) run(bufferSize int) {
	var (
		r = buffer{
			ch: l.ch,
			b:  make([]byte, bufferSize),
		}
		b0, b1 []byte
	)
	for {
		b0, b1 = r.Read()
		if len(b0) > 0 {
			l.writer.Write(b0)
		}
		if len(b1) > 0 {
			l.writer.Write(b1)
		}
	}
}

type buffer struct {
	ch <-chan []byte
	b  []byte
}

func (b *buffer) Read() (b0, b1 []byte) {
	var (
		count = len(b.b)
		n, rn int
		data  []byte
	)
	for n == 0 {
		data = <-b.ch
		rn = len(data)
		if rn < 1 {
			continue
		}
		if rn >= count {
			b0 = data
			return
		}
		b0 = b.b
		n = copy(b0, data)
	}
	for {
		select {
		case data = <-b.ch:
		default:
			b0 = b0[:n]
			return
		}
		rn = copy(b0[n:], data)
		n += rn
		data = data[rn:]
		if len(data) != 0 {
			b1 = data
			return
		}
	}
	return
}
