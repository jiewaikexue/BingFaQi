package log

import (
	"fmt"
	"github.com/shengyanli1982/law"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

/*
传统异步文件写入,
效果非常满意 但是我还是选择接入grayLog
2024年4月19日22:34:32
*/

func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02-15:04:05"))
}

func GetLawLogger(ObjectName string, Filename string, Mod string) (*law.WriteAsyncer, *os.File) {
	Filename = "log/" + ObjectName + "_" + Filename + ".log"
	// Create a new WriteAsyncer instance using os.Stdout
	file, err := os.OpenFile(Filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Println("日志文件打开失败:", err)
	}
	if Mod == "console" {
		file = os.Stdout
	}
	aw := law.NewWriteAsyncer(file, nil)
	return aw, file
}

func GetZapLogger(A *law.WriteAsyncer, Level string, UseJson bool) *zap.Logger {
	aw := A
	encoderCfg := zapcore.EncoderConfig{
		//MessageKey:     "msg",                         // 消息的键名
		//LevelKey:       "level",                       // 级别的键名
		//NameKey:        "logger",                      // 记录器名的键名
		//EncodeLevel:    zapcore.LowercaseLevelEncoder, // 级别的编码器
		//EncodeTime:     zapcore.ISO8601TimeEncoder,    // 时间的编码器
		//EncodeDuration: zapcore.StringDurationEncoder, // 持续时间的编码器
		//当存储的格式为JSON的时候这些作为可以key
		MessageKey:    "message",
		LevelKey:      "atomicLevel",
		TimeKey:       "time",
		NameKey:       "logger",
		CallerKey:     "caller",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		//以上字段输出的格式
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     CustomTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	// 使用 WriteAsyncer 创建一个 zapcore.WriteSyncer 实例
	// Create a zapcore.WriteSyncer instance using WriteAsyncer
	zapAsyncWriter := zapcore.AddSync(aw)

	// 使用编码器配置和 WriteSyncer 创建一个 zapcore.Core 实例
	// Create a zapcore.Core instance using the encoder configuration and WriteSyncer
	var enc zapcore.Encoder
	if UseJson == true {
		enc = zapcore.NewJSONEncoder(encoderCfg)
	} else {
		enc = zapcore.NewConsoleEncoder(encoderCfg)
	}
	if Level == "debug" {
		zapCore := zapcore.NewCore(enc, zapAsyncWriter, zapcore.DebugLevel)
		// 使用 Core 创建一个 zap.Logger 实例
		// Create a zap.Logger instance using Cores
		return zap.New(zapCore)
	} else if Level == "info" {
		zapCore := zapcore.NewCore(enc, zapAsyncWriter, zapcore.InfoLevel)
		return zap.New(zapCore)
	} else {
		zapCore := zapcore.NewCore(enc, zapAsyncWriter, zapcore.InfoLevel)
		return zap.New(zapCore)
	}
}
