package cli

import (
	"context"
	"errors"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/chromedp/chromedp"
	"github.com/tak-sh/tak/generated/go/api/provider/v1beta1"
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
			&cli.GenericFlag{
				Name: "operation",
				Value: &EnumValue{
					Enum: []string{
						string(engine.OperationListAccounts),
						string(engine.OperationDownloadTransactions),
						string(engine.OperationLogin),
					},
				},
				Usage: "Select a specific operation to run for the script.",
			},
			&cli.StringFlag{
				Name:  "account",
				Usage: "Name of the account to target.",
				Action: func(c *cli.Context, s string) error {
					if s == "" && engine.Operation(c.Generic("operation").(*EnumValue).selected) == engine.OperationDownloadTransactions {
						return except.NewInvalid("the '--account' flag is required for the %s operation", string(engine.OperationDownloadTransactions))
					}
					return nil
				},
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

			str := renderer.NewStream()
			eq := engine.NewEventQueue()
			stpper := debug.NewFactory()

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

			prov, err := provider.New(c, mani, stpper, script.WithChromeOpts(chromeOpts...))
			if err != nil {
				return err
			}

			acctCtx, acctCancel := context.WithCancelCause(cmd.Context)
			op := engine.Operation(cmd.String("operation"))
			progMsg := fmt.Sprintf("Debugging %s...", op.ActionString())
			switch op {
			case engine.OperationLogin:
				go func() {
					var err error
					defer func() {
						acctCancel(err)
					}()
					err = prov.Login(acctCtx)
				}()
			case engine.OperationDownloadTransactions:
				go func() {
					var err error
					defer func() {
						acctCancel(err)
					}()
					a := cmd.String("account")

					err = prov.DownloadTransactions(acctCtx, a)
					if err != nil {
						err = errors.Join(fmt.Errorf("failed to downloading transactions for account %s", a), err)
						return
					}
				}()
			case engine.OperationListAccounts:
				go func() {
					var err error
					defer func() {
						acctCancel(err)
					}()
					var accts []*v1beta1.Account
					accts, err = prov.ListAccounts(acctCtx)
					fmt.Println(accts)
					return
				}()
			}

			scriptComp := ui.NewScriptComponent(progMsg, str, eq, logger)
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

			str := renderer.NewStream()
			eq := engine.NewEventQueue()
			stpper := stepper.NewFactory()

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

			prov, err := provider.New(c, acctRaw, stpper, script.WithChromeOpts(chromeOpts...), script.WithScreenshotAfter(cmd.Bool("debug")))
			if err != nil {
				return err
			}

			bubble := ui.NewScriptComponent(fmt.Sprintf("Adding your %s account...", acctRaw.GetMetadata().GetName()), str, eq, logger)
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

			acctCtx, acctCancel := context.WithCancelCause(cmd.Context)
			go func() {
				var err error
				defer func() {
					acctCancel(err)
				}()

				err = prov.Login(c)
				if err != nil {
					return
				}

				var acct []*v1beta1.Account
				acct, err = prov.ListAccounts(c)
				if err != nil {
					logger.Error("Failed to list accounts.", slog.String("err", err.Error()))
					return
				}

				for _, a := range acct {
					err = prov.DownloadTransactions(c, a.Name)
					if err != nil {
						logger.Error("Failed to download transactions.", slog.String("name", a.Name))
						return
					}
				}
			}()

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
