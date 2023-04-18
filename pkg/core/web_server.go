package core

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/vogtp/som/pkg/core/log"
	"golang.org/x/exp/slog"
)

// WebServer is a wrapper for a webserver
type WebServer struct {
	log      *slog.Logger
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
	w.log = c.log.With(log.Component, "web")

	if !IsFreePort(w.port) {
		w.log.Warn("Port is not free, using a random port", "port", w.port)
		w.port = 0
	}

	if w.port < 1 {
		p, err := GetFreePort()
		w.log.Warn("Found free port", "port", p)
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
		w.log.Info("Cannot start webserver, it is already running")
		return
	}
	w.started = true
	w.log = w.log.With("listen_adr", w.url)
	w.log.Warn("Webserver starting")
	go func() {
		err := w.srv.ListenAndServe()
		w.mu.Lock()
		defer w.mu.Unlock()
		w.started = false
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				w.log.Warn("http server stopped", log.Error, err)
			} else {
				w.log.Error("http server error", log.Error, err)
			}
		}
	}()
}

// Stop the webserver
func (w *WebServer) Stop() {
	w.log.Warn("Stopping web server")
	if err := w.srv.Shutdown(context.Background()); err != nil {
		w.log.Warn("cannot shutdown webserver", log.Error, err)
	}
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
