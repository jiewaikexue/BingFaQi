package GrayLogWithZapByTCP

import (
	"2024_4_15/Dingyi"
	"encoding/json"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net"
	"time"
)

type GelfLogger struct {
	conn net.Conn
}

func NewGelfLogger(address string) (*GelfLogger, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	return &GelfLogger{conn: conn}, nil
}

// 这个ok
func (g *GelfLogger) Write(p []byte) (n int, err error) {
	// 解析日志并生成 GELF 消息
	var m map[string]interface{}
	if err = json.Unmarshal(p, &m); err != nil {
		return 0, err
	}

	// 确保包含必要的 GELF 字段
	m["version"] = "1.1"
	m["host"] = Dingyi.GrayLogAddress // 使用 logger 字段作为 host，需要根据实际情况调整
	m["short_message"] = m["message"].(string)
	if _, ok := m["timestamp"]; !ok {
		m["timestamp"] = float64(time.Now().UnixNano()) / 1e9
	}

	// 序列化为 JSON
	msg, err := json.Marshal(m)
	if err != nil {
		return 0, err
	}

	// 发送消息
	msg = append(msg, byte('\x00')) // GELF TCP 消息以 null 字节结尾
	if _, err = g.conn.Write(msg); err != nil {
		return 0, err
	}
	return len(p), nil
}

//func (g *GelfLogger) Write(p []byte) (n int, err error) {
//	var m map[string]interface{}
//	if err = json.Unmarshal(p, &m); err != nil {
//		return 0, fmt.Errorf("JSON unmarshal error: %v", err)
//	}
//
//	m["version"] = "1.1"
//	m["host"] = Dingyi.GrayLogShowName // 替换为实际的主机名或预先定义的变量
//
//	if msg, ok := m["message"].(string); ok {
//		m["short_message"] = msg
//	} else {
//		m["short_message"] = "Default message"
//	}
//
//	m["timestamp"] = float64(time.Now().UnixNano()) / 1e9
//
//	msg, err := json.Marshal(m)
//	if err != nil {
//		return 0, fmt.Errorf("JSON marshal error: %v", err)
//	}
//
//	msg = append(msg, byte('\x00'))
//	if _, err = g.conn.Write(msg); err != nil {
//		return 0, fmt.Errorf("TCP write error: %v", err)
//	}
//	return len(p), nil
//}

//这个也ok
//func (g *GelfLogger) Write(p []byte) (n int, err error) {
//	var jsonData map[string]interface{}
//	if err := json.Unmarshal(p, &jsonData); err != nil {
//		return 0, fmt.Errorf("JSON unmarshal error: %v", err)
//	}
//
//	gelfMessage := map[string]interface{}{
//		"version": "1.1",
//		"host":    Dingyi.GrayLogShowName, // 默认值，如果 jsonData 中没有 host
//	}
//
//	// 动态设置字段，覆盖默认值
//	for key, value := range jsonData {
//		switch key {
//		case "host", "message":
//			gelfMessage[key] = value
//		default:
//			// 将自定义字段加入，以 `_` 开头
//			gelfMessage[key] = value
//		}
//	}
//
//	// 特别处理 message 字段到 short_message
//	if msg, ok := jsonData["message"].(string); ok {
//		gelfMessage["short_message"] = msg
//	}
//
//	// 序列化 GELF 消息
//	msg, err := json.Marshal(gelfMessage)
//	if err != nil {
//		return 0, fmt.Errorf("JSON marshal error: %v", err)
//	}
//
//	// 发送消息
//	msg = append(msg, byte('\x00')) // GELF TCP 消息以 null 字节结尾
//	if _, err = g.conn.Write(msg); err != nil {
//		return 0, fmt.Errorf("TCP write error: %v", err)
//	}
//	return len(p), nil
//}

func (g *GelfLogger) Sync() error {
	return nil
}

func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02-15:04:05"))
}
func GetZapLogger(level string, useJson bool, gelfAddress string) *zap.Logger {
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
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     CustomTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}

	var enc zapcore.Encoder
	if useJson {
		enc = zapcore.NewJSONEncoder(encoderCfg)
	} else {
		enc = zapcore.NewConsoleEncoder(encoderCfg)
	}

	var logLevel zapcore.Level
	switch level {
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
