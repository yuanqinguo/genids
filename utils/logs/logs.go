package logs

import (
	"genids/utils"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"time"
)

var LogSystem = logrus.New()

func InitLog(systemLogPath string, logMaxAge int) error {
	err := createLogPath(systemLogPath)
	if err != nil {
		return err
	}

	// error log
	error_log, err := rotatelogs.New(systemLogPath, rotatelogs.WithMaxAge(time.Duration(logMaxAge)*24*time.Hour))
	if err != nil {
		return err
	}

	//系统运行日志
	LogSystem.SetOutput(error_log)
	LogSystem.SetFormatter(&logrus.JSONFormatter{TimestampFormat: utils.SysTimeform})
	LogSystem.SetReportCaller(false)

	return nil
}

func createLogPath(path string) error {
	dirPath, err := filepath.Abs(filepath.Dir(path))
	if err != nil {
		return err
	}

	_, err = os.Stat(dirPath)
	if err == nil {
		return nil
	}

	if os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, os.ModePerm)
		return err
	}

	return err
}
