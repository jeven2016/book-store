package yzs8

import (
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logLevelMap = map[string]zapcore.Level{
	"DEBUG":  zapcore.DebugLevel,
	"INFO":   zapcore.InfoLevel,
	"WARN":   zapcore.WarnLevel,
	"ERROR":  zapcore.ErrorLevel,
	"DPANIC": zapcore.DPanicLevel,
	"PANIC":  zapcore.PanicLevel,
	"FATAL":  zapcore.FatalLevel,
}

func SetupLog() *zap.Logger {
	level, ok := logLevelMap["INFO"]
	if !ok {
		panic("the log_level is invalid, it only supports: DEBUG,INFO,WARN, ERROR,DPANIC,PANIC,FATAL.  ")
	}

	atomicLevel := zap.NewAtomicLevelAt(level)

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "line",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 短路径编码器
		// EncodeCaller:   zapcore.FullCallerEncoder, // 全路径编码器
		EncodeName: zapcore.FullNameEncoder,
	}

	// 日志轮转
	writer := &lumberjack.Logger{
		// 日志名称
		Filename: "log.log",
		// 日志大小限制，单位MB
		MaxSize: 10,
		// 历史日志文件保留天数
		MaxAge: 30,
		// 最大保留历史日志数量
		MaxBackups: 10,
		// 本地时区
		LocalTime: true,
		// 历史日志文件压缩
		Compress: true,
	}

	syncers := []zapcore.WriteSyncer{zapcore.AddSync(writer)}
	if true {
		// 同时输出到控制台
		syncers = append(syncers, zapcore.AddSync(os.Stdout))
	}

	zapCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(syncers...),
		atomicLevel,
	)

	// 开启开发模式，堆栈跟踪
	options := []zap.Option{zap.AddCaller()}

	// 开启文件及行号
	options = append(options, zap.Development())

	// 设置默认携带的字段
	// options = append(options, zap.Fields(zap.String("serviceName", cfg.General.ServiceName)))

	//替换zap的全局logger对象，之后可以已zap的API调用
	//zap.S().Info("Data encoded => ", data)
	//zap.L().Info()
	logger := zap.New(zapCore, options...)
	zap.ReplaceGlobals(logger)

	return logger
}
