package cli

import (
	"errors"
	"github.com/tak-sh/tak/generated/go/api/account/v1beta1"
	"github.com/tak-sh/tak/pkg/contexts"
	"github.com/tak-sh/tak/pkg/except"
	"github.com/tak-sh/tak/pkg/headless"
	"github.com/tak-sh/tak/pkg/headless/script"
	"github.com/tak-sh/tak/pkg/prompt"
	"github.com/tak-sh/tak/pkg/protoenc"
	"github.com/tak-sh/tak/pkg/settings"
	"github.com/tak-sh/tak/pkg/ui"
	"github.com/urfave/cli/v2"
	"log/slog"
	"os"
	"path/filepath"
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
		},
		Action: func(cmd *cli.Context) error {
			ss := cmd.String("screenshots")
			fp := cmd.Args().First()
			logger := contexts.GetLogger(cmd.Context)

			_, err := os.Stat(fp)
			if err != nil {
				return errors.Join(except.NewNotFound("failed to find account file %s", fp), err)
			}

			acct := new(v1beta1.Account)
			dir, name := filepath.Split(fp)
			err = protoenc.UnmarshalFile(acct, name, os.DirFS(dir))
			if err != nil {
				return errors.Join(except.NewInvalid("%s is not a valid account file", fp), err)
			}

			s, err := script.New(acct.GetSpec().GetLoginScript())
			if err != nil {
				return errors.Join(except.NewInvalid("failed to compile your login script"), err)
			}

			s.Debug = cmd.Bool("debug")

			str := prompt.NewStream()
			ct, err := headless.NewContext(cmd.Context, str, headless.WithScreenshotsDir(ss))
			if err != nil {
				return err
			}

			user := ui.NewReadWriterUI(os.Stdin, os.Stdout, str)

			_, err = user.Start(cmd.Context)
			if err != nil {
				logger.Error("Failed to start the UI.", slog.String("err", err.Error()))
				return errors.Join(except.NewInternal("failed to start the UI"), err)
			}

			err = script.Run(ct, s)
			if err != nil {
				return err
			}

			return nil
		},
	}
	return cmd
}
