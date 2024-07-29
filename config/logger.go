package config

import (
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func InitLogger() {
    Logger = logrus.New()
    
    path := "logs/application_log"
    
    writer, err := rotatelogs.New(
        path+".%Y%m%d%H%M",
        rotatelogs.WithLinkName(path),
        rotatelogs.WithMaxAge(7*24*time.Hour),
        rotatelogs.WithRotationTime(24*time.Hour),
    )
    if err != nil {
        logrus.Errorf("Failed to initialize log file: %v", err)
        return
    }

    lfHook := lfshook.NewHook(lfshook.WriterMap{
        logrus.InfoLevel:  writer,
        logrus.ErrorLevel: writer,
        logrus.FatalLevel: writer,
        logrus.PanicLevel: writer,
    }, &logrus.TextFormatter{})
    
    Logger.AddHook(lfHook)

    Logger.SetOutput(os.Stdout)
    Logger.SetLevel(logrus.DebugLevel)
}
