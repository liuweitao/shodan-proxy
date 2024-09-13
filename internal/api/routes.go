package api

import (
    "net/http"

    "shodan-proxy/internal/middleware"
    "shodan-proxy/internal/config"
)

func SetupRoutes() *http.ServeMux {
    mux := http.NewServeMux()

    fs := http.FileServer(http.Dir("./public"))
    mux.Handle("/static/", http.StripPrefix("/static/", fs))

    mux.HandleFunc("/", middleware.IPCheckMiddleware(ShodanProxy))
    mux.HandleFunc("/login", LoginHandler)
    mux.HandleFunc("/admin", middleware.AuthMiddleware(AdminHandler))
    mux.HandleFunc("/update-config", middleware.AuthMiddleware(UpdateConfigHandler))
    mux.HandleFunc("/tools/httpheaders", HandleHTTPHeaders) 
    mux.HandleFunc("/tools/myip", HandleMyIP) 

    mux.HandleFunc("/get-config", middleware.AuthMiddleware(config.ServeConfig))

    mux.HandleFunc("/api/shodan-keys", middleware.AuthMiddleware(config.ServeShodanKeys))

    mux.HandleFunc("/logout", LogoutHandler)

    return mux
}
