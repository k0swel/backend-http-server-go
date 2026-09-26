package webserver

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"
	"web_server_go/db"
	"web_server_go/utils"

	"github.com/bradfitz/gomemcache/memcache"
)

const serverHeader = "k0swel-server-go"

type WebServer_s struct {
	server          *http.Server
	mux             *http.ServeMux
	Db              *db.SQLite3
	ctx             context.Context
	Memcache_client *memcache.Client
}

func CreateWebServer(config WebServerConfig) *WebServer_s {
	timeouts := config.Timeout
	server := http.Server{
		Addr:                         config.Host,
		DisableGeneralOptionsHandler: config.HandleOptions,
		ReadTimeout:                  time.Duration(timeouts.ReadHeaderTimeout) * time.Second,
		WriteTimeout:                 time.Duration(timeouts.WriteTimeout) * time.Second,
		IdleTimeout:                  time.Duration(timeouts.IdleTimeout) * time.Second,
		ErrorLog:                     utils.Logger,
	}

	s := &WebServer_s{
		server: &server,
		mux:    http.NewServeMux(),
	}
	server.Handler = middleware(s.mux)
	return s
}

func (s *WebServer_s) Run() error {
	s.Handle("POST", "/api/v1/register", s.RegisterHandler)
	s.Handle("POST", "/api/v1/auth", s.LoginHandler)
	s.Handle("GET", "/api/v1/account", s.GetAccountHandler)
	s.Handle("GET", "/api/v1/hostname", s.GetBackendHostname)
	return s.server.ListenAndServe()
}

func (s *WebServer_s) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		utils.Logger.Printf("graceful shutdown failed: %v, closing hard", err)
		if err := s.server.Close(); err != nil {
			utils.Logger.Printf("hard close failed: %v, exiting", err)
			os.Exit(7)
		}
	}
}

/* ---------- Middleware ---------- */

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setCORSHeaders(w) // Уcтановка стандартных CORS заголовков (не PREFLIGHT запрос)
		w.Header().Set("Server", serverHeader)

		if r.Method == http.MethodOptions {
			setCORSHeadersOptionsPreflight(w)
			w.WriteHeader(http.StatusNoContent)
			return
		}

		time_now := time.Now() // Время обработки запроса
		next.ServeHTTP(w, r)
		utils.Logger.Printf("Запрос [%s (%s). Client-IP: %s] завершился за %v \n", r.Method, r.RequestURI, r.Header.Get("X-Real-IP"), time.Since(time_now)) // Время обработки запроса
	})
}

// Установка CORS заголовков в MiddleWARE
func setCORSHeadersOptionsPreflight(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", fmt.Sprintf("https://%s", utils.GetEnv()["DOMAIN_NAME"]))
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Max-Age", "86400")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}

// Установка CORS заголовков
func setCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
}

// Handle регистрирует хендлер на путь с проверкой метода.
func (s *WebServer_s) Handle(method string, path string, h http.HandlerFunc) {
	s.mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h(w, r)
	})
}
