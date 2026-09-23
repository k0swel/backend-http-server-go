package webserver

import (
	"net/http"
	"web_server_go/db"
)

type WebServerConfig struct {
	Host          string   // "host:port"
	Timeout       Timeouts // таймауты
	HandleOptions bool     // true = OPTIONS запросы переходят в наш мидлваре, false - библиотека сама отвечат 204 на OPTIONS
}

type WebServerHandler struct {
	method string
	url    string
}

type WebServer_i interface {
	//CreateWebServer(config WebServerConfig) *webServer
	Stop()
	Handle(method string, path string, handler http.HandlerFunc)
	Run(database *db.SQLite3) error
}

type Timeouts struct {
	WriteTimeout      int // 2 сек поставить
	ReadHeaderTimeout int // 2 сек поставить
	IdleTimeout       int // 60 сек поставить. HTTP-keepalive
}
