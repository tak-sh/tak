package cli

import (
	"github.com/tak-sh/tak/pkg/contexts"
	"github.com/tak-sh/tak/pkg/settings"
	"github.com/urfave/cli/v2"
	"log/slog"
	"os"
	"os/signal"
	"strings"
)

func New(version string) *cli.App {
	app := &cli.App{
		Name:    "tak",
		Usage:   "Local-first finance tracker CLI.",
		Version: version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "settings",
				Usage: "The path to the settings folder. Created if it does not exist.",
				Value: settings.Default.SettingsDir(),
				EnvVars: []string{
					"TAK_SETTINGS",
				},
				Action: func(context *cli.Context, s string) error {
					_, err := os.Stat(s)
					if err != nil {
						_ = os.MkdirAll(s, os.ModePerm)
					}
					return nil
				},
			},
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "Set this to enable debug mode when reporting an issue.",
			},
			&cli.StringSliceFlag{
				Name:   "log_outputs",
				Usage:  "List of either file paths, stdout, or stderr",
				Hidden: true,
			},
		},
		Description: "Private, local-first, cli-based finance tracker with automatic account syncing.",
		Before: func(context *cli.Context) error {
			set := settings.New(version, context.String("settings"), settings.UserCacheDir)

			lvl := slog.LevelInfo
			if context.Bool("debug") {
				lvl = slog.LevelDebug
			}

			var outputs []string
			if lo := context.StringSlice("log_outputs"); len(lo) > 0 {
				outputs = lo
			} else {
				outputs = append(outputs, set.LogFile())
			}

			logger := settings.NewLogger(lvl, outputs)
			context.Context = settings.WithSettings(context.Context, set)
			context.Context = contexts.WithLogger(context.Context, logger)

			context.Context, _ = signal.NotifyContext(context.Context, os.Interrupt, os.Kill)

			logger.Info("Running command.", slog.String("cmd", strings.Join(context.Args().Slice(), " ")))
			return nil
		},
		Commands: []*cli.Command{
			NewAccountCommand(),
		},
		Suggest: true,
	}

	return app
}
