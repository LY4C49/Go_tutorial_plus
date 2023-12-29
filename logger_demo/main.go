package main

import (
	"net/http"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

//	func SetupLogger() {
//		logFileLocation, _ := os.OpenFile("./test.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0744)
//		log.SetOutput(logFileLocation)
//	}
var logger *zap.Logger
var sugarLogger *zap.SugaredLogger

//	func InitLogger() {
//		logger, _ = zap.NewProduction()
//	}
func simpleHttpGet(url string) {
	resp, err := http.Get(url)
	if err != nil {
		sugarLogger.Error(
			"Error fetching url..",
			zap.String("url", url),
			zap.Error(err))
	} else {
		sugarLogger.Info("Success..",
			zap.String("statusCode", resp.Status),
			zap.String("url", url))
		resp.Body.Close()
	}
}

func InitLogger() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel) //core： Encoder, WriteSync, 日志级别

	logger = zap.New(core, zap.AddCaller())
	sugarLogger = logger.Sugar()
}

func getEncoder() zapcore.Encoder {
	//return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()) //json格式
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncodeLevel = zapcore.CapitalLevelEncoder
	//return zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig()) //普通的终端输出流类似
	return zapcore.NewConsoleEncoder(config)
}

func getLogWriter() zapcore.WriteSyncer {
	//file, _ := os.Create("./test.log")
	//file, _ := os.OpenFile("./test.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0744)

	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./test.log",
		MaxSize:    1,     //MB
		MaxBackups: 5,     //备份数量
		MaxAge:     30,    //最大备份天数（最多保存多少天的）
		Compress:   false, //是否压缩
	}
	return zapcore.AddSync(lumberJackLogger)
}

func main() {
	InitLogger()
	defer logger.Sync() //退出时将缓存的数据刷入文件
	for {
		simpleHttpGet("www.google.com")
		simpleHttpGet("http://www.google.com")
	}

}
