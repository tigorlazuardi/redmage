package app

import (
	"io/fs"

	"github.com/pocketbase/pocketbase"

	_ "github.com/tigorlazuardi/redmage/migrations"
)

type Redmage struct {
	App    *pocketbase.PocketBase
	Public fs.FS
}

func Start(public fs.FS) error {
	app := pocketbase.New()
	rm := Redmage{
		App:    app,
		Public: public,
	}
	return rm.Serve()
}
