package cli

import (
	"context"
	"errors"
	"github.com/tak-sh/tak/pkg/account"
	"github.com/tak-sh/tak/pkg/contexts"
	"github.com/tak-sh/tak/pkg/except"
	"github.com/tak-sh/tak/pkg/headless/engine"
	"github.com/tak-sh/tak/pkg/headless/script"
	"github.com/tak-sh/tak/pkg/renderer"
	"github.com/tak-sh/tak/pkg/settings"
	"github.com/tak-sh/tak/pkg/ui"
	"github.com/urfave/cli/v2"
	"log/slog"
	"os"
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

			acct, err := account.LoadFile(fp)
			if err != nil {
				return err
			}

			s, err := script.New(acct.GetSpec().GetLoginScript())
			if err != nil {
				return errors.Join(except.NewInvalid("failed to compile your login script"), err)
			}

			s.ScreenShotAfter = cmd.Bool("debug")

			str := renderer.NewStream()
			eq := engine.NewEventQueue()

			bubble := ui.NewBubbleUI(acct, str, eq)

			uiCtx, err := bubble.Start(cmd.Context, os.Stdin, os.Stdout)
			if err != nil {
				logger.Error("Failed to start the UI.", slog.String("err", err.Error()))
				return errors.Join(except.NewInternal("failed to start the UI"), err)
			}

			c, err := engine.NewContext(cmd.Context, str, engine.NewEvaluator(eq), engine.ContextOpts{
				ScreenshotDir: ss,
			})
			if err != nil {
				return err
			}

			scriptCtx, err := script.Run(c, s, script.WithHeadless(!cmd.Bool("mfa")))
			if err != nil {
				return err
			}

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
