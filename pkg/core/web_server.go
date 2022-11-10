package core

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/vogtp/go-hcl"
)

// WebServer is a wrapper for a webserver
type WebServer struct {
	hcl      hcl.Logger
	port     int
	basepath string
	url      string
	srv      *http.Server
	mux      *http.ServeMux
	mid      []Middleware
	mu       sync.Mutex
	started  bool
}

func (w *WebServer) init(c *Core) {
	w.hcl = c.hcl.Named("web")

	if !IsFreePort(w.port) {
		w.hcl.Warnf("Port %v is not free, using a random port", w.port)
		w.port = 0
	}

	if w.port < 1 {
		p, err := GetFreePort()
		w.hcl.Warnf("Found free port to run on: %d", p)
		if err != nil {
			panic(err)
		}
		if p < 1 {
			panic(fmt.Errorf("no free port found"))
		}
		w.port = p
	}

	w.mux = http.NewServeMux()
	adr := fmt.Sprintf(":%d", w.port)
	w.srv = &http.Server{
		Addr:              adr,
		ReadTimeout:       1 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		Handler:           w.mux,
	}
	w.basepath = strings.TrimSuffix(w.basepath, "/")
	if len(w.basepath) < 2 {
		w.basepath = ""
	} else if !strings.HasPrefix(w.basepath, "/") {
		w.basepath = fmt.Sprintf("/%s", w.basepath)
	}
	ip := "localhost"
	i, err := GetOutboundIP()
	if err == nil {
		ip = i.String()
	}
	w.url = fmt.Sprintf("http://%s%s%s", ip, w.srv.Addr, w.basepath)
}

// Start the webserver
func (w *WebServer) Start() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.started {
		w.hcl.Info("Cannot start webserver, it is already running")
		return
	}
	w.started = true
	w.hcl.Warnf("Webserver listen on %s", w.url)
	go func() {
		err := w.srv.ListenAndServe()
		w.mu.Lock()
		defer w.mu.Unlock()
		w.started = false
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				w.hcl.Warnf("http server stopped: %v", err)
			} else {
				w.hcl.Errorf("http server error: %v", err)
			}
		}
	}()
}

// Stop the webserver
func (w *WebServer) Stop() {
	w.hcl.Warn("Stopping web server")
	w.srv.Shutdown(context.Background())
}

// BasePath is the url path
func (w *WebServer) BasePath() string {
	return w.basepath
}

// Middleware func
type Middleware func(http.HandlerFunc) http.HandlerFunc

// AddMiddleware adds middleware
func (w *WebServer) AddMiddleware(m Middleware) {
	if w.mid == nil {
		w.mid = make([]Middleware, 0)
	}
	w.mid = append(w.mid, m)
}

// Handle registers the handler for the given pattern.
// If a handler already exists for pattern, Handle panics.
func (w *WebServer) Handle(pattern string, handler http.Handler) {
	h := handler
	for _, m := range w.mid {
		h = m(h.ServeHTTP)
	}
	w.mux.Handle(w.basepath+pattern, h)
}

// HandleFunc registers the handler function for the given pattern.
func (w *WebServer) HandleFunc(pattern string, handler http.HandlerFunc) {
	w.Handle(pattern, http.HandlerFunc(handler))
}

// URL returns the url the server is listening on
func (w *WebServer) URL() string {
	return w.url
}
