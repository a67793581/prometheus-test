package logger

import (
	"context"
	"log"
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"prometheus-test/lib/util"
)

var alogger *zap.SugaredLogger
var blogger *zap.SugaredLogger

type LoggerConf struct {
	Level         string `toml:"level"`
	Business      string `toml:"business"`
	Access        string `toml:"access"`
	BusinessLink  string `toml:"business_link"`
	Size          int    `toml:"size"`
	AccessLink    string `toml:"access_link"`
	RotationCount uint   `toml:"rotation_count"`
	//
	Console bool `toml:"console"`
}

func Init(cfg LoggerConf) error {
	level := zap.InfoLevel
	switch cfg.Level {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "error":
		level = zap.ErrorLevel
	case "warn":
		level = zap.WarnLevel
	case "fatal":
		level = zap.FatalLevel
	default:
	}

	l, err := initLogger("stdout", "stdout", level, cfg.Size, cfg.RotationCount, zap.AddCaller(), zap.AddCallerSkip(2))
	if err != nil {
		return err
	}
	blogger = l.Sugar()
	alogger = l.Sugar()
	return nil
}

func initLogger(logFile string, logFileLink string, level zapcore.Level, size int, rotationCount uint, options ...zap.Option) (*zap.Logger, error) {
	var w zapcore.WriteSyncer
	w = zapcore.AddSync(os.Stdout)
	if logFile != "stdout" {
		isK8s := os.Getenv("K8S_ENV")
		if isK8s == "" {
			logFile = "./logs/business-%Y%m%d.log-%H%M"
			logFileLink = "./logs/business.log"
		}
		rotator, err := rotatelogs.New(
			logFile,
			rotatelogs.WithLinkName(logFileLink),
			rotatelogs.WithRotationSize(int64(size)*1024*1024*1024),
			rotatelogs.WithRotationTime(time.Hour),
			rotatelogs.WithRotationCount(rotationCount),
		)
		if err != nil {
			return nil, err
		}
		// add the encoder config and rotator to create a new zap logger
		w = zapcore.AddSync(rotator)
	}

	Encoder := GetEncoder()

	core := zapcore.NewCore(
		// zapcore.NewJSONEncoder(enc),
		Encoder,
		w,
		level,
	)
	return zap.New(core, options...), nil
}

func GetEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(
		zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,      // "\n"
			EncodeLevel:    cEncodeLevel,                   // ，:InfoLevel "info"
			EncodeTime:     cEncodeTime,                    //
			EncodeDuration: zapcore.SecondsDurationEncoder, // ，Duration
			// EncodeCaller:   zapcore.ShortCallerEncoder,     //
		})
}
func cEncodeTime(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(logTmFmtWithMS) + "|")
}

func cEncodeLevel(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + level.CapitalString() + "]|")
}

func initAccessLogger(logFile string, logFileLink string, level zapcore.Level, size int, rotationCount uint, options ...zap.Option) (*zap.Logger, error) {
	var w zapcore.WriteSyncer
	w = zapcore.AddSync(os.Stdout)
	if logFile != "stdout" {
		isK8s := os.Getenv("K8S_ENV")
		if isK8s == "" {
			logFile = "./logs/access-%Y%m%d.log-%H%M"
			logFileLink = "./logs/access.log"
		}
		rotator, err := rotatelogs.New(
			logFile,
			rotatelogs.WithLinkName(logFileLink),
			rotatelogs.WithRotationSize(int64(size)*1024*1024*1024),
			rotatelogs.WithRotationTime(time.Hour),
			rotatelogs.WithRotationCount(rotationCount),
		)
		if err != nil {
			return nil, err
		}
		// add the encoder config and rotator to create a new zap logger
		w = zapcore.AddSync(rotator)
	}

	Encoder := GetAccessEncoder()

	core := zapcore.NewCore(
		Encoder,
		w,
		level,
	)
	return zap.New(core, options...), nil
}

func GetAccessEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(
		zapcore.EncoderConfig{
			MessageKey: "msg",
		})
}

const (
	logTmFmtWithMS = "2006-01-02 15:04:05"
)

func AccessInfo(c context.Context, msg string, args ...interface{}) {
	alogger.Infow(msg, args...)
}

func Debug(c context.Context, msg string, args ...interface{}) {
	blogger.Debugw(util.GetRequestId(c)+"|"+msg, args...)
}

func Debugf(c context.Context, template string, args ...interface{}) {
	blogger.Debugf(util.GetRequestId(c)+"|"+template, args...)
}

func Warn(c context.Context, msg string, args ...interface{}) {
	blogger.Warnw(util.GetRequestId(c)+"|"+msg, args...)
}

func Warnf(c context.Context, msg string, args ...interface{}) {
	blogger.Warnf(util.GetRequestId(c)+"|"+msg, args...)
}

func Info(c context.Context, msg string, args ...interface{}) {
	blogger.Infow(util.GetRequestId(c)+"|"+msg, args...)
}

func Infof(c context.Context, msg string, args ...interface{}) {
	blogger.Infof(util.GetRequestId(c)+"|"+msg, args...)
}

func Error(c context.Context, msg string, args ...interface{}) {
	blogger.Errorw(util.GetRequestId(c)+"|"+msg, args...)
}

func Errorf(c context.Context, template string, args ...interface{}) {
	blogger.Errorf(util.GetRequestId(c)+"|"+template, args...)
}

func Fatal(c context.Context, msg string, args ...interface{}) {
	blogger.Fatalw(util.GetRequestId(c)+"|"+msg, args...)
}

func NotCtxInfo(msg string, args ...interface{}) {
	blogger.Infow("notCtx|"+msg, args...)
}
func NotCtxInfof(msg string, args ...interface{}) {
	blogger.Infof("notCtx|"+msg, args...)
}
func NotCtxFatal(msg string, args ...interface{}) {
	blogger.Fatalw("notCtx|"+msg, args...)
}
func NotCtxFatalf(msg string, args ...interface{}) {
	blogger.Fatalf("notCtx|"+msg, args...)
}
func NotCtxError(msg string, args ...interface{}) {
	blogger.Errorw("notCtx|"+msg, args...)
}
func NotCtxErrorf(msg string, args ...interface{}) {
	blogger.Errorf("notCtx|"+msg, args...)
}

func Close() {
	err := blogger.Sync()
	if err != nil {
		log.Printf("[DS]BizLogger Close Sync failed,err=%v", err)
	}
	err = alogger.Sync()
	if err != nil {
		log.Printf("[DS]AccessLogger Close Sync failed,err=%v", err)
	}
}

func GetBizLogger() *zap.SugaredLogger {
	return blogger
}
