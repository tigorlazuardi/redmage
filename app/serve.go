package app

import (
	"os"
	"strings"

	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/tigorlazuardi/redmage/app/routes"
)

func (rm *Redmage) Serve() error {
	// serves static files from the provided public dir (if exists)
	rm.App.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		r := routes.Routes{}
		e.Router.GET("/", r.Home)
		if rm.Public != nil {
			e.Router.GET("/*", apis.StaticDirectoryHandler(rm.Public, false))
		}
		return nil
	})

	// loosely check if it was executed using "go run"
	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())
	migratecmd.MustRegister(rm.App, rm.App.RootCmd, migratecmd.Config{
		// enable auto creation of migration files when making collection changes in the Admin UI
		// (the isGoRun check is to enable it only during development)
		Automigrate: isGoRun,
	})
	return rm.App.Start()
}
