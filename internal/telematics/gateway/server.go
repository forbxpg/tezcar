// SPDX-License-Identifier: BUSL-1.1

// Package gateway implements a Telematics Gateway server.
package gateway

import (
	"errors"
	"log/slog"
	"net"
)

// Server is a Telematics Gateway server.
// It is responsible for accepting incoming connections and
// handling them.
// It uses a logger to log messages.
type Server struct {
	Logger *slog.Logger
}

// Serve listens for incoming connections and handles them.
// It returns an error if the server fails to start.
func (s *Server) Serve(ln net.Listener) error {
	for {
		conn, err := ln.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return nil
			}
			return err
		}
		go s.handleConn(conn)
	}
}

// handleConn handles an incoming connection.
func (s *Server) handleConn(conn net.Conn) {
	defer func() {
		_ = conn.Close()
		s.Logger.Debug("client disconnected", "remote", conn.RemoteAddr())
	}()
	s.Logger.Info("client connected", "remote", conn.RemoteAddr())
}
