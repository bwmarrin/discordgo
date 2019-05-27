package discordgo

import (
	"fmt"
	"log"
	"runtime"
	"strings"
)

type LogLevel int

const (

	// LogError level is used for critical errors that could lead to data loss
	// or panic that would not be returned to a calling function.
	LogError LogLevel = iota

	// LogWarning level is used for very abnormal events and errors that are
	// also returned to a calling function.
	LogWarning

	// LogInformational level is used for normal non-error activity
	LogInformational

	// LogDebug level is for very detailed non-error activity.  This is
	// very spammy and will impact performance.
	LogDebug
)

// Logger provides a generic logger interface to make it easy to inject custom loggers
type Logger interface {
	Log(level LogLevel, format string, v ...interface{})
}

// DefaultLogger returns a logger implementing the interface
func DefaultLogger(l LogLevel) Logger {
	return defaultLogger{
		level: l,
	}
}

type defaultLogger struct {
	level LogLevel
}

func (d defaultLogger) Log(l LogLevel, format string, v ...interface{}) {
	if l < d.level {
		return
	}

	pc, file, line, _ := runtime.Caller(2)

	files := strings.Split(file, "/")
	file = files[len(files)-1]

	name := runtime.FuncForPC(pc).Name()
	fns := strings.Split(name, ".")
	name = fns[len(fns)-1]

	msg := fmt.Sprintf(format, v...)

	log.Printf("[DG%d] %s:%d:%s() %s\n", l, file, line, name, msg)
}

func (s *Session) log(msgLevel LogLevel, format string, v ...interface{}) {
	s.Logger.Log(msgLevel, format, v...)
}

func (v *VoiceConnection) log(msgLevel LogLevel, format string, a ...interface{}) {
	v.session.Logger.Log(msgLevel, format, a...)
}
