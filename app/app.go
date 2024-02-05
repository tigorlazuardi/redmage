package app

import (
	"io/fs"

	"github.com/pocketbase/pocketbase"

	_ "github.com/tigorlazuardi/redmage/migrations"
)

type Redmage struct {
	App           *pocketbase.PocketBase
	HTMLTemplates fs.FS
}

func Start(embeds fs.FS) error {
	app := pocketbase.New()
	rm := Redmage{
		App:           app,
		HTMLTemplates: embeds,
	}
	return rm.Serve()
}
