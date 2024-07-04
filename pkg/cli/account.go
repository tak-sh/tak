package cli

import (
	"context"
	"errors"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/chromedp/chromedp"
	"github.com/tak-sh/tak/pkg/contexts"
	"github.com/tak-sh/tak/pkg/debug"
	"github.com/tak-sh/tak/pkg/except"
	"github.com/tak-sh/tak/pkg/headless/engine"
	"github.com/tak-sh/tak/pkg/headless/script"
	"github.com/tak-sh/tak/pkg/headless/step/stepper"
	"github.com/tak-sh/tak/pkg/provider"
	"github.com/tak-sh/tak/pkg/renderer"
	"github.com/tak-sh/tak/pkg/settings"
	"github.com/tak-sh/tak/pkg/ui"
	"github.com/tak-sh/tak/pkg/ui/keyregistry"
	"github.com/urfave/cli/v2"
	"log/slog"
	"os"
	"time"
)

func NewAccountCommand() *cli.Command {
	cmd := &cli.Command{
		Name:    "account",
		Aliases: []string{"a", "acct"},
		Usage:   "Manage your accounts.",
		Subcommands: []*cli.Command{
			NewGetAccountCommand(),
			NewAddAccountCommand(),
			NewAccountSyncCommand(),
			NewDebugAccountCommand(),
		},
	}
	return cmd
}

func NewGetAccountCommand() *cli.Command {
	cmd := &cli.Command{
		Name:      "get",
		Aliases:   []string{"g"},
		Usage:     "List accounts or get a single account.",
		Args:      true,
		ArgsUsage: "The name of the account to get.",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "remote",
				Usage:   "List all available accounts in the GitHub repo.",
				Aliases: []string{"r"},
			},
		},
		Action: func(context *cli.Context) error {
			return nil
		},
	}

	return cmd
}

func NewAddAccountCommand() *cli.Command {
	cmd := &cli.Command{
		Name:        "add",
		Aliases:     []string{"a"},
		Usage:       "Add a new account.",
		Args:        true,
		ArgsUsage:   "The name of the account or path to an account file.",
		Description: "Accepts either the name of an account or a path. Paths must contain a '/' character e.g. ./chase.yaml. Run 'tak acct get -r' to see a list of available accounts.",
		Action: func(context *cli.Context) error {
			return nil
		},
	}

	return cmd
}

func NewDebugAccountCommand() *cli.Command {
	cmd := &cli.Command{
		Name:        "debug",
		Aliases:     []string{"de"},
		Usage:       "Debug an account manifest",
		Args:        true,
		ArgsUsage:   "The path to an account file.",
		Description: "Run and test every step of a new account in an interactive terminal window.",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "skip_login",
				Usage: "Skip the login step for the account provider.",
			},
			&cli.BoolFlag{
				Name:  "skip_download",
				Usage: "Skip the download transactions step for the account provider.",
			},
		},
		Action: func(cmd *cli.Context) error {
			ss := cmd.String("screenshots")
			fp := cmd.Args().First()
			logger := contexts.GetLogger(cmd.Context)

			mani, err := provider.LoadFile(fp)
			if err != nil {
				return err
			}

			prov, err := provider.New(mani)
			if err != nil {
				return err
			}

			str := renderer.NewStream()
			eq := engine.NewEventQueue()
			stpper := debug.NewFactory()

			scriptComp := ui.NewScriptComponent(prov.GetMetadata().GetName(), str, eq, logger)
			debugComp := ui.NewDebugComponent(stpper, scriptComp)
			app := ui.NewApp(debugComp)
			app.Help.Keys = keyregistry.DebugKeys
			p := tea.NewProgram(
				app,
				tea.WithContext(cmd.Context),
				tea.WithAltScreen(),
				tea.WithInput(os.Stdin),
				tea.WithOutput(os.Stdout),
			)

			uiCtx, err := ui.Run(cmd.Context, p)
			if err != nil {
				logger.Error("Failed to start the UI.", slog.String("err", err.Error()))
				return errors.Join(except.NewInternal("failed to start the UI"), err)
			}

			c, err := engine.NewContext(cmd.Context, str, engine.NewEvaluator(eq, 10*time.Second), engine.ContextOpts{
				ScreenshotDir: ss,
			})
			if err != nil {
				return err
			}

			chromeOpts := []chromedp.ExecAllocatorOption{
				chromedp.Flag("headless", false),
			}

			if cdd := settings.Default.Get().GetChromeDataDirectory(); cdd != "" {
				chromeOpts = append(chromeOpts, chromedp.UserDataDir(cdd))
			}

			acctCtx := prov.Run(c, stpper,
				provider.WithSkipLogin(cmd.Bool("skip_login")),
				provider.WithSkipDownloadTransactions(cmd.Bool("skip_download")),
				provider.WithScriptOpts(
					script.WithChromeOpts(chromeOpts...),
				),
			)

			select {
			case <-uiCtx.Done():
				return context.Cause(uiCtx)
			case <-acctCtx.Done():
				return context.Cause(acctCtx)
			}
		},
	}

	return cmd
}

func NewAccountSyncCommand() *cli.Command {
	cmd := &cli.Command{
		Name:        "sync",
		Aliases:     []string{"s"},
		Usage:       "Sync your accounts.",
		Description: "Sync the transactions from your account to your local machine. If no account is specified, all are synced.",
		Args:        true,
		ArgsUsage:   "The name of the account to sync.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "screenshots",
				Usage: "The directory to store screenshots of account syncs.",
				Value: settings.Default.ScreenshotDir(),
				Action: func(context *cli.Context, s string) error {
					_, err := os.Stat(s)
					if err != nil {
						_ = os.MkdirAll(s, os.ModePerm)
					}
					logger := contexts.GetLogger(context.Context)
					logger.Info("Storing screenshots.", slog.String("dir", s))
					return nil
				},
			},
			&cli.BoolFlag{
				Name:  "mfa",
				Usage: "Enables support for MFA by running chrome's GUI.",
			},
		},
		Action: func(cmd *cli.Context) error {
			ss := cmd.String("screenshots")
			fp := cmd.Args().First()
			logger := contexts.GetLogger(cmd.Context)

			acctRaw, err := provider.LoadFile(fp)
			if err != nil {
				return err
			}

			acct, err := provider.New(acctRaw)
			if err != nil {
				return err
			}

			str := renderer.NewStream()
			eq := engine.NewEventQueue()

			bubble := ui.NewScriptComponent(acctRaw.GetMetadata().GetName(), str, eq, logger)
			app := ui.NewApp(bubble)
			p := tea.NewProgram(
				app,
				tea.WithContext(cmd.Context),
				tea.WithAltScreen(),
				tea.WithInput(os.Stdin),
				tea.WithOutput(os.Stdout),
			)

			uiCtx, err := ui.Run(cmd.Context, p)
			if err != nil {
				logger.Error("Failed to start the UI.", slog.String("err", err.Error()))
				return errors.Join(except.NewInternal("failed to start the UI"), err)
			}

			c, err := engine.NewContext(cmd.Context, str, engine.NewEvaluator(eq, 10*time.Second), engine.ContextOpts{
				ScreenshotDir: ss,
			})
			if err != nil {
				return err
			}

			chromeOpts := []chromedp.ExecAllocatorOption{
				chromedp.Flag("headless", !cmd.Bool("mfa")),
			}

			if cdd := settings.Default.Get().GetChromeDataDirectory(); cdd != "" {
				chromeOpts = append(chromeOpts, chromedp.UserDataDir(cdd))
			}

			stpper := stepper.NewFactory()
			scriptCtx := acct.Run(c, stpper,
				provider.WithScriptOpts(
					script.WithScreenshotAfter(cmd.Bool("debug")),
					script.WithChromeOpts(chromeOpts...),
				),
			)
			select {
			case <-uiCtx.Done():
				return context.Cause(uiCtx)
			case <-scriptCtx.Done():
				return context.Cause(scriptCtx)
			}
		},
	}
	return cmd
}
