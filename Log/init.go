package Log

import (
	"io"
	"os"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func getWriter(filename string) io.Writer {
	// 生成rotatelogs的Logger 实际生成的文件名 demo.log.YYmmddHH
	// demo.log是指向最新日志的链接
	// 保存7天内的日志，每1小时(整点)分割一次日志
	hook, err := rotatelogs.New(
		strings.Replace(filename, ".log", "", -1)+"-%Y%m%d%H.log", // 没有使用go风格反人类的format格式
		//rotatelogs.WithLinkName(filename),
		//rotatelogs.WithMaxAge(time.Hour*24*7),
		rotatelogs.WithRotationTime(time.Hour),
		//rotatelogs.WithRotationCount(100),
	)

	if err != nil {
		panic(err)
	}
	return hook
}

func Init() {

	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel
	})

	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	infoWriter := getWriter("./demo_info.log")
	errorWriter := getWriter("./demo_error.log")
	encoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
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

	//consoleDebugging := zapcore.Lock(os.Stdout)
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.DebugLevel
	})

	// for human operators.

	var allCore []zapcore.Core
	allCore = append(allCore, zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel))
	allCore = append(allCore, zapcore.NewCore(encoder, zapcore.AddSync(errorWriter), errorLevel))
	allCore = append(allCore, zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), lowPriority))
	core := zapcore.NewTee(
		allCore...,
	)

	log := zap.New(core, zap.AddCaller())
	sugar := log.Sugar()
	for i := 0; i < 1000; i++ {
		sugar.Info("111")
		sugar.Error("111")
	}
	defer sugar.Sync()
}
