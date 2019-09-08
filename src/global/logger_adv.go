package global

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func SetAdvLogger() error {
	// 打印什么日志级别
	errorsLogPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})

	fullLogPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.DebugLevel
	})

	// 打印到什么地方
	fullLogFile, err := os.OpenFile("./gin.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	errorLogFile, err := os.OpenFile("./err.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	//In particular, *os.Files must be locked before use.
	logToFullLogFile := zapcore.Lock(fullLogFile)
	LogToErrorLogFile := zapcore.Lock(errorLogFile)
	LogToConsole := zapcore.Lock(os.Stdout)

	// 其他 io.writer 使用 zapcore.AddSync
	//zapcore.AddSync(kafka io.writer)

	// 打印设置,两种打印方式 json / console
	var consoleCfg = zapcore.EncoderConfig{
		MessageKey: "message", // 消息字段名
		LevelKey:   "level",   // 级别字段名
		TimeKey:    "time",    // 时间字段名
		//CallerKey:     "file",    // 记录源码文件的字段名
		//StacktraceKey: "trace",   // 记录trace

		//// Caller的编码器,FullCallerEncoder,ShortCallerEncoder
		//EncodeCaller: zapcore.FullCallerEncoder,

		// 大写彩色编码
		EncodeLevel: zapcore.CapitalColorLevelEncoder,

		// 编码时间字符串的格式
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
	}

	var jsonCfg = zapcore.EncoderConfig{
		MessageKey: "message", // 消息字段名
		LevelKey:   "level",   // 级别字段名
		TimeKey:    "time",    // 时间字段名
		//CallerKey:     "file",    // 记录源码文件的字段名
		//StacktraceKey: "trace",   // 记录trace

		////Caller的编码器,FullCallerEncoder,ShortCallerEncoder
		//EncodeCaller: zapcore.FullCallerEncoder,

		// 大写编码
		EncodeLevel: zapcore.CapitalLevelEncoder,

		// 编码时间字符串的格式
		EncodeTime: func(t time.Time, p zapcore.PrimitiveArrayEncoder) {
			p.AppendString(t.Format("2006-01-02 15:04:05"))
		},
	}

	// 用什么格式打印 —— json / console
	jsonEncoder := zapcore.NewJSONEncoder(jsonCfg)
	consoleEncoder := zapcore.NewConsoleEncoder(consoleCfg)

	// 设置Logger
	core := zapcore.NewTee(
		// zapcore.NewCore 参数 —— 用什么格式打印 | 打印到什么地方 | 打印什么日志级别
		zapcore.NewCore(jsonEncoder, LogToErrorLogFile, errorsLogPriority), // json格式 打印 lv>=error 信息到 error log
		zapcore.NewCore(jsonEncoder, logToFullLogFile, fullLogPriority),    // json格式 打印 lv>=debug 信息到 full log
		zapcore.NewCore(consoleEncoder, LogToConsole, fullLogPriority),     // console格式 打印 lv>=debug 到 os.stdout
	)

	// 如果需要用到 caller 和 Stacktrace 需要在这里添加
	Logger = zap.New(core) // 不使用 caller 和 Stacktrace
	//Logger = zap.New(core, zap.AddCaller()) // 只使用 caller
	//Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)) // 使用 caller 和 Stacktrace

	// 必须 sync
	defer Logger.Sync()

	return nil
}
