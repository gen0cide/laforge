package cli

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

	// Boldgreen defines a color printer
	Boldgreen = color.New(color.FgHiGreen, color.Bold).SprintfFunc()

	// Boldwhite defines a color printer
	Boldwhite = color.New(color.FgHiWhite, color.Bold).SprintfFunc()

	// Boldred defines a color printer
	Boldred = color.New(color.FgHiRed, color.Bold).SprintfFunc()

	// Boldyellow defines a color printer
	Boldyellow = color.New(color.FgHiYellow, color.Bold).SprintfFunc()

	// Boldcyan defines a color printer
	Boldcyan = color.New(color.FgHiCyan, color.Bold).SprintfFunc()

	// Boldb defines a color printer
	Boldb = color.New(color.FgHiBlue, color.Bold).SprintfFunc()

	// Boldg defines a color printer
	Boldg = color.New(color.FgHiGreen, color.Bold).SprintfFunc()

	// Boldw defines a color printer
	Boldw = color.New(color.FgHiWhite, color.Bold).SprintfFunc()

	// Boldr defines a color printer
	Boldr = color.New(color.FgHiRed, color.Bold).SprintfFunc()

	// Boldy defines a color printer
	Boldy = color.New(color.FgHiYellow, color.Bold).SprintfFunc()

	// Boldc defines a color printer
	Boldc = color.New(color.FgHiCyan, color.Bold).SprintfFunc()

	// Boldm defines a color printer
	Boldm = color.New(color.FgHiMagenta, color.Bold).SprintfFunc()

	// Britw defines a color printer
	Britw = color.New(color.FgHiWhite).SprintfFunc()

	// Normb defines a color printer
	Normb = color.New(color.FgBlue).SprintfFunc()

	// Nocol defines a color printer
	Nocol = color.New(color.Reset).SprintfFunc()

	// Boldblue defines a color printer
	Boldblue = color.New(color.FgHiBlue, color.Bold).SprintfFunc()

	// Plainblue defines a color printer
	Plainblue = color.New(color.FgHiBlue)
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
		Boldwhite("["),
		Boldblue(w.Prog),
		Boldwhite(":"),
		Boldgreen(strings.ToLower(w.Name)),
		Boldwhite("]"),
		string(p),
	)
	written, err := io.Copy(color.Output, strings.NewReader(output))
	return int(written), err
}
