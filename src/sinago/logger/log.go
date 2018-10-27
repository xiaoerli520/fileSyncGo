package logger

import (
	"os"
	"log"
	"sinago/utils"
)


const (
	INFO = iota
	NOTICE
	WARN
	ERROR
)

type Logger struct {

	path string

	level int

}


func SetLogger(logPath string, level int) (*Logger) {

	var logger = new(Logger)

	logger.path = logPath

	logger.level = level

	return logger

}

func (logger *Logger) Err(content string) {
	if ERROR >= logger.level {
		var time = utils.GetDate()
		var fileName = logger.path + "error-" + time + ".log"
		logFile, err := os.OpenFile(fileName, os.O_APPEND | os.O_CREATE | os.O_RDWR, 0666)
		if err != nil {
			panic(err)
		}
		dbgLog := log.New(logFile, "[ERROR] ", log.LstdFlags)
		dbgLog.Println(content)
	}
}

func (logger *Logger) Warn(content string) {
	if WARN >= logger.level {
		var time = utils.GetDate()
		var fileName = logger.path + "warn-" + time + ".log"
		logFile, err := os.OpenFile(fileName, os.O_APPEND | os.O_CREATE | os.O_RDWR, 0666)
		if err != nil {
			panic(err)
		}
		dbgLog := log.New(logFile, "[WARN] ", log.LstdFlags)
		dbgLog.Println(content)
	}
}

func (logger *Logger) Notice(content string) {
	if NOTICE >= logger.level {
		var time = utils.GetDate()
		var fileName = logger.path + "notice-" + time + ".log"
		logFile, err := os.OpenFile(fileName, os.O_APPEND | os.O_CREATE | os.O_RDWR, 0666)
		if err != nil {
			panic(err)
		}
		dbgLog := log.New(logFile, "[NOTICE] ", log.LstdFlags)
		dbgLog.Println(content)
	}
}

func (logger *Logger) Info(content string) {
	if INFO >= logger.level {
		var time = utils.GetDate()
		var fileName = logger.path + "info-" + time + ".log"
		logFile, err := os.OpenFile(fileName, os.O_APPEND | os.O_CREATE | os.O_RDWR, 0666)
		if err != nil {
			panic(err)
		}
		dbgLog := log.New(logFile, "[INFO] ", log.LstdFlags)
		dbgLog.Println(content)
	}

}


