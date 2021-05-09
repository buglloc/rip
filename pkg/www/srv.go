package www

import (
	"context"
	"embed"
	"fmt"
	"net/http"

	log "github.com/buglloc/simplelog"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/websocket"

	"github.com/buglloc/rip/v2/pkg/cfg"
	"github.com/buglloc/rip/v2/pkg/hub"
)

//go:embed static/*
var staticFiles embed.FS
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type HttpSrv struct {
	https      http.Server
	tokens     *tokenManager
	inShutdown bool
}

func NewHttpSrv() *HttpSrv {
	srv := &HttpSrv{
		tokens: NewTokenManager(),
	}

	srv.https = http.Server{
		Addr:    cfg.HttpAddr,
		Handler: srv.router(),
	}

	return srv
}

func (s *HttpSrv) Addr() string {
	return s.https.Addr
}

func (s *HttpSrv) ListenAndServe() error {
	err := s.https.ListenAndServe()
	if s.inShutdown {
		return nil
	}

	return err
}

func (s *HttpSrv) Shutdown(ctx context.Context) error {
	s.inShutdown = true
	return s.https.Shutdown(ctx)
}

func (s *HttpSrv) router() http.Handler {
	var staticFS = http.FS(staticFiles)
	fs := http.FileServer(staticFS)

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	r.Get("/", serveStatic(staticFS, "/static/index.html"))
	r.Get("/session", serveStatic(staticFS, "/static/session.html"))
	r.Get("/ping", serveStatic(staticFS, "/static/pong.txt"))

	r.Get("/start", func(w http.ResponseWriter, r *http.Request) {
		token, err := s.tokens.NewToken()
		if err != nil {
			log.Error("can't create new token", "err", err)
			http.Error(w, "can't start new session", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/session?token=%s", token), http.StatusTemporaryRedirect)
	})

	r.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if token == "" {
			http.NotFound(w, r)
			return
		}

		channelID, err := s.tokens.ParseToken(token)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, fmt.Sprintf("can't upgrade request: %v", err), http.StatusInternalServerError)
			return
		}

		hub.Register(ws, channelID)
	})

	r.Mount("/static/", fs)
	return r
}
