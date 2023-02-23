package logger

import (
	"github.com/powerpuffpenguin/easy-grpc/core/cnf"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}

var emptyAtomicLevel = zap.NewAtomicLevel()

// 對 zap.Logger 的包裝
type Wrap struct {
	noCopy noCopy
	// zap 日誌
	*zap.Logger

	// 檔案日誌
	fileLevel zap.AtomicLevel
	// 控制檯日誌
	consoleLevel zap.AtomicLevel
}

func New(options *cnf.Logger, zapOptions ...zap.Option) *Wrap {
	bufferSize := options.BufferSize
	if bufferSize < 128 {
		bufferSize = 1024 * 32
	}

	var cores []zapcore.Core
	fileLevel := zap.NewAtomicLevel()
	consoleLevel := zap.NewAtomicLevel()
	if options.FileLevel == "" {
		fileLevel.SetLevel(zap.FatalLevel)
	} else if e := fileLevel.UnmarshalText([]byte(options.FileLevel)); e != nil {
		fileLevel.SetLevel(zap.FatalLevel)
	}
	cores = append(cores, zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		newLoggerCache(zapcore.AddSync(&lumberjack.Logger{
			Filename:   options.Filename,
			MaxSize:    options.MaxSize, // megabytes
			MaxBackups: options.MaxBackups,
			MaxAge:     options.MaxDays, // days
		}), bufferSize),
		fileLevel,
	))

	// console
	if options.ConsoleLevel == "" {
		consoleLevel.SetLevel(zap.FatalLevel)
	} else if e := consoleLevel.UnmarshalText([]byte(options.ConsoleLevel)); e != nil {
		consoleLevel.SetLevel(zap.FatalLevel)
	}
	cores = append(cores, zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		newLoggerCache(monitor, bufferSize),
		consoleLevel,
	))

	if options.Caller {
		zapOptions = append(zapOptions, zap.AddCaller())
	}

	return &Wrap{
		Logger: zap.New(
			zapcore.NewTee(cores...),
			zapOptions...,
		),
		fileLevel:    fileLevel,
		consoleLevel: consoleLevel,
	}
}
func (l *Wrap) Attach(src *Wrap) {
	l.Logger = src.Logger
	l.fileLevel = src.fileLevel
	l.consoleLevel = src.consoleLevel
}

// Detach Logger
func (l *Wrap) Detach() {
	l.Logger = nil
	l.fileLevel = emptyAtomicLevel
	l.consoleLevel = emptyAtomicLevel
}

func (l *Wrap) FileLevel() zap.AtomicLevel {
	return l.fileLevel
}

func (l *Wrap) ConsoleLevel() zap.AtomicLevel {
	return l.consoleLevel
}
