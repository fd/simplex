package http

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

type (
	Server struct {
		server   *http.Server
		listener net.Listener
		mux      *http.ServeMux
	}

	handler struct {
		*Server
	}
)

func New(addr string) *Server {
	s := &Server{
		server: &http.Server{
			Addr:         addr,
			ReadTimeout:  1 * time.Minute,
			WriteTimeout: 2 * time.Minute,
		},
		mux: http.NewServeMux(),
	}

	s.server.Handler = &handler{s}

	return s
}

func (s *Server) Start() error {
	addr, err := net.ResolveTCPAddr("tcp", s.server.Addr)
	if err != nil {
		return err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}

	s.listener = l
	go s.serve_with_listener()

	return nil
}

// TODO(fd) Environment should handle logging and errors
func (s *Server) serve_with_listener() {
	err := s.server.Serve(s.listener)
	if err != nil {
		panic(err)
	}
}

func (s *Server) Stop() error {
	return s.listener.Close()
}

func (s *Server) Handle(pattern string, handler http.Handler) {
	if len(pattern) == 0 {
		pattern = "/"
	}

	if pattern[len(pattern)-1] == '/' {
		prefix := pattern[:len(pattern)-1]
		handler = http.StripPrefix(prefix, handler)
	}

	s.mux.Handle(pattern, handler)
}

func (s *Server) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	s.Handle(pattern, http.HandlerFunc(handler))
}

func (s *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	tw := &http_TrackingResponseWriter{ResponseWriter: w, req: req}
	defer tw.LogHTTPTransaction()
	s.mux.ServeHTTP(tw, req)
}

type http_TrackingResponseWriter struct {
	http.ResponseWriter
	req    *http.Request
	status int
}

func (w *http_TrackingResponseWriter) LogHTTPTransaction() {
	if h := w.req.Host; h != "" && w.req.URL.Host == "" {
		w.req.URL.Host = h
	}
	if w.req.URL.Scheme == "" {
		w.req.URL.Scheme = "http"
	}
	fmt.Printf("%s %s %d\n", w.req.Method, w.req.URL, w.status)
}

func (w *http_TrackingResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}
