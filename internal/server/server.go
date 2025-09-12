package server

import (
	"fmt"
	"io"
	"net"

	"github.com/imim77/httpfromtcp/internal/request"
	"github.com/imim77/httpfromtcp/internal/response"
)

type serverState string

const (
	StateRunning serverState = "running"
	StateClosed  serverState = "stopped"
)

type Server struct {
	State   serverState
	Handler Handler
}

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

type Handler func(w *response.Writer, req *request.Request)

func runConnection(s *Server, conn io.ReadWriteCloser) {
	defer conn.Close()
	responseWriter := response.NewWriter(conn)
	r, err := request.RequestFromReader(conn)
	if err != nil {
		responseWriter.WriteStatusLine(response.StatusBadRequest)
		responseWriter.WriteHeaders(*response.GetDefaultHeaders(0))
		return
	}
	s.Handler(responseWriter, r)
}

func runServer(s *Server, listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if s.State == StateClosed {
			return
		}
		if err != nil {
			return
		}
		go runConnection(s, conn)

	}
}

func Serve(port uint16, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	s := &Server{State: StateRunning, Handler: handler}
	go runServer(s, listener)
	return s, nil
}

func (s *Server) Close() error {
	s.State = StateClosed
	return nil
}
