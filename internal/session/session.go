package session

import (
    "net/http"
    "time"
    "sync"
    "github.com/google/uuid"
)

var (
    sessions = make(map[string]string)
    mutex    sync.RWMutex
)

func GenerateSessionToken() string {
    // 生成一个唯一的会话令牌
    return uuid.New().String()
}

func SaveSession(token, username string) {
    mutex.Lock()
    defer mutex.Unlock()
    sessions[token] = username
}

func GetSession(token string) (string, bool) {
    mutex.RLock()
    defer mutex.RUnlock()
    username, found := sessions[token]
    return username, found
}

func ClearSession(w http.ResponseWriter) {
    http.SetCookie(w, &http.Cookie{
        Name:     "session_token",
        Value:    "",
        Expires:  time.Now().Add(-1 * time.Hour),
        HttpOnly: true,
        Path:     "/",
    })
}
