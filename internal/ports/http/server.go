package http

import (
	"fmt"
	"net/http"
)

func NewServer(port int) *http.Server {
	return &http.Server{
		Addr: fmt.Sprintf(":%d", port),
	}
}
