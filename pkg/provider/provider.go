package provider

import (
	"context"
	"errors"
	"github.com/eddieowens/opts"
	"github.com/tak-sh/tak/generated/go/api/provider/v1beta1"
	"github.com/tak-sh/tak/pkg/contexts"
	"github.com/tak-sh/tak/pkg/except"
	"github.com/tak-sh/tak/pkg/headless/engine"
	"github.com/tak-sh/tak/pkg/headless/script"
	"github.com/tak-sh/tak/pkg/headless/step/stepper"
	"github.com/tak-sh/tak/pkg/protoenc"
	"github.com/tak-sh/tak/pkg/utils/grpcutils"
	"github.com/tak-sh/tak/pkg/validate"
	"log/slog"
	"os"
	"path/filepath"
)

func New(prov *v1beta1.Manifest) (a *Manifest, err error) {
	a = &Manifest{
		Manifest: prov,
	}

	a.Login, err = script.New(a.GetSpec().GetLogin().GetScript())
	if err != nil {
		return nil, errors.Join(errors.New("login script"), err)
	}

	a.DownloadTransactions, err = script.New(a.GetSpec().GetDownloadTransactions().GetScript())
	if err != nil {
		return nil, errors.Join(errors.New("download transactions script"), err)
	}

	return a, nil
}

type RunOpts struct {
	ScriptOpts               []opts.Opt[script.RunOpts]
	SkipLogin                bool
	SkipDownloadTransactions bool
}

func (r RunOpts) DefaultOptions() RunOpts {
	return RunOpts{}
}

func WithSkipLogin(b bool) opts.Opt[RunOpts] {
	return func(r *RunOpts) {
		r.SkipLogin = b
	}
}

func WithSkipDownloadTransactions(b bool) opts.Opt[RunOpts] {
	return func(r *RunOpts) {
		r.SkipDownloadTransactions = b
	}
}

func WithScriptOpts(o ...opts.Opt[script.RunOpts]) opts.Opt[RunOpts] {
	return func(r *RunOpts) {
		r.ScriptOpts = append(r.ScriptOpts, o...)
	}
}

var _ grpcutils.ProtoWrapper[*v1beta1.Manifest] = &Manifest{}
var _ validate.Validator = &Manifest{}

type Manifest struct {
	*v1beta1.Manifest
	Login                *script.Script
	DownloadTransactions *script.Script
}

func (p *Manifest) Run(c *engine.Context, stepperFact stepper.Factory, o ...opts.Opt[RunOpts]) context.Context {
	ctx, cancel := context.WithCancelCause(c.Context)
	op := opts.DefaultApply(o...)

	go func() {
		var err error
		defer func() {
			cancel(err)
		}()
		logger := contexts.GetLogger(c.Context)
		if !op.SkipLogin {
			stper := stepperFact.NewStepper(p.Login.Signals, p.Login.Steps)
			err = script.Run(c, p.Login, stper, op.ScriptOpts...)
			if err != nil {
				logger.Error("Failed to run login script.", slog.String("err", err.Error()))
				return
			}
		}

		if !op.SkipDownloadTransactions {
			stper := stepperFact.NewStepper(p.DownloadTransactions.Signals, p.DownloadTransactions.Steps)
			err = script.Run(c, p.DownloadTransactions, stper, op.ScriptOpts...)
			if err != nil {
				logger.Error("Failed to run download transactions script.", slog.String("err", err.Error()))
				return
			}
		}
	}()

	return ctx
}

func (p *Manifest) Validate() error {
	err := p.Login.Validate()
	if err != nil {
		return errors.Join(errors.New("login script"), err)
	}

	err = p.DownloadTransactions.Validate()
	if err != nil {
		return errors.Join(errors.New("download transactions script"), err)
	}

	return nil
}

func (p *Manifest) ToProto() *v1beta1.Manifest {
	return p.Manifest
}

func LoadFile(fp string) (*v1beta1.Manifest, error) {
	_, err := os.Stat(fp)
	if err != nil {
		return nil, errors.Join(except.NewNotFound("failed to find account file %s", fp), err)
	}

	acct := new(v1beta1.Manifest)
	dir, name := filepath.Split(fp)
	err = protoenc.UnmarshalFile(acct, name, os.DirFS(dir))
	if err != nil {
		return nil, errors.Join(except.NewInvalid("%s is not a valid account file", fp), err)
	}

	return acct, nil
}
