package provider

import (
	"context"
	"errors"
	"github.com/eddieowens/opts"
	"github.com/tak-sh/tak/generated/go/api/provider/v1beta1"
	"github.com/tak-sh/tak/pkg/except"
	"github.com/tak-sh/tak/pkg/headless/engine"
	"github.com/tak-sh/tak/pkg/headless/script"
	"github.com/tak-sh/tak/pkg/headless/step/stepper"
	"github.com/tak-sh/tak/pkg/protoenc"
	"github.com/tak-sh/tak/pkg/utils/grpcutils"
	"github.com/tak-sh/tak/pkg/validate"
	"os"
	"path/filepath"
)

func New(c *engine.Context, prov *v1beta1.Manifest, s stepper.Factory, o ...opts.Opt[script.RunOpts]) (a *Manifest, err error) {
	a = &Manifest{
		Manifest:       prov,
		Ctx:            c,
		Opts:           o,
		StepperFactory: s,
	}

	a.LoginScript, err = script.New(a.GetSpec().GetLogin().GetScript())
	if err != nil {
		return nil, errors.Join(errors.New("login field"), err)
	}

	a.DownloadTransactionsScript, err = script.New(a.GetSpec().GetDownloadTransactions().GetScript())
	if err != nil {
		return nil, errors.Join(errors.New("download_transactions field"), err)
	}

	a.ListAccountsScript, err = script.New(a.GetSpec().GetListAccounts().GetScript())
	if err != nil {
		return nil, errors.Join(errors.New("list_accounts field"), err)
	}

	a.listAccountMapping, err = newAccountMapping(a.GetSpec().GetListAccounts().GetOutputs().GetAccount())
	if err != nil {
		return nil, errors.Join(errors.New(`"account" field for "list_accounts" spec`), err)
	}

	return a, nil
}

var _ grpcutils.ProtoWrapper[*v1beta1.Manifest] = &Manifest{}
var _ validate.Validator = &Manifest{}
var _ Provider = &Manifest{}

type Manifest struct {
	*v1beta1.Manifest
	LoginScript                *script.Script
	DownloadTransactionsScript *script.Script
	ListAccountsScript         *script.Script

	StepperFactory stepper.Factory
	Opts           []opts.Opt[script.RunOpts]
	Ctx            *engine.Context

	listAccountMapping *AccountMapping
}

func (p *Manifest) ListAccounts(ctx context.Context) ([]*v1beta1.Account, error) {
	c := p.Ctx.WithContext(ctx)
	stper := p.StepperFactory.NewStepper(p.ListAccountsScript.Signals, p.ListAccountsScript.Steps)

	p.Ctx.Evaluator.EventQueue() <- &engine.ChangeOperationEvent{
		To: engine.OperationListAccounts,
	}

	err := script.Run(c, p.ListAccountsScript, stper, p.Opts...)
	if err != nil {
		return nil, err
	}

	accts := make([]*v1beta1.Account, 0)
	c.TemplateData.ForEach(p.GetSpec().GetListAccounts().GetOutputs().GetForEach(), func(r *engine.TemplateData) {
		accts = append(accts, p.listAccountMapping.Render(r))
	})

	return accts, nil
}

func (p *Manifest) Login(ctx context.Context) error {
	c := p.Ctx.WithContext(ctx)
	stper := p.StepperFactory.NewStepper(p.LoginScript.Signals, p.LoginScript.Steps)

	p.Ctx.Evaluator.EventQueue() <- &engine.ChangeOperationEvent{
		To: engine.OperationLogin,
	}

	err := script.Run(c, p.LoginScript, stper, p.Opts...)
	if err != nil {
		return err
	}

	return nil
}

func (p *Manifest) DownloadTransactions(ctx context.Context, acctName string) error {
	c := p.Ctx.WithContext(ctx)
	stper := p.StepperFactory.NewStepper(p.DownloadTransactionsScript.Signals, p.DownloadTransactionsScript.Steps)

	p.Ctx.Evaluator.EventQueue() <- &engine.ChangeOperationEvent{
		To:      engine.OperationDownloadTransactions,
		Message: "for account " + acctName,
	}

	err := script.Run(c, p.DownloadTransactionsScript, stper, p.Opts...)
	if err != nil {
		return err
	}

	return nil
}

func (p *Manifest) Validate() error {
	err := p.LoginScript.Validate()
	if err != nil {
		return errors.Join(errors.New("login script"), err)
	}

	err = p.DownloadTransactionsScript.Validate()
	if err != nil {
		return errors.Join(errors.New("download transactions script"), err)
	}

	err = p.ListAccountsScript.Validate()
	if err != nil {
		return errors.Join(errors.New("list accounts script"), err)
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
