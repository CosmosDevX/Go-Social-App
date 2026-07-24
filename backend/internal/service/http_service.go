// Package service
package service

import (
	"net/http"
	"time"
)

type HTTPService struct {
	Server *http.Server
}

func NewHTTPService(handler http.Handler) HTTPService {
	return HTTPService{
		Server: &http.Server{
			Addr:              ":8080",
			Handler:           handler,
			ReadTimeout:       5 * time.Second,
			ReadHeaderTimeout: 2 * time.Second,
			WriteTimeout:      10 * time.Second,
			IdleTimeout:       120 * time.Second,
			MaxHeaderBytes:    1024 * 1024,
		},
	}
}
