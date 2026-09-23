package webserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"web_server_go/users"
	"web_server_go/utils"
	webserver_utils "web_server_go/web_server/utils"

	"github.com/bradfitz/gomemcache/memcache"
)

/* ---------- Регистрация ---------- */

func (s *WebServer_s) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var payload users.User
	err := webserver_utils.BodyToJSON(r.Body, &payload)
	utils.Logger.Printf("%s\n", payload.Password)
	if err != nil {
		utils.Logger.Printf("Не удалось конвертировать Body запроса в JSON")
	}
	utils.Logger.Printf("register attempt: email=%s", payload.Email)

	_, err = s.Db.GetUserByEmail(ctx, payload.Email)
	switch {
	case err == nil:
		webserver_utils.WriteError(w, http.StatusConflict, "email already exists")
		return
	case errors.Is(err, users.ErrUserNotFoundDB):
		// ок, регистрируем
	default:
		utils.Logger.Printf("GetUserByEmail: %v", err)
		webserver_utils.WriteError(w, http.StatusInternalServerError, "database error")
		return
	}
	if err := s.Db.CreateUser(ctx, payload); err != nil {
		utils.Logger.Printf("CreateUser: %v", err)
		webserver_utils.WriteError(w, http.StatusInternalServerError, "cannot create user")
		return
	} else {
		utils.Logger.Printf("Пользователь %s успешно добавлен в базу.\n", payload.Email)
	}
	cookie := http.Cookie{Name: "session-id", Value: webserver_utils.GenerateSessionCookie(), Path: "/", MaxAge: int(24 * 3600), HttpOnly: true, SameSite: http.SameSiteLaxMode}
	http.SetCookie(w, &cookie)
	s.Memcache_client.Set(&memcache.Item{Key: cookie.Value, Value: []byte(payload.Email), Expiration: 24 * 3600})
	data_from_memcache, _ := s.Memcache_client.Get(cookie.Value)
	fmt.Printf("%v\n", data_from_memcache)
	webserver_utils.WriteJSON(w, http.StatusCreated, map[string]string{
		"status":  "ok",
		"message": "registration is successful",
	})
}

/* Обработчик события авторизации*/
func (s *WebServer_s) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var payload users.User
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		webserver_utils.WriteError(w, http.StatusBadRequest, "invalid Payload: ")
		utils.Logger.Printf("Ошибка. Пользователь не отправил никакой Payload")
		return
	}
	user, err := s.Db.GetUserByEmail(r.Context(), payload.Email)
	switch {
	case errors.Is(err, users.ErrUserNotFoundDB):
		webserver_utils.WriteError(w, http.StatusUnauthorized, "authorization denied")
		utils.Logger.Printf("Пользователь %s не найден в базе\n", payload.Email)
		return
	case err != nil:
		utils.Logger.Printf("GetUserByEmail: %v", err)
		webserver_utils.WriteError(w, http.StatusInternalServerError, "database error")
		return
	}

	passwordHash, _ := s.Db.GetPasswordHash(r.Context(), user.Email)
	user.PasswordHash = string(passwordHash)
	if !utils.VerifyPassword([]byte(payload.Password), []byte(user.PasswordHash)) {
		webserver_utils.WriteError(w, http.StatusUnauthorized, "authorization denied")
		utils.Logger.Printf("Пользователь %s неверно ввёл пароль.\n", payload.Email)
		return
	}

	// --- Установка куки и сохранение сессии в memcache ---
	cookie := http.Cookie{
		Name:     "session-id",
		Value:    webserver_utils.GenerateSessionCookie(),
		Path:     "/",
		MaxAge:   int(24 * 3600),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookie)

	if err := s.Memcache_client.Set(&memcache.Item{
		Key:        cookie.Value,
		Value:      []byte(user.Email),
		Expiration: 24 * 3600,
	}); err != nil {
		utils.Logger.Printf("Memcache Set: %v", err)
	}

	webserver_utils.WriteJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"message": "authorization is successful",
	})
}

/*	---------Получение информации об аккаунте--------------- */
func (s *WebServer_s) GetAccountHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session-id")
	if err == http.ErrNoCookie {
		utils.Logger.Printf("Клиент (%s) не направил куки. Будет возвращён 301 /\n", r.Header.Get("X-Real-IP"))
		http.Redirect(w, r, fmt.Sprintf("https://%s", utils.GetEnv()["DOMAIN_NAME"]), 301)
		return
	} else if err == nil {
		email_from_cookie, err := s.Memcache_client.Get(cookie.Value)
		if err == memcache.ErrCacheMiss {
			utils.Logger.Printf("Сеанс пользователя (%s) истёк. Перенаправляю его на страницу авторизации\n", r.Header.Get("X-Real-IP"))
			return
		}
		user, err := s.Db.GetUserByEmail(r.Context(), string(email_from_cookie.Value))
		if err == users.ErrUserNotFoundDB {
			webserver_utils.WriteError(w, http.StatusUnauthorized, "authorization denied")
			http.Redirect(w, r, fmt.Sprintf("https://%s", utils.GetEnv()["DOMAIN_NAME"]), 301)
			return
		} else if err != nil {
			utils.Logger.Printf("Ошибка при получении информации о %s (%v)\n", email_from_cookie.Value, err)
			webserver_utils.WriteError(w, http.StatusUnauthorized, err.Error())
			return
		}
		webserver_utils.WriteJSON(w, 200, map[string]string{"first_name": user.First_name, "last_name": user.Last_name, "father_name": user.Father_name})
		return
	}
}
