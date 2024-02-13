package app

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/tigorlazuardi/redmage/app/config"
	"github.com/tigorlazuardi/redmage/app/routes"
)

func getOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return ""
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

func (rm *Redmage) Serve() error {
	// serves static files from the provided public dir (if exists)
	rm.App.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		e.Router.Use(middleware.Gzip())
		e.Server.Addr = fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
		r := routes.Routes{
			Config: cfg,
		}
		if r.Config.HotReload {
			e.Router.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
				return func(c echo.Context) error {
					c.Response().Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
					return next(c)
				}
			})
			r.HotReload = make(chan struct{}, 1<<4)
			r.HotReload <- struct{}{}
		}
		r.Register(e.Router)
		e.Router.GET("/*", apis.StaticDirectoryHandler(rm.Public, false))
		slog.Info(fmt.Sprintf("Server outbound host: http://%s:%d", getOutboundIP(), cfg.Server.Port))
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
