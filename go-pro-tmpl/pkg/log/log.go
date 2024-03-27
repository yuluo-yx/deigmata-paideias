package log

import (
	"io"
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/writer"
)

var SysLog = logrus.New()

func InitLogger(logPath, logLevel string, consoleLog bool) {

	if logPath == "" {
		logPath = "./apt-registry.log"
	}

	_, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		logrus.Fatalln("open log file err ", err)
	}

	sysLog := logrus.New()
	innerLogLevel, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logrus.Fatalln("config log level error", err)
		panic("config log level error!")
	}
	sysLog.SetLevel(innerLogLevel)

	var out io.Writer
	if consoleLog {
		out = os.Stdout
	} else {
		out = io.MultiWriter()
	}
	sysLog.Out = io.MultiWriter()

	logWriter, err := rotatelogs.New(
		logPath+".%Y%m%d.log",
		rotatelogs.WithLinkName(logPath),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(1*time.Hour),
	)
	sysLog.SetFormatter(&logrus.JSONFormatter{})
	writerMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}
	sysLog.AddHook(lfshook.NewHook(writerMap, &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	}))
	sysLog.AddHook(&writer.Hook{
		Writer: out,
		LogLevels: []logrus.Level{
			logrus.InfoLevel,
			logrus.FatalLevel,
			logrus.DebugLevel,
			logrus.WarnLevel,
			logrus.ErrorLevel,
			logrus.PanicLevel,
		},
	})

	sysLog = sysLog

}
