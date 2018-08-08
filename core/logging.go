package core

import (
	"fmt"
	"io"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/fatih/color"
	"github.com/gen0cide/gscript/logger"
	"github.com/gen0cide/gscript/logger/standard"
)

var (
	// Logger is a global singleton logger used by Laforge
	Logger logger.Logger

	boldgreen  = color.New(color.FgHiGreen, color.Bold).SprintfFunc()
	boldwhite  = color.New(color.FgHiWhite, color.Bold).SprintfFunc()
	boldred    = color.New(color.FgHiRed, color.Bold).SprintfFunc()
	boldyellow = color.New(color.FgHiYellow, color.Bold).SprintfFunc()
	boldcyan   = color.New(color.FgHiCyan, color.Bold).SprintfFunc()
	boldb      = color.New(color.FgHiBlue, color.Bold).SprintfFunc()
	boldg      = color.New(color.FgHiGreen, color.Bold).SprintfFunc()
	boldw      = color.New(color.FgHiWhite, color.Bold).SprintfFunc()
	boldr      = color.New(color.FgHiRed, color.Bold).SprintfFunc()
	boldy      = color.New(color.FgHiYellow, color.Bold).SprintfFunc()
	boldc      = color.New(color.FgHiCyan, color.Bold).SprintfFunc()
	boldm      = color.New(color.FgHiMagenta, color.Bold).SprintfFunc()
	britw      = color.New(color.FgHiWhite).SprintfFunc()
	normb      = color.New(color.FgBlue).SprintfFunc()
	nocol      = color.New(color.Reset).SprintfFunc()
	boldblue   = color.New(color.FgHiBlue, color.Bold).SprintfFunc()
	plainblue  = color.New(color.FgHiBlue)
)

var (
	intLogger    *internalLogger
	globalProg   = `LAFORGE`
	startName    = `cli`
	defaultLevel = logrus.WarnLevel
)

type internalLogger struct {
	internal *logrus.Logger
	writer   *logWriter
	prog     string
	context  string
}

type logWriter struct {
	Name string
	Prog string
}

func init() {
	base := standard.NewStandardLogger(nil, "LAFORGE", "cli", false, false)
	baseSL := base.Logger
	writer := &logWriter{Prog: globalProg, Name: startName}
	baseSL.Out = writer
	baseSL.SetLevel(defaultLevel)
	logger := &internalLogger{
		internal: baseSL,
		writer:   writer,
	}
	intLogger = logger
	Logger = base
}

// SetLogLevel allows you to override the logging level for the Laforge global logger
func SetLogLevel(level string) {
	switch strings.ToLower(level) {
	case "debug":
		intLogger.internal.SetLevel(logrus.DebugLevel)
	case "info":
		intLogger.internal.SetLevel(logrus.InfoLevel)
	case "warn":
		intLogger.internal.SetLevel(logrus.WarnLevel)
	case "error":
		intLogger.internal.SetLevel(logrus.ErrorLevel)
	case "fatal":
		intLogger.internal.SetLevel(logrus.FatalLevel)
	}
}

// SetLogName allows you to override the log name parameter (part after LAFORGE in log output)
func SetLogName(name string) {
	intLogger.writer.Name = name
}

func (w *logWriter) Write(p []byte) (int, error) {
	output := fmt.Sprintf(
		"%s%s%s%s%s %s",
		boldwhite("["),
		boldblue(w.Prog),
		boldwhite(":"),
		boldgreen(strings.ToLower(w.Name)),
		boldwhite("]"),
		string(p),
	)
	written, err := io.Copy(color.Output, strings.NewReader(output))
	return int(written), err
}
