package http

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
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
	tw := http_NewTrackingResponseWriter(w, req)

	defer tw.LogHTTPTransaction()
	s.mux.ServeHTTP(tw, req)
}

type http_TrackingResponseWriter struct {
	http.ResponseWriter
	req    *http.Request
	url    string
	start  time.Time
	status int
}

func http_NewTrackingResponseWriter(w http.ResponseWriter, req *http.Request) *http_TrackingResponseWriter {
	var u url.URL
	u = *req.URL

	if h := req.Host; h != "" && u.Host == "" {
		u.Host = h
	}
	if u.Scheme == "" {
		u.Scheme = "http"
	}

	return &http_TrackingResponseWriter{
		ResponseWriter: w,
		req:            req,
		url:            u.String(),
		start:          time.Now(),
		status:         200,
	}
}

func (w *http_TrackingResponseWriter) LogHTTPTransaction() {
	duration := time.Now().Sub(w.start)
	fmt.Printf("[method: %s, status: %d, duration: %s] %s\n", w.req.Method, w.status, duration, w.url)
}

func (w *http_TrackingResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}
