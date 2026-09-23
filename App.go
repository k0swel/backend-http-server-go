package main

import (
	_ "web_server_go/users"
	utils "web_server_go/utils"
)

var env map[string]string = utils.GetEnv() // Переменные окружения,переданные процессу
// env["SQLITE3_FILE"] - файл СУБД sqlite3
// env["SQLITE3_TABLE"] - единственная таблица нашего веб-приложения - таблица Users
// env["MEMCACHED_ADDR"] - адрес сервиса MEMCACHED "127.0.0.1:11211"

func Run() {
	var sqlite3file string = env["SQLITE3_FILE"]
	var sqlite3table string = env["SQLITE3_TABLE"]
	var memcached_addr string = env["MEMCACHED_ADDR"]
	var services Service
	resultMemcached := services.InitializeMemcached(memcached_addr)
	sql := services.InitializeSQL(sqlite3file, sqlite3table)
	services.CreateWebServer(resultMemcached, sql, 2, 2, 60)
	if sql == nil {
		utils.Logger.Printf("Возникла ошибка при инициализации базы данных\n")
	}
	utils.Logger.Printf("Запуск веб-сервера\n")
	err := services.Webserver.Run()
	if err != nil {
		utils.Logger.Printf("Возникла ошибка при запуксе веб-сервера %v\n", err)
	}

}
