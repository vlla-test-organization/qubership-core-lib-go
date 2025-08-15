package logging

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	lock "github.com/viney-shih/go-lock"
	"github.com/vlla-test-organization/qubership-core-lib-go/v3/configloader"
)

// List of predefined log Levels
const (
	LvlCrit Lvl = iota
	LvlError
	LvlWarn
	LvlInfo
	LvlDebug
)

// Lvl is a type for predefined log levels.
type Lvl int

type LogLevels map[string]string

type logger struct {
	maxLvl    Lvl
	name      string
	logFormat func(r *Record) []byte
	// mutex for sync log
	mu              *lock.ChanMutex
	rwLockForMaxLvl sync.RWMutex
}

var (
	registeredLoggers sync.Map
	once              sync.Once
	globalLogFormat   = DefaultFormat.format
	defaultMaxLevel   = LvlInfo
	envNameRegexp     = regexp.MustCompile(`[^a-zA-Z0-9_]+`)
)

// A Logger writes key/value pairs to a Handler
type Logger interface {
	GetLevel() Lvl

	// SetLevel updates the logger to set specific max level to write for
	SetLevel(maxLvl Lvl)
	SetLogFormat(logFormat func(r *Record) []byte)
	SetMessageFormat(fn messageFmt)

	// Log a Message at the given level
	Debug(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	DebugC(ctx context.Context, format string, args ...interface{})

	Info(format string, args ...interface{})
	Infof(format string, args ...interface{})
	InfoC(ctx context.Context, format string, args ...interface{})

	Warn(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	WarnC(ctx context.Context, format string, args ...interface{})

	Error(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	ErrorC(ctx context.Context, format string, args ...interface{})

	Panic(format string, args ...interface{})
	Panicf(format string, args ...interface{})
	PanicC(ctx context.Context, format string, args ...interface{})
}

func watch() {
	_, err := configloader.Subscribe(func(event configloader.Event) error {
		if event.Type == configloader.InitedEventT || event.Type == configloader.RefreshedEventT {
			rootLevel := defineRootLvl(defaultMaxLevel.String())
			registeredLoggers.Range(func(key, value interface{}) bool {
				logLevel := definePackageLvl(value.(*logger).name)
				if logLevel == "" {
					logLevel = rootLevel
				}
				lvl, _ := lvlFromString(logLevel)
				if currLogger, ok := registeredLoggers.Load(key); ok {
					currLogger.(*logger).SetLevel(lvl)
					registeredLoggers.Store(key, currLogger)
				}
				return true
			})
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func readLvlFromConfig(pkg string) string {
	if configloader.IsConfigLoaderInited() {
		packageLvl := definePackageLvl(pkg)
		if packageLvl != "" {
			return packageLvl
		}
		defaultLvl := defineRootLvl(defaultMaxLevel.String())
		return defaultLvl
	} else {
		pkgEnvSuffix := envNameRegexp.ReplaceAllString(strings.ToUpper(pkg), "_")
		envHierarchy := [4]string{
			"LOGGING_LEVEL_" + pkgEnvSuffix,
			"LOG_LEVEL_PACKAGE_" + pkgEnvSuffix,
			"LOGGING_LEVEL_ROOT",
			"LOG_LEVEL",
		}
		for _, env := range envHierarchy {
			if res, isExist := os.LookupEnv(env); isExist {
				return res
			}
		}
		return defaultMaxLevel.String()
	}
}

func definePackageLvl(pkg string) string {
	logPkgString := fmt.Sprintf("logging.level.%s", pkg)
	level := configloader.GetOrDefaultString(logPkgString, "")
	if level == "" {
		logPkgString := fmt.Sprintf("log.level.package.%s", pkg)
		level = configloader.GetOrDefaultString(logPkgString, "")
	}
	return level
}

func defineRootLvl(defaultValue string) string {
	logRootString := "logging.level.root"
	level := configloader.GetOrDefaultString(logRootString, "")
	if level == "" {
		logRootString = "logging.level.ROOT"
		level = configloader.GetOrDefaultString(logRootString, "")
	}
	if level == "" {
		logRootString = "log.level"
		level = configloader.GetOrDefaultString(logRootString, defaultValue)
	}
	return level
}

func GetLogger(name string) Logger {
	once.Do(watch)
	if l, ok := registeredLoggers.Load(name); ok {
		return l.(*logger)
	}
	l := new(logger)
	l.name = name
	l.mu = lock.NewChanMutex()
	maxLvl := readLvlFromConfig(name)
	if lvl, ok := lvlFromString(maxLvl); ok {
		l.maxLvl = lvl
	} else {
		l.maxLvl = defaultMaxLevel
		l.Warn("wrong log level logFormat: %s, falling to '"+defaultMaxLevel.String()+"'", maxLvl)
	}
	registeredLoggers.Store(name, l)
	return l
}

func GetLogLevels() LogLevels {
	logLevels := make(LogLevels)
	rootLvl, _ := lvlFromString(readLvlFromConfig(""))
	logLevels["ROOT"] = strings.ToUpper(rootLvl.String())
	registeredLoggers.Range(func(key, value any) bool {
		logLevels[key.(string)] = strings.ToUpper(value.(Logger).GetLevel().String())
		return true
	})
	return logLevels
}

func (l *logger) GetLevel() Lvl {
	return l.readMaxLvlWithRLock()
}

func (l *logger) SetLevel(maxLvl Lvl) {
	l.rwLockForMaxLvl.Lock()
	l.maxLvl = maxLvl
	l.rwLockForMaxLvl.Unlock()
}

func (l *logger) SetLogFormat(logFormat func(r *Record) []byte) {
	l.logFormat = logFormat
}

func (l *logger) SetMessageFormat(fn messageFmt) {
	logFormat := defaultFormat{}
	logFormat.SetMessageFormat(fn)
	l.logFormat = logFormat.format
}

func SetLogFormat(format func(r *Record) []byte) {
	globalLogFormat = format
}

func (l *logger) format(r *Record) []byte {
	var bytes []byte
	if l.logFormat != nil {
		bytes = l.logFormat(r)
	} else {
		bytes = globalLogFormat(r)
	}
	return bytes
}

// Returns the name of a Lvl
func (l Lvl) String() string {
	switch l {
	case LvlDebug:
		return "debug"
	case LvlInfo:
		return "info"
	case LvlWarn:
		return "warn"
	case LvlError:
		return "error"
	case LvlCrit:
		return "fatal"
	default:
		panic("bad level")
	}
}

// lvlFromString returns the appropriate Lvl from a string name.
// Useful for parsing command line args and configuration files.
func lvlFromString(lvlString string) (Lvl, bool) {
	lvlLowCase := strings.ToLower(lvlString)
	switch lvlLowCase {
	case "debug":
		return LvlDebug, true
	case "info":
		return LvlInfo, true
	case "warn":
		return LvlWarn, true
	case "error":
		return LvlError, true
	case "fatal":
		return LvlCrit, true
	default:
		return defaultMaxLevel, false
	}
}

func setLogLevel(lvl string, packageName string) error {
	desiredLogger, loggerIsFound := registeredLoggers.Load(packageName)
	if !loggerIsFound {
		return errors.New("Logger with name " + packageName + " not found")
	}
	if newLevel, isLevelExists := lvlFromString(lvl); isLevelExists {
		desiredLogger.(*logger).SetLevel(newLevel)
		registeredLoggers.Store(packageName, desiredLogger)
		return nil
	}
	return errors.New("Can't set lvl " + lvl + " for logger " + packageName)
}

// mutexes are used to guarantee that
// only a single Log operation can proceed at a Time. It's necessary
// for thread-safe concurrent writes.
func (l *logger) log(ctx context.Context, lvl Lvl, wr io.Writer, sFormat string, args ...interface{}) error {
	if lvl <= l.readMaxLvlWithRLock() {
		r := &Record{
			PackageName: l.name,
			Time:        time.Now(),
			Lvl:         lvl,
			Message:     fmt.Sprintf(sFormat, args...),
			Ctx:         ctx,
		}

		if l.mu.TryLockWithTimeout(5 * time.Second) {
			defer l.mu.Unlock()
			_, err := wr.Write(l.format(r))
			return err
		} else {
			defaultFormatForError := new(defaultFormat)
			printErrorLogInDefaultFormat(wr, *r, *defaultFormatForError)
			_, err := wr.Write(defaultFormatForError.format(r))
			return err
		}
	}
	return nil
}

func (l *logger) Debug(format string, args ...interface{}) {
	l.log(nil, LvlDebug, os.Stdout, format, args...)
}
func (l *logger) Debugf(format string, args ...interface{}) {
	l.Debug(format, args...)
}
func (l *logger) DebugC(ctx context.Context, format string, args ...interface{}) {
	l.log(ctx, LvlDebug, os.Stdout, format, args...)
}
func (l *logger) Info(format string, args ...interface{}) {
	l.log(nil, LvlInfo, os.Stdout, format, args...)
}
func (l *logger) Infof(format string, args ...interface{}) {
	l.Info(format, args...)
}
func (l *logger) InfoC(ctx context.Context, format string, args ...interface{}) {
	l.log(ctx, LvlInfo, os.Stdout, format, args...)
}
func (l *logger) Warn(format string, args ...interface{}) {
	l.log(nil, LvlWarn, os.Stdout, format, args...)
}
func (l *logger) Warnf(format string, args ...interface{}) {
	l.Warn(format, args...)
}
func (l *logger) WarnC(ctx context.Context, format string, args ...interface{}) {
	l.log(ctx, LvlWarn, os.Stdout, format, args...)
}
func (l *logger) Error(format string, args ...interface{}) {
	l.log(nil, LvlError, os.Stdout, format, args...)
}
func (l *logger) Errorf(format string, args ...interface{}) {
	l.Error(format, args...)
}
func (l *logger) ErrorC(ctx context.Context, format string, args ...interface{}) {
	l.log(ctx, LvlError, os.Stdout, format, args...)
}
func (l *logger) Panic(format string, args ...interface{}) {
	l.log(nil, LvlCrit, os.Stdout, format, args...)
	panic(fmt.Sprintf(format, args...))
}
func (l *logger) Panicf(format string, args ...interface{}) {
	l.Panic(format, args...)
	panic(fmt.Sprintf(format, args...))
}
func (l *logger) PanicC(ctx context.Context, format string, args ...interface{}) {
	l.log(ctx, LvlCrit, os.Stdout, format, args...)
	panic(fmt.Sprintf(format, args...))
}

func (l *logger) readMaxLvlWithRLock() Lvl {
	l.rwLockForMaxLvl.RLock()
	defer l.rwLockForMaxLvl.RUnlock()
	return l.maxLvl
}

func printErrorLogInDefaultFormat(wr io.Writer, r Record, DefaultFormatForError defaultFormat) {
	wr.Write(DefaultFormatForError.format(&Record{
		PackageName: constructCallerValue(r.Ctx, r.PackageName),
		Time:        time.Time{},
		Lvl:         1,
		Message: "Possibility of deadlock or circular dependency. " +
			"Perhaps, you use wrong custom log format",
		Ctx: nil,
	}))
}
