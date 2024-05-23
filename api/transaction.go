package api

import (
	"context"
	"database/sql"

	"github.com/stephenafamo/bob"
	"github.com/tigorlazuardi/redmage/pkg/errs"
)

type txAble interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type executor func(exec bob.Executor) error

func (api *API) withTransaction(ctx context.Context, f executor) (err error) {
	tx, err := api.txAble.BeginTx(ctx, nil)
	if err != nil {
		return errs.Wrapw(err, "failed to begin transaction")
	}

	exec := bob.New(tx)
	err = f(exec)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
