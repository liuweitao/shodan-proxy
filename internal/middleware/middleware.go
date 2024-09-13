package middleware

import (
	"net/http"

	"shodan-proxy/internal/utils"
	"shodan-proxy/internal/session"
	"log"
	"shodan-proxy/internal/config"  // 添加这行
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		sessionToken := cookie.Value
		username, found := session.GetSession(sessionToken)
		if !found {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// 可以在这里添加日志
		log.Printf("Authenticated user %s accessing admin page", username)

		next.ServeHTTP(w, r)
	}
}

func IPCheckMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientIP := utils.GetClientIP(r)
		log.Printf("Checking IP: %s", clientIP)
		log.Printf("X-Forwarded-For: %s", r.Header.Get("X-Forwarded-For"))
		log.Printf("X-Real-IP: %s", r.Header.Get("X-Real-IP"))
		log.Printf("RemoteAddr: %s", r.RemoteAddr)
		
		config.ConfigMutex.RLock()
		log.Printf("Allowed IPs: %v", config.GlobalConfig.AllowedIPs)
		config.ConfigMutex.RUnlock()
		
		if !utils.IsIPAllowed(clientIP) {
			log.Printf("IP %s is not allowed", clientIP)
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		log.Printf("IP %s is allowed", clientIP)
		next.ServeHTTP(w, r)
	}
}
