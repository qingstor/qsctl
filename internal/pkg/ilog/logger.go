package ilog

import (
	"io"
	"os"

	"github.com/qingstor/log"
	"github.com/qingstor/log/level"
)

// InitLoggerWithDebug is util to init stderr logger with debug level or default
func InitLoggerWithDebug(d bool) *log.Logger {
	lvl := level.Warn // default level
	if d {
		lvl = level.Empty
	}
	return InitLoggerWithLevelAndWriter(lvl, os.Stderr)
}

// InitLoggerWithLevelAndWriter is util to init logger with given level and writer
func InitLoggerWithLevelAndWriter(l level.Level, w io.Writer) *log.Logger {
	e := log.ExecuteMatchWrite(
		// Only print log that level is higher than Debug.
		log.MatchHigherLevel(l),
		// Write into w.
		w,
	)
	tf, _ := log.NewText(&log.TextConfig{
		// Use unix timestamp nano for time
		TimeFormat: log.TimeFormatUnixNano,
		// Use upper case level
		LevelFormat: level.UpperCase,
		EntryFormat: "[{level}] - {time} {value}",
	})
	return log.New().WithExecutor(e).WithTransformer(tf)
}
