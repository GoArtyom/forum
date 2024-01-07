package server

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"

	"forum/config"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(cfg *config.Config, handler http.Handler) error {
	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}

	s.httpServer = &http.Server{
		Addr:           cfg.Port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20, // 1 MB
		IdleTimeout:    time.Minute,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		TLSConfig:      tlsConfig,
	}
	return s.httpServer.ListenAndServeTLS(cfg.TLSCart, cfg.TLSKey)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
