package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"
)

type Server struct {
	srv *http.Server
}

func New(handler http.Handler, port string) *Server {
	address := net.JoinHostPort("0.0.0.0", port)

	srv := &http.Server{
		Addr:    address,
		Handler: handler,
	}

	return &Server{
		srv: srv,
	}
}

func (s *Server) Run(ctx context.Context) error {
	go func() {
		err := s.srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			err = fmt.Errorf("Listen and serve error: %w", err)
		}
	}()

	<-ctx.Done()

	downCtx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := s.srv.Shutdown(downCtx); err != nil {
		return fmt.Errorf("shotdown: %w", err)
	}

	return nil
}
