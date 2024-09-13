package api

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"  // 添加这行
	"net/http"
	"strings"
	"time"

	"shodan-proxy/internal/config"
	"shodan-proxy/internal/session"
	"shodan-proxy/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

func AdminHandler(w http.ResponseWriter, r *http.Request) {
	// 添加方法检查
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tmpl, err := template.ParseFiles("templates/admin.html")
	if err != nil {
		log.Printf("Error parsing admin template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	config.ConfigMutex.RLock()
	data := struct {
		Config config.Config
	}{
		Config: config.GlobalConfig,
	}
	config.ConfigMutex.RUnlock()

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing admin template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// 处理 GET 请求，显示登录页面
		tmpl, err := template.ParseFiles("templates/login.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	case http.MethodPost:
		// 处理 POST 请求，验证登录
		var loginData struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		err := json.NewDecoder(r.Body).Decode(&loginData)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		config.ConfigMutex.RLock()
		defer config.ConfigMutex.RUnlock()

		if loginData.Username == config.GlobalConfig.AdminUser.Username {
			err := bcrypt.CompareHashAndPassword([]byte(config.GlobalConfig.AdminUser.Password), []byte(loginData.Password))
			if err == nil {
				sessionToken := session.GenerateSessionToken()
				session.SaveSession(sessionToken, loginData.Username)
				http.SetCookie(w, &http.Cookie{
					Name:     "session_token",
					Value:    sessionToken,
					HttpOnly: true,
					Path:     "/",
				})
				w.WriteHeader(http.StatusOK)
			} else {
				http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			}
		} else {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func UpdateConfigHandler(w http.ResponseWriter, r *http.Request) {
	var newConfig config.Config
	err := json.NewDecoder(r.Body).Decode(&newConfig)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 如果密码为空，保持原密码不变
	if newConfig.AdminUser.Password == "" {
		newConfig.AdminUser.Password = config.GlobalConfig.AdminUser.Password
	} else {
		// 如果提供了新密码，对其进行哈希处理
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newConfig.AdminUser.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}
		newConfig.AdminUser.Password = string(hashedPassword)
	}

	// 更新全局配置
	config.GlobalConfig = newConfig

	// 保存配置
	err = config.SaveConfig()
	if err != nil {
		http.Error(w, "Error saving configuration", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Configuration updated successfully"))
}

func validateConfig(cfg *config.Config) error {
	log.Printf("Validating config: %+v", cfg)

	// 移除对 ShodanAPIKeys 的验证，因为它们现在单独管理
	// 其他验证逻辑保持不变
	if cfg.AdminUser.Username == "" {
		return fmt.Errorf("admin username cannot be empty")
	}

	// 如果提供了新密码，进行哈希处理
	if cfg.AdminUser.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(cfg.AdminUser.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("error hashing password: %v", err)
		}
		cfg.AdminUser.Password = string(hashedPassword)
	}

	return nil
}

func HandleHTTPHeaders(w http.ResponseWriter, r *http.Request) {
	headers := make(map[string]string)
	for name, values := range r.Header {
		headers[name] = strings.Join(values, ", ")
	}

	jsonResponse, err := json.Marshal(headers)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func HandleMyIP(w http.ResponseWriter, r *http.Request) {
	ip := utils.GetClientIP(r)
	w.Write([]byte(ip))
}

// LogoutHandler handles user logout
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Clear the session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	w.WriteHeader(http.StatusOK)
}
