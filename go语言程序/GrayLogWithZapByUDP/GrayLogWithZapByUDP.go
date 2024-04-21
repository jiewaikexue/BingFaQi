package GrayLogWithZapByUDP

import (
	"2024_4_15/Dingyi"
	"bytes"
	gelf "github.com/Graylog2/go-gelf/gelf"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02-15:04:05"))
}

type GelfLogger struct {
	writer *gelf.Writer
}

func NewGelfLogger(address string) (*GelfLogger, error) {
	writer, err := gelf.NewWriter(address)
	if err != nil {
		return nil, err
	}
	return &GelfLogger{writer: writer}, nil
}

func (g *GelfLogger) Write(p []byte) (n int, err error) {
	// 构造一个基本的 GELF 消息
	m := gelf.Message{
		Version:  "1.1",
		Host:     Dingyi.GrayLogShowName, // 应替换为实际主机名
		Short:    string(p),              // 使用整个日志消息作为短消息
		Full:     string(p),              // 同样，使用整个日志消息作为完整消息
		TimeUnix: float64(time.Now().UnixNano()) / 1e9,
	}

	// 从日志数据中尝试提取日志级别（如果需要可以进行更复杂的解析）
	levelPos := bytes.Index(p, []byte(`"level":"`))
	if levelPos != -1 {
		levelStart := levelPos + 9 // 跳过 `"level":"` 部分
		levelEnd := bytes.IndexByte(p[levelStart:], '"') + levelStart
		if levelEnd > levelStart {
			levelStr := p[levelStart:levelEnd]
			level, lerr := zapcore.ParseLevel(string(levelStr))
			if lerr == nil {
				m.Level = int32(level)
			}
		}
	}

	// 发送消息
	if err = g.writer.WriteMessage(&m); err != nil {
		return 0, err
	}
	return len(p), nil
}

func (g *GelfLogger) Sync() error {
	return nil
}

func GetZapLogger(Level string, UseJson bool, gelfAddress string) *zap.Logger {
	gelfLogger, err := NewGelfLogger(gelfAddress)
	if err != nil {
		panic(err)
	}

	encoderCfg := zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "atomicLevel",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     CustomTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}

	var enc zapcore.Encoder
	if UseJson {
		enc = zapcore.NewJSONEncoder(encoderCfg)
	} else {
		enc = zapcore.NewConsoleEncoder(encoderCfg)
	}

	var logLevel zapcore.Level
	switch Level {
	case "debug":
		logLevel = zapcore.DebugLevel
	case "info":
		logLevel = zapcore.InfoLevel
	default:
		logLevel = zapcore.InfoLevel // Default to Info if level not specified
	}
	zapCore := zapcore.NewCore(enc, zapcore.AddSync(gelfLogger), logLevel)
	return zap.New(zapCore)
}
