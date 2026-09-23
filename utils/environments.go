package utils

// SQLITE3_TABLE - таблица SQLITE3, к которой будет подключение
// SQLITE3_FILEPATH - путь к файлу Базы Данных
import (
	"os"
	"strings"
)

var ENV []string = os.Environ()

func GetEnv() map[string]string {
	var environments map[string]string = make(map[string]string)
	for i := 0; i < len(ENV); i++ {
		var key string = strings.Split(ENV[i], "=")[0]
		var value string = strings.Split(ENV[i], "=")[1]
		environments[key] = value
	}
	return environments
}
