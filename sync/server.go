package sync

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

//Server represents a server
type Server struct {
	*http.Server
	termination chan bool
}

func (s *Server) shutdown() {
	<-s.termination
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := s.Server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}

//StopOnSiginals stops server on siginal
func (s *Server) StopOnSiginals(siginals ...os.Signal) {
	notification := make(chan os.Signal, 1)
	signal.Notify(notification, siginals...)
	<-notification
	s.Stop()
}

//Stop stop server
func (s *Server) Stop() {
	s.termination <- true
}

//NewServer creates a new server
func NewServer(service Service, port int) *Server {
	router := NewRouter(service)
	return &Server{
		termination: make(chan bool, 1),
		Server: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: router,
		},
	}
}
