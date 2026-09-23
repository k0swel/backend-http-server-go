package main

import (
	"context"
	"web_server_go/db"
	"web_server_go/utils"
	webserver "web_server_go/web_server"

	memcached "github.com/bradfitz/gomemcache/memcache"
)

type Service struct {
	Database  *db.SQLite3           // интерфейс для взаимодействия с СУБД SQLITE3
	Webserver webserver.WebServer_s // интерфейс для взаимодействия с HTTP
	Memcache  *memcached.Client     // Интерфейс для взаимодействия с Memcached
}

// Инициаализация Service.Webserver
func (s *Service) CreateWebServer(memcache_client *memcached.Client, db *db.SQLite3, writeTimeout int, readHeaderTimeout int, idleTimeout int) {
	// some body...
	config := webserver.WebServerConfig{Host: "127.0.0.1:8081", Timeout: webserver.Timeouts{WriteTimeout: writeTimeout, ReadHeaderTimeout: readHeaderTimeout, IdleTimeout: idleTimeout}, HandleOptions: true}
	s.Webserver = *webserver.CreateWebServer(config)
	s.Webserver.Memcache_client = memcache_client
	s.Webserver.Db = db
}

// Инициализация Service.database. ОБЯЗАТЕЛЬНЫЕ параметры:
// filename - путь к файлу sqlite3
// table - таблица, которую необходимо открыть.
// В случае невозможности зайти в БД, будет os.Exit(2)
func (s *Service) InitializeSQL(filename string, table string) *db.SQLite3 {
	ctx := context.Background()
	db, err := db.NewSQLite3(ctx, filename, table)
	if err != nil {
		utils.Logger.Printf("Ошибка при инициализации базы данных %v...\n", err)
		return nil
	} else {
		utils.Logger.Printf("Успешное подключение к %s(%s)", filename, table)
	}
	s.Database = db
	return db
}

func (s *Service) InitializeMemcached(server ...string) *memcached.Client {
	memcacheClient := memcached.New(server...)
	err := memcacheClient.Ping()
	if err != nil {
		utils.Logger.Printf("Error during test connection to memcached servers (%v).", err)
	}
	s.Memcache = memcacheClient
	return s.Memcache
}
