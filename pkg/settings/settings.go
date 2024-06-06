package settings

import (
	"github.com/tak-sh/tak/pkg/utils/fileutils"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

var Default Settings

const (
	DirName            = ".tak"
	LogDirName         = "logs"
	LogFilename        = "tak.log"
	AccountsDirName    = "accounts"
	AccountDataDirName = "data"
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
	HTMLDir(accountName string) string
	LogFile() string
	IsDev() bool
	AccountDir(account string) string
	AccountDataDir(account string) string
}

type settings struct {
	Version string
	SetDir  string
	Cache   string
	Log     string
}

func (s *settings) HTMLDir(accountName string) string {
	basePath := fileutils.FindUpwardFrom("accounts", "", "")
	if basePath == "" {
		basePath, _ = os.Getwd()
	}

	fp := filepath.Join(basePath, accountName, "html")
	_ = os.MkdirAll(fp, os.ModePerm)

	return fp
}

func (s *settings) AccountDataDir(account string) string {
	d := filepath.Join(s.AccountDir(account), AccountDataDirName)
	_ = os.MkdirAll(d, os.ModePerm)
	return d
}

func (s *settings) AccountDir(account string) string {
	d := filepath.Join(s.SettingsDir(), AccountsDirName, account)
	_ = os.MkdirAll(d, os.ModePerm)
	return d
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
