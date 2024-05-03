package api

import (
	"context"

	"github.com/stephenafamo/bob"
	"github.com/tigorlazuardi/redmage/pkg/errs"
)

type executor func(exec bob.Executor) error

func (api *API) withTransaction(ctx context.Context, f executor) (err error) {
	tx, err := api.sqldb.BeginTx(ctx, nil)
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
