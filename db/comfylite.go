package db

import (
	"github.com/davidroman0O/comfylite3"
	"github.com/tigorlazuardi/redmage/config"
	"github.com/tigorlazuardi/redmage/pkg/errs"
)

func NewComfy(cfg *config.Config) (*comfylite3.ComfyDB, error) {
	target := cfg.String("db.string")
	db, err := comfylite3.Comfy(comfylite3.WithPath(target))
	if err != nil {
		return db, errs.Wrapf(err, "failed to create/open comfy db at %q", target)
	}
	return db, nil
}
