package app

import (
	"context"
	"flowcargo/internal/shared/logger"
	"net/http"
)

type Server struct {
	srv *http.Server // Define your server fields here
}

func wireServer() Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})
	return Server{
		srv: &http.Server{
			Addr:    ":8080",
			Handler: mux,
		},
	}
}

func (s Server) getAddress() string {
	return s.srv.Addr
}

func (s Server) start(l logger.Logger) error {
	l.Infof("Starting server on address %s", s.getAddress())
	return s.srv.ListenAndServe()
}

func (s Server) shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
