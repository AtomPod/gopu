package log

import (
	"context"
	"io"
	"os"

	"github.com/ngs24313/gopu/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	//DefaultLogger default logger
	DefaultLogger  *zap.Logger = zap.NewExample()
	customerWriter []io.Writer
)

func init() {
	logger, err := zap.NewDevelopmentConfig().Build()
	if err != nil {
		panic(err)
	}
	DefaultLogger = logger
}

//AddWriter add writer to logger
func AddWriter(w io.Writer) {
	customerWriter = append(customerWriter, w)
}

//Init initiaize default logger
func Init(conf *config.Config, redirectStd bool) {
	logConf := conf.Logger

	mode := logConf.Mode
	if mode == "" {
		mode = conf.Mode
	}

	var logger *zap.Logger
	if mode == "release" {

		rotateLogger := &lumberjack.Logger{
			Filename:   logConf.Filename,
			MaxAge:     logConf.MaxAge,
			MaxBackups: logConf.MaxBackups,
			MaxSize:    logConf.MaxSize,
		}

		syncers := []zapcore.WriteSyncer{
			zapcore.AddSync(rotateLogger),
			zapcore.AddSync(os.Stdout),
		}

		if len(customerWriter) > 0 {
			for _, w := range customerWriter {
				syncers = append(syncers, zapcore.AddSync(w))
			}
		}

		syncer := zapcore.NewMultiWriteSyncer(syncers...)
		encoder := zap.NewProductionEncoderConfig()
		core := zapcore.NewCore(zapcore.NewJSONEncoder(encoder), syncer, zapcore.InfoLevel)
		logger = zap.New(core)

	} else {

		config := zap.NewDevelopmentEncoderConfig()
		config.EncodeLevel = zapcore.CapitalColorLevelEncoder

		var syncer zapcore.WriteSyncer
		syncer = zapcore.AddSync(os.Stdout)

		if len(customerWriter) > 0 {
			syncers := []zapcore.WriteSyncer{
				zapcore.AddSync(os.Stdout),
			}
			for _, w := range customerWriter {
				syncers = append(syncers, zapcore.AddSync(w))
			}
			syncer = zapcore.NewMultiWriteSyncer(syncers...)
		}

		core := zapcore.NewCore(
			zapcore.NewConsoleEncoder(config),
			syncer,
			zapcore.DebugLevel)

		logger = zap.New(core, zap.AddStacktrace(zap.ErrorLevel))
	}

	if redirectStd {
		zap.RedirectStdLogAt(logger, zap.DebugLevel)
	}

	DefaultLogger = logger
}

//Logger get default logger
func Logger(ctx context.Context) *zap.Logger {
	return DefaultLogger
}

//Debug print debug log
func Debug(msg string, fields ...zap.Field) {
	DefaultLogger.Debug(msg, fields...)
}

//Info print info log
func Info(msg string, fields ...zap.Field) {
	DefaultLogger.Info(msg, fields...)
}

//Warn print warn log
func Warn(msg string, fields ...zap.Field) {
	DefaultLogger.Warn(msg, fields...)
}

//Error print error log
func Error(msg string, fields ...zap.Field) {
	DefaultLogger.Error(msg, fields...)
}

//Panic print panic log
func Panic(msg string, fields ...zap.Field) {
	DefaultLogger.Panic(msg, fields...)
}

//Fatal print panic log
func Fatal(msg string, fields ...zap.Field) {
	DefaultLogger.Fatal(msg, fields...)
}
