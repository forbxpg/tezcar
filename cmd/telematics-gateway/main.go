// SPDX-License-Identifier: BUSL-1.1

// Package main implements the Telematics Gateway server.
package main

import (
	"context"
	"flag"
	"log/slog"
	"net"
	"os"

	"github.com/forbxpg/tezcar/internal/telematics/gateway"
)

// main is the entry point for the Telematics Gateway server.
// It creates a logger and a listener, and then starts the server.
func main() {
	addr := flag.String("addr", ":5027", "address to listen on")
	flag.Parse()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	var lc net.ListenConfig
	ln, err := lc.Listen(context.Background(), "tcp", *addr)
	if err != nil {
		logger.Error("listen", "err", err)
		os.Exit(1)
	}
	logger.Info("listening", "addr", ln.Addr())
	srv := &gateway.Server{Logger: logger}
	if err := srv.Serve(ln); err != nil {
		logger.Error("serve", "err", err)
		os.Exit(1)
	}
}
