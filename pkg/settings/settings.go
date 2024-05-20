package settings

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

var Default Settings

const (
	DirName     = ".tak"
	LogDirName  = "logs"
	LogFilename = "tak.log"
)

var (
	UserConfigDir, _ = os.UserConfigDir()
	UserHomeDir, _   = os.UserHomeDir()
	UserCacheDir, _  = os.UserCacheDir()

	MaxLogFileSize = 10 * 1024 * 1024 // 10Mb
)

func init() {
	Default = New("", UserConfigDir, UserCacheDir)
}

func NewLogger(level slog.Level, outputs []string) *slog.Logger {
	w := make([]io.Writer, 0, len(outputs))
	for _, v := range outputs {
		switch v {
		case "stdout":
			w = append(w, os.Stdout)
		case "stderr":
			w = append(w, os.Stderr)
		default:
			w = append(w, &lumberjack.Logger{
				Filename:   v,
				MaxSize:    MaxLogFileSize,
				MaxBackups: 3,
			})
		}
	}

	logger := slog.New(slog.NewJSONHandler(io.MultiWriter(w...), &slog.HandlerOptions{
		AddSource: true,
		Level:     level,
	}))

	return logger
}

func New(version, configDir, cacheDir string) Settings {
	logDir := filepath.Join(cacheDir, LogDirName)
	out := &settings{
		Version: version,
		SetDir:  filepath.Join(configDir, DirName),
		Cache:   cacheDir,
		Log:     filepath.Join(logDir, LogFilename),
	}

	_ = os.MkdirAll(logDir, os.ModePerm)
	_ = os.MkdirAll(out.SetDir, os.ModePerm)

	return out
}

type Settings interface {
	CacheDir() string
	SettingsDir() string
	ScreenshotDir() string
	LogFile() string
	IsDev() bool
}

type settings struct {
	Version string
	SetDir  string
	Cache   string
	Log     string
}

func (s *settings) LogFile() string {
	return s.Log
}

func (s *settings) ScreenshotDir() string {
	return filepath.Join(s.SettingsDir(), "screenshots")
}

func (s *settings) IsDev() bool {
	return s.Version == ""
}

func (s *settings) CacheDir() string {
	return s.Cache
}

func (s *settings) SettingsDir() string {
	return s.SetDir
}
