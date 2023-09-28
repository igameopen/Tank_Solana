package log

import (
	"context"
	"io"
	"strings"
	"time"

	"transferSrv/infra/config"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in
	// production.
	DebugLevel zapcore.Level = iota - 1
	// InfoLevel is the default logging priority.
	InfoLevel
	// WarnLevel logs are more important than Info, but don't need individual
	// human review.
	WarnLevel
	// ErrorLevel logs are high-priority. If an application is running smoothly,
	// it shouldn't generate any error-level logs.
	ErrorLevel
	// DPanicLevel logs are particularly important errors. In development the
	// logger panics after writing the message.
	DPanicLevel
	// PanicLevel logs a message, then panics.
	PanicLevel
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel
)

var logger *zap.Logger
var path, level, encoding string

func transLeve(lv string) zapcore.Level {
	switch strings.ToLower(lv) {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn":
		return WarnLevel
	case "error":
		return ErrorLevel
	case "dpanic":
		return DPanicLevel
	case "panic":
		return PanicLevel
	case "fatal":
		return FatalLevel
	default:
		return InfoLevel
	}
}

func setOption() {
	path = config.GetCnf().LogCnf.Output

	level = config.GetCnf().LogCnf.Level
	encoding = "json"

	// path = "../../log/"
	// level = "debug"
}

func Init() {

	setOption()

	encoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		NameKey:     "name",
		EncodeLevel: zapcore.CapitalLevelEncoder,
		TimeKey:     "ts",

		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		CallerKey:    "file",
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	})

	//infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
	//	return lvl <= zapcore.InfoLevel
	//})

	warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.WarnLevel
	})

	// infoWriter := getWriter(path + "info")
	warnWriter := getWriter(path + "warn")
	allWriter := getWriter(path + "log")

	atom := zap.NewAtomicLevel()
	atom.SetLevel(transLeve(level))

	// 最后创建具体的Logger
	core := zapcore.NewTee(
		// zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel), //按日志等级分开写，不受配置等级限制
		zapcore.NewCore(encoder, zapcore.AddSync(warnWriter), warnLevel), //按日志等级分开写，不受配置等级限制
		zapcore.NewCore(encoder, zapcore.AddSync(allWriter), atom),       //按日志等级写
	)

	// before := func(e zapcore.Entry) error {
	// 	fmt.Println("before ", e.Message)
	// 	e.Message = "before " + e.Message
	// 	return nil
	// }

	// logger = zap.New(core, zap.Hooks(before), zap.AddCaller())
	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
}

func getWriter(filename string) io.Writer {
	hook, err := rotatelogs.New(
		filename+"-%Y%m%d"+".log",
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Hour*24*7),
		rotatelogs.WithRotationTime(time.Hour*24),
	)

	if err != nil {
		panic(err)
	}
	return hook
}

func GetLogger() *zap.SugaredLogger {
	return logger.Sugar()
}

// 可通过context将trace-id渗透到日志
func GetLoggerWithCtx(ctx context.Context) *zap.SugaredLogger {
	return logger.WithOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		trace := ctx.Value("trace-id")
		if trace == nil {
			trace = 0
		}
		if tmp, ok := trace.(string); ok {
			return c.With([]zap.Field{zap.String("trace-id", tmp)})
		}
		return c
	})).Sugar()
}

func WithContext(ctx context.Context) *zap.SugaredLogger {
	return GetLoggerWithCtx(ctx)
}

// Debug uses fmt.Sprint to construct and log a message.
func Debug(args ...interface{}) {
	GetLogger().Debug(args...)
}

// Info uses fmt.Sprint to construct and log a message.
func Info(args ...interface{}) {
	GetLogger().Info(args...)
}

// Warn uses fmt.Sprint to construct and log a message.
func Warn(args ...interface{}) {
	GetLogger().Warn(args...)
}

// Error uses fmt.Sprint to construct and log a message.
func Error(args ...interface{}) {
	GetLogger().Error(args...)
}

// DPanic uses fmt.Sprint to construct and log a message. In development, the
// logger then panics. (See DPanicLevel for details.)
func DPanic(args ...interface{}) {
	GetLogger().DPanic(args...)
}

// Panic uses fmt.Sprint to construct and log a message, then panics.
func Panic(args ...interface{}) {
	GetLogger().Panic(args...)
}

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
func Fatal(args ...interface{}) {
	GetLogger().Fatal(args...)
}

// Debugf uses fmt.Sprintf to log a templated message.
func Debugf(template string, args ...interface{}) {
	GetLogger().Debugf(template, args...)
}

// Infof uses fmt.Sprintf to log a templated message.
func Infof(template string, args ...interface{}) {
	GetLogger().Infof(template, args...)
}

// Warnf uses fmt.Sprintf to log a templated message.
func Warnf(template string, args ...interface{}) {
	GetLogger().Warnf(template, args...)
}

// Errorf uses fmt.Sprintf to log a templated message.
func Errorf(template string, args ...interface{}) {
	GetLogger().Errorf(template, args...)
}

// DPanicf uses fmt.Sprintf to log a templated message. In development, the
// logger then panics. (See DPanicLevel for details.)
func DPanicf(template string, args ...interface{}) {
	GetLogger().DPanicf(template, args...)
}

// Panicf uses fmt.Sprintf to log a templated message, then panics.
func Panicf(template string, args ...interface{}) {
	GetLogger().Panicf(template, args...)
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func Fatalf(template string, args ...interface{}) {
	GetLogger().Fatalf(template, args...)
}

func Sync() {
	GetLogger().Sync()
}
