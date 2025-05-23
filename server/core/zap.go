package core

import (
	"log"
	"os"
	"server/global"
	
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitLogger 初始化并返回一个基于配置设置的新 zap.Logger 实例
func InitLogger() *zap.Logger {
	zapCfg := global.Config.Zap
	
	// 创建一个用于日志输出的 writeSyncer
	writeSyncer := getLogWriter(zapCfg.Filename, zapCfg.MaxSize, zapCfg.MaxBackups, zapCfg.MaxAge)
	
	// 如果配置了控制台输出，则添加控制台输出
	if zapCfg.IsConsolePrint {
		writeSyncer = zapcore.NewMultiWriteSyncer(writeSyncer, zapcore.AddSync(os.Stdout))
	}
	
	// 创建日志格式化的编码器
	encoder := getEncoder()
	
	// 根据配置确定日志级别
	var logLevel zapcore.Level
	
	if err := logLevel.UnmarshalText([]byte(zapCfg.Level)); err != nil {
		log.Fatalf("Failed to parse log level: %v", err)
	}
	
	// 创建核心和日志实例
	core := zapcore.NewCore(encoder, writeSyncer, logLevel)
	logger := zap.New(core, zap.AddCaller())
	return logger
}

// getLogWriter 返回一个 zapcore.WriteSyncer，该写入器利用 lumberjack 包，实现日志的滚动记录 ：轮转（Log Rotation）是一种日志管理技术，当日志文件达到指定大小或时间限制时，系统会自动将当前日志文件归档，并创建一个新的日志文件继续记录。github.com/natefinch/lumberjack包实现了这一功能，支持按文件大小（MaxSize）、保留文件数量（MaxBackups）和保留天数（MaxAge）进行日志轮转。
func getLogWriter(filename string, maxSize, maxBackups, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,   // 日志文件的位置
		MaxSize:    maxSize,    // 在进行切割之前，日志文件的最大大小（以MB为单位）
		MaxBackups: maxBackups, // 保留旧文件的最大个数
		MaxAge:     maxAge,     // 保留旧文件的最大天数
	}
	return zapcore.AddSync(lumberJackLogger) // AddSync(io *io.Writer) WriteSyncer
}

// getEncoder 返回一个为生产日志配置的 JSON 编码器
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}
