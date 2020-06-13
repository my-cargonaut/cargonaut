package http

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"golang.org/x/net/netutil"
)

// Server implements a http server which handles requests using the configured
// Router.
type Server struct {
	lis net.Listener
	srv *http.Server

	doneCh chan struct{}
	errCh  chan error
}

// NewServer creates a new http server listening on the configured address and
// using the configured http.Handler to serve requests. An error is returned if
// the address the server should listen on is already in use.
func NewServer(log *log.Logger, addr string, handler http.Handler) (*Server, error) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("create listener: %w", err)
	}
	lis = netutil.LimitListener(lis, 512)

	srv := &http.Server{
		Addr:         addr,
		Handler:      h2c.NewHandler(handler, &http2.Server{}),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
		ErrorLog:     log,
	}

	return &Server{
		lis: lis,
		srv: srv,

		doneCh: make(chan struct{}),
		errCh:  make(chan error),
	}, nil
}

// ListenAddr returns the address the http server is listening on.
func (s *Server) ListenAddr() net.Addr {
	return s.lis.Addr()
}

// ListenError returns the receive-only channel which signals errors during http
// server startup.
func (s *Server) ListenError() <-chan error {
	return s.errCh
}

// Run starts the http server in a separate goroutine and pushes errors into the
// errors channel. This method is non-blocking. Use the ListenError method to
// listen for errors which occure during startup.
func (s *Server) Run() {
	go func() {
		if err := s.srv.Serve(s.lis); err != nil && err != http.ErrServerClosed {
			s.errCh <- fmt.Errorf("serve http: %w", err)
		}
		close(s.errCh)
		close(s.doneCh)
	}()
}

// Shutdown stops the http server gracefully.
func (s *Server) Shutdown(ctx context.Context) error {
	// Shutdown the http server.
	if err := s.srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown http server: %w", err)
	}

	// Wait for the context to expire or a graceful shutdown signaled by the
	// done channel.
	select {
	case <-s.doneCh:
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil
}
