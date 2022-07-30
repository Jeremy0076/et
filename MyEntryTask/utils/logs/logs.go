package logs

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

// 1.日志级别
// todo 2.日志分割 3.异步打印日志

const (
	kLogPath = "../utils/logs/"
	kLogExt  = "_log.log"
	kLayout  = "2006-01-02"
)

const (
	LevelError = iota
	LevelWarning
	LevelInfo
	LevelDebug
)

type Logger struct {
	level int
	err   *log.Logger
	warn  *log.Logger
	info  *log.Logger
	debug *log.Logger
}

func NewLogger(level int) (*Logger, error) {
	ll := &Logger{}
	ll.level = level
	logPath := kLogPath + time.Now().Format(kLayout) + kLogExt
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("open file err: [%v]\n", err)
		return nil, err
	}

	w := io.Writer(file)

	// 打印date:hh:mm:ss 行号xxx.go
	flag := log.LstdFlags | log.Lshortfile
	ll.err = log.New(w, "[E]", flag)
	ll.warn = log.New(w, "[W]", flag)
	ll.info = log.New(w, "[I]", flag)
	ll.debug = log.New(w, "[D]", flag)
	return ll, nil
}

func (ll *Logger) Error(format string, v ...interface{}) {
	if LevelError > ll.level {
		return
	}
	msg := fmt.Sprintf(format, v...)
	//ll.err.Printf(msg)
	ll.err.Output(2, msg)
}

func (ll *Logger) Warn(format string, v ...interface{}) {
	if LevelError > ll.level {
		return
	}
	msg := fmt.Sprintf(format, v...)
	//ll.warn.Printf(msg)
	ll.warn.Output(2, msg)
}

func (ll *Logger) Info(format string, v ...interface{}) {
	if LevelError > ll.level {
		return
	}
	msg := fmt.Sprintf(format, v...)
	//ll.info.Printf(msg)
	ll.info.Output(2, msg)
}

func (ll *Logger) Debug(format string, v ...interface{}) {
	if LevelError > ll.level {
		return
	}
	msg := fmt.Sprintf(format, v...)
	//ll.debug.Printf(msg)
	ll.debug.Output(2, msg)
}
