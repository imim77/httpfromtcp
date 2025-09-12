package server

import (
	"bytes"
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

type Handler func(w io.Writer, req *request.Request) *HandlerError

func runConnection(s *Server, conn io.ReadWriteCloser) {
	defer conn.Close()
	headers := response.GetDefaultHeaders(0)
	r, err := request.RequestFromReader(conn)
	if err != nil {
		response.WriteStatusLine(conn, response.StatusBadRequest)
		response.WriteHeaders(conn, headers)
		return
	}

	writer := bytes.NewBuffer([]byte{})
	handlerError := s.Handler(writer, r)

	var body []byte = nil
	var status response.StatusCode = response.StatusOK
	if handlerError != nil {
		status = handlerError.StatusCode
		body = []byte(handlerError.Message)
	} else {
		body = writer.Bytes()
	}

	headers.Replace("Content-Length", fmt.Sprintf("%d", len(body)))

	response.WriteStatusLine(conn, status)
	response.WriteHeaders(conn, headers)
	conn.Write(body)

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
