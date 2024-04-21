// 放弃使用

package Suglog

import (
	"github.com/natefinch/lumberjack"
	"github.com/shengyanli1982/law"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"time"
	//	go get github.com/shengyanli1982/law
)

//var sugarLogger *zap.SugaredLogger

//func InitLogger() *zap.SugaredLogger {
//	encoder := getEncoder()
//	writeSyncer := getLogWriter()
//	core := zapcore.NewCore(encoder, writeSyncer, zapcore.InfoLevel)
//	//level := zap.NewAtomicLevel()
//	//zap.AddCaller()  添加将调用函数信息记录到日志中的功能。
//	logger := zap.New(core, zap.AddCaller())
//	sugarLogger := logger.Sugar()
//	if sugarLogger == nil {
//		return nil
//	}
//	return sugarLogger
//}
//
//func InitLogger_Info(ProjectName string, AW_TXT *law.WriteAsyncer, AW_CTL *law.WriteAsyncer) *zap.Logger {
//	encoder := getEncoder()
//	//writeSyncer_txt := getLogWriter_TXT(ProjectName, AW_TXT)
//	writeSyncer_ctl := getLogWriter_CTL(AW_CTL)
//	// 设置日志级别为 INFO
//	//atomicLevel_txt := zap.NewAtomicLevelAt(zapcore.InfoLevel)
//	atomicLevel_ctl := zap.NewAtomicLevelAt(zapcore.DebugLevel)
//	//双端同时记录
//	core := zapcore.NewTee(
//		//zapcore.NewCore(encoder, writeSyncer_txt, atomicLevel_txt),
//		zapcore.NewCore(encoder, writeSyncer_ctl, atomicLevel_ctl), // 将日志同时输出到控制台
//	)
//
//	// 添加 Caller 选项，记录调用函数信息到日志中
//	logger := zap.New(core, zap.AddCaller())
//	// 创建 SugaredLogger
//	sugarLogger := logger
//	if sugarLogger == nil {
//		return nil
//	}
//	return sugarLogger
//}

//func getEncoder() zapcore.Encoder {
//	encoderConfig := zap.NewProductionEncoderConfig()
//	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder        // 修改时间编码器
//	encoderConfig.EncodeDuration = zapcore.StringDurationEncoder // 持续时间的编码器
//	// 在日志文件中使用大写字母记录日志级别
//
//	// ==============================================
//	//	NewConsoleEncoder 打印更符合人们观察的方式
//	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
//	return zapcore.NewConsoleEncoder(encoderConfig)
//	// ==============================================
//	//			Json文件格式返回
//	//return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
//	// ==============================================
//}

// ===============================
//	没有插件lumberjack
// ===============================

func getLogWriter_TXT(ProjectName string, aw_txt *law.WriteAsyncer) zapcore.WriteSyncer {
	return zapcore.AddSync(aw_txt)
}

// ===============================
//
//	引入插件lumberjack
//
// ===============================
//
//	func getLogWriter_TXT(ProjectName string) zapcore.WriteSyncer {
//		Name := "Suglog/" + ProjectName + ".log"
//		lumberJackLogger := &lumberjack.Logger{
//			Filename:   Name,
//			MaxSize:    1000,
//			MaxBackups: 5,
//			MaxAge:     30,
//			Compress:   false,
//		}
//		aw := law.NewWriteAsyncer(lumberJackLogger, nil)
//		return zapcore.AddSync(aw)
//	}
func getLogWriter_CTL(aw_ctl *law.WriteAsyncer) zapcore.WriteSyncer {
	return zapcore.AddSync(aw_ctl)
}

// https://zhuanlan.zhihu.com/p/672671600  完全异步

var L *zap.Logger

const (
	// _defaultBufferSize specifies the default size used by Buffer.
	_defaultBufferSize = 256 * 4096 // 256 kB

	// _defaultFlushInterval specifies the default flush interval for
	// Buffer.
	_defaultFlushInterval = 30 * time.Second
)

func InitLogger(loggerConfig Config) {
	L = loggerConfig.Build()
}

func (lc *Config) parseLevel() zap.AtomicLevel {
	level, err := zap.ParseAtomicLevel(lc.Level)
	if err != nil {
		log.Panicf("init level failed level %s err %v", lc.Level, err)
	}
	return level
}

type Config struct {
	//日志级别 debug info warn panic
	Level string
	//panic时候 是否显示堆栈 panic级别的日志输出堆栈信息。
	Stacktrace bool
	//添加调用者信息
	AddCaller bool
	//调用链，往上多少级 ，在一些中间件，对日志有包装，可以通过这个选项指定。
	CallerShip int
	//输出到哪里标准输出console,还是文件file
	Mode string
	//文件名称加路径
	FileName string
	//error级别的日志输入到不同的地方
	ErrorFileName string
	// 日志文件大小 单位MB 默认500MB
	MaxSize int
	//日志保留天数
	MaxAge int
	//日志最大保留的个数
	MaxBackup int
	//异步日志 日志将先输入到内存到，定时批量落盘。如果设置这个值，要保证在程序退出的时候调用Sync(),在开发阶段不用设置为true。
	Async bool
	//是否 输出json格式的数据，JSON格式相对于console格式，不方便阅读，但是对机器更加友好
	//最佳实践，在开发的时候json为false,mode为console
	Json bool
	//是否日志压缩
	Compress    bool
	options     []zap.Option
	atomicLevel zap.AtomicLevel
}

func (lc *Config) UpdateLevel(level zapcore.Level) {
	lc.atomicLevel.SetLevel(level)
}

func (lc *Config) Build() *zap.Logger {
	if lc.Mode == "file" && lc.FileName == "" {
		log.Printf("file mode, but file name is empty")
	}
	var (
		ws      zapcore.WriteSyncer
		errorWs zapcore.WriteSyncer
		encoder zapcore.Encoder
	)
	encoderConfig := zapcore.EncoderConfig{
		//当存储的格式为JSON的时候这些作为可以key
		MessageKey:    "message",
		LevelKey:      "atomicLevel",
		TimeKey:       "time",
		NameKey:       "logger",
		CallerKey:     "caller",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		//以上字段输出的格式
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     CustomTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	if lc.Mode == "console" {
		ws = zapcore.Lock(os.Stdout)
		errorWs = zapcore.Lock(os.Stderr)
		//输出到控制台彩色。
		if !lc.Json {
			encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		}
	} else {
		normalConfig := &lumberjack.Logger{
			Filename:   lc.FileName,
			MaxSize:    lc.MaxSize,
			MaxAge:     lc.MaxAge,
			MaxBackups: lc.MaxBackup,
			LocalTime:  true,
			Compress:   lc.Compress,
		}
		if lc.ErrorFileName != "" {
			errorConfig := &lumberjack.Logger{
				Filename:   lc.ErrorFileName,
				MaxSize:    lc.MaxSize,
				MaxAge:     lc.MaxAge,
				MaxBackups: lc.MaxBackup,
				LocalTime:  true,
				Compress:   lc.Compress,
			}
			errorWs = zapcore.Lock(zapcore.AddSync(errorConfig))
		}

		ws = zapcore.Lock(zapcore.AddSync(normalConfig))

	}
	if lc.Async {
		ws = &zapcore.BufferedWriteSyncer{
			WS:            ws,
			Size:          _defaultBufferSize,
			FlushInterval: _defaultFlushInterval,
		}
		if errorWs != nil {
			errorWs = &zapcore.BufferedWriteSyncer{
				WS:            errorWs,
				Size:          _defaultBufferSize,
				FlushInterval: _defaultFlushInterval,
			}
		}

	}
	if lc.Json {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}
	var (
		core zapcore.Core
	)
	atomicLevel := lc.parseLevel()
	if lc.ErrorFileName != "" && lc.Mode == "file" {
		lowCore := zapcore.NewCore(encoder, ws, atomicLevel)
		c := []zapcore.Core{lowCore}
		if errorWs != nil {
			highCore := zapcore.NewCore(encoder, errorWs, zapcore.ErrorLevel)
			c = append(c, highCore)
		}
		core = zapcore.NewTee(c...)
	} else {
		core = zapcore.NewCore(encoder, ws, atomicLevel)
	}
	logger := zap.New(core)
	//是否新增调用者信息
	if lc.AddCaller {
		lc.options = append(lc.options, zap.AddCaller())
		if lc.CallerShip != 0 {
			lc.options = append(lc.options, zap.AddCallerSkip(lc.CallerShip))
		}
	}
	//当错误时是否添加堆栈信息
	if lc.Stacktrace {
		lc.options = append(lc.options, zap.AddStacktrace(zap.PanicLevel))
	}

	lc.atomicLevel = atomicLevel
	return logger.WithOptions(lc.options...)

}

func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02-15:04:05"))
}
