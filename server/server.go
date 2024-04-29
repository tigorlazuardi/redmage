package server

import (
	"context"
	"errors"
	"io/fs"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tigorlazuardi/redmage/api"
	"github.com/tigorlazuardi/redmage/config"
	"github.com/tigorlazuardi/redmage/pkg/caller"
	"github.com/tigorlazuardi/redmage/pkg/errs"
	"github.com/tigorlazuardi/redmage/pkg/log"
	"github.com/tigorlazuardi/redmage/server/routes"
)

type Server struct {
	server    *http.Server
	config    *config.Config
	PublicDir fs.FS
}

func (srv *Server) Start(exit <-chan struct{}) error {
	errch := make(chan error, 1)
	caller := caller.New(3)
	go func() {
		log.New(context.Background()).Caller(caller).Info(
			"starting http server", "address", "http://"+srv.server.Addr,
			"outbound_ip", "http://"+GetOutboundIP().String()+":"+srv.config.String("http.port"),
		)
		errch <- srv.server.ListenAndServe()
	}()

	select {
	case <-exit:
		log.New(context.Background()).Caller(caller).Info("received exit signal. shutting down server")
		ctx, cancel := context.WithTimeout(context.Background(), srv.config.Duration("http.shutdown_timeout"))
		defer cancel()
		return srv.server.Shutdown(ctx)
	case err := <-errch:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return errs.Wrap(err)
	}
}

func New(cfg *config.Config, api *api.API, publicDir fs.FS) *Server {
	router := chi.NewRouter()

	routes := routes.Routes{
		API:       api,
		Config:    cfg,
		PublicDir: publicDir,
	}

	routes.Register(router)

	server := &http.Server{
		Handler: router,
		Addr:    cfg.String("http.host") + ":" + cfg.String("http.port"),
	}

	return &Server{server: server, config: cfg}
}

// Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
