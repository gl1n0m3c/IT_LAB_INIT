package log

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"strings"
)

type Logs struct {
	InfoLogger  *zerolog.Logger
	ErrorLogger *zerolog.Logger
}

func UnitFormatter() {
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05"

	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		path := strings.Split(file, "IT_LAB_INIT")
		return fmt.Sprintf("%s:%d", fmt.Sprintf("IT_LAB_INIT%s", path[len(path)-1]), line)
	}
}

func InitLoggers() (*Logs, *os.File, *os.File) {
	UnitFormatter()

	loggerInfoFile, err := os.OpenFile("log/info.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		panic("Error opening info log file")
	}

	loggerErrorFile, err := os.OpenFile("log/error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		panic("Error opening error log file")
	}

	infoLogger := zerolog.New(loggerInfoFile).With().Timestamp().Caller().Logger()
	errorLogger := zerolog.New(loggerErrorFile).With().Timestamp().Caller().Logger()

	log := &Logs{
		InfoLogger:  &infoLogger,
		ErrorLogger: &errorLogger,
	}

	return log, loggerInfoFile, loggerErrorFile
}
