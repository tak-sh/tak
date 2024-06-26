package settings

import (
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	"github.com/tak-sh/tak/generated/go/api/settings/v1beta1"
	"github.com/tak-sh/tak/pkg/protoenc"
	"github.com/tak-sh/tak/pkg/utils/fileutils"
	"github.com/tak-sh/tak/pkg/utils/ptr"
	"google.golang.org/protobuf/proto"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

var Default Settings

const (
	DefaultSettingsEnvVar = "TAK_SETTINGS"
	DirName               = ".tak"
	LogDirName            = "logs"
	LogFilename           = "tak.log"
	FileName              = "tak.yaml"
	ChromeUserDataDir     = "chrome_user_data"
)

var (
	UserHomeDir, _  = os.UserHomeDir()
	UserCacheDir, _ = os.UserCacheDir()

	MaxLogFileSize = 10 * 1024 * 1024 // 10Mb
)

func init() {
	Default = New("", filepath.Join(UserHomeDir, DirName), UserCacheDir)
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
		SetDir:  configDir,
		Cache:   cacheDir,
		Log:     filepath.Join(logDir, LogFilename),
		Config:  koanf.New("."),
	}
	out.SetFS = os.DirFS(out.SetDir)

	_ = os.MkdirAll(logDir, os.ModePerm)
	_ = os.MkdirAll(out.SetDir, os.ModePerm)

	return out
}

type Settings interface {
	SettingsDir() string
	ScreenshotDir() string
	LogFile() string
	IsDev() bool
	LoadFile() (*v1beta1.Settings, error)
	Get() *v1beta1.Settings
}

type settings struct {
	Version string
	SetDir  string
	SetFS   fs.FS
	Cache   string
	Log     string
	Config  *koanf.Koanf

	userHome string
	sett     *v1beta1.Settings
}

func (s *settings) Get() *v1beta1.Settings {
	if s.sett == nil {
		def := newDefaultSettings(s.SetDir)
		s.sett, _ = s.LoadFile()
		proto.Merge(s.sett, def)

		if s.sett.ChromeDataDirectory != nil {
			s.sett.ChromeDataDirectory = ptr.Ptr(fileutils.ExpandHome(s.userHome, *s.sett.ChromeDataDirectory))
		}

		_ = s.Config.Load(structs.Provider(s.sett, "json"), nil)

		_ = s.Config.Load(env.Provider("TAK_", ".", func(s string) string {
			return strings.ReplaceAll(strings.ToLower(strings.TrimPrefix(s, "TAK_")), "_", ".")
		}), nil)

		_ = s.Config.Unmarshal("", s.sett)
	}
	return s.sett
}

func newDefaultSettings(settingsDir string) *v1beta1.Settings {
	return &v1beta1.Settings{
		ChromeDataDirectory: ptr.Ptr(filepath.Join(settingsDir, ChromeUserDataDir)),
	}
}

func (s *settings) LoadFile() (*v1beta1.Settings, error) {
	out := new(v1beta1.Settings)
	return out, protoenc.UnmarshalFile(out, FileName, s.SetFS)
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
