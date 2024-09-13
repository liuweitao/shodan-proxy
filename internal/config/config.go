package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/yaml.v2"
	"encoding/json"
	"net/http"
	"time"
)

var (
	GlobalConfig = Config{
		// ... 其他配置 ...
		AdminUser: AdminUser{
			Username: "admin", // 确保这里设置了正确的默认用户名
			Password: "shodanproxy", // 确保这里设置了正确的默认密码
		},
	}
	CurrentKey   int
	KeyMutex     sync.Mutex
	ConfigMutex  sync.RWMutex
	ShodanKeys []string
	ShodanKeysMutex sync.RWMutex
)

type Config struct {
	BlockedPaths   []string   `json:"blocked_paths" yaml:"blocked_paths"`
	AllowedIPs     []string   `json:"allowed_ips" yaml:"allowed_ips"`
	TrustedProxies []string   `json:"trusted_proxies" yaml:"trusted_proxies"`
	AdminUser      AdminUser  `json:"admin_user" yaml:"admin_user"`
}

type AdminUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

const (
	DefaultAdminUsername = "admin"
	DefaultAdminPassword = "shodanproxy"
)

func LoadConfig() error {
	configPath := "config/config.yaml"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := createDefaultConfig(configPath); err != nil {
            return err
        }
	}

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("error reading config file: %v", err)
	}

	ConfigMutex.Lock()
	defer ConfigMutex.Unlock()

	err = yaml.Unmarshal(data, &GlobalConfig)
	if err != nil {
		return fmt.Errorf("error unmarshaling config: %v", err)
	}

	if GlobalConfig.AdminUser.Username == "" {
		GlobalConfig.AdminUser.Username = DefaultAdminUsername
	}

	if GlobalConfig.AdminUser.Password == "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(DefaultAdminPassword), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("error hashing default password: %v", err)
		}
		GlobalConfig.AdminUser.Password = string(hashedPassword)
	}

	// 加载 Shodan keys
	return LoadShodanKeys()
}

func createDefaultConfig(configPath string) error {
	defaultConfig := Config{
		BlockedPaths:  []string{"/account", "/org", "/api-info"},
		AdminUser: AdminUser{
			Username: "",
			Password: "",
		},
		AllowedIPs:     []string{},
		TrustedProxies: []string{},
	}

	data, err := yaml.Marshal(defaultConfig)
	if err != nil {
		return fmt.Errorf("error marshaling default config: %v", err)
	}

	err = os.MkdirAll(filepath.Dir(configPath), 0755)
	if err != nil {
		return fmt.Errorf("error creating config directory: %v", err)
	}

	err = ioutil.WriteFile(configPath, data, 0644)
	if err != nil {
		return fmt.Errorf("error writing default config file: %v", err)
	}

	log.Println("Created default config file at", configPath)
	return nil
}

func LoadShodanKeys() error {
    configPath := "config/shodan_keys.yaml"
    if _, err := os.Stat(configPath); os.IsNotExist(err) {
        return createDefaultShodanKeysConfig(configPath)
    }

    data, err := ioutil.ReadFile(configPath)
    if err != nil {
        return fmt.Errorf("error reading Shodan keys file: %v", err)
    }

    lines := strings.Split(string(data), "\n")
    var comments []string
    var keys []string

    for _, line := range lines {
        trimmedLine := strings.TrimSpace(line)
        if strings.HasPrefix(trimmedLine, "#") || trimmedLine == "" {
            comments = append(comments, line)
        } else {
            // 移除 "- " 前缀并去除空白
            key := strings.TrimSpace(strings.TrimPrefix(trimmedLine, "-"))
            if key != "" {
                keys = append(keys, key)
            }
        }
    }

    ShodanKeysMutex.Lock()
    defer ShodanKeysMutex.Unlock()

    ShodanKeys = keys
    log.Printf("Loaded %d Shodan keys", len(ShodanKeys))
    return nil
}

func createDefaultShodanKeysConfig(configPath string) error {
    defaultContent := `# Shodan API Keys
# Add your Shodan API keys below, one per line
# Example:
# - XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
`

    err := os.MkdirAll(filepath.Dir(configPath), 0755)
    if err != nil {
        return fmt.Errorf("error creating config directory: %v", err)
    }

    err = ioutil.WriteFile(configPath, []byte(defaultContent), 0644)
    if err != nil {
        return fmt.Errorf("error writing default Shodan keys file: %v", err)
    }

    log.Println("Created default Shodan keys file at", configPath)
    return nil
}

func SaveShodanKeys() error {
    ShodanKeysMutex.Lock()
    defer ShodanKeysMutex.Unlock()
    return saveShodanKeysUnsafe()
}

func saveShodanKeysUnsafe() error {
    log.Printf("Saving %d Shodan keys", len(ShodanKeys))

    data, err := ioutil.ReadFile("config/shodan_keys.yaml")
    if err != nil && !os.IsNotExist(err) {
        return fmt.Errorf("error reading existing Shodan keys file: %v", err)
    }

    var comments []string
    if err == nil {
        lines := strings.Split(string(data), "\n")
        for _, line := range lines {
            trimmedLine := strings.TrimSpace(line)
            if strings.HasPrefix(trimmedLine, "#") || trimmedLine == "" {
                comments = append(comments, line)
            } else {
                break
            }
        }
    }

    // 移除注释末尾的空行
    for len(comments) > 0 && strings.TrimSpace(comments[len(comments)-1]) == "" {
        comments = comments[:len(comments)-1]
    }

    var output strings.Builder
    for _, comment := range comments {
        output.WriteString(comment + "\n")
    }

    // 如果有注释且有密钥，添加一个空行
    if len(comments) > 0 && len(ShodanKeys) > 0 {
        output.WriteString("\n")
    }

    for _, key := range ShodanKeys {
        output.WriteString("- " + key + "\n")
    }

    // 移除最后一个可能的空行
    outputStr := strings.TrimSpace(output.String()) + "\n"

    err = ioutil.WriteFile("config/shodan_keys.yaml", []byte(outputStr), 0644)
    if err != nil {
        log.Printf("Error writing Shodan keys file: %v", err)
        return fmt.Errorf("error writing Shodan keys file: %v", err)
    }

    log.Printf("Shodan keys saved successfully")
    return nil
}

func AddShodanKey(key string) error {
    // 首先验证 key 是否有效
    if err := validateShodanKey(key); err != nil {
        return fmt.Errorf("invalid Shodan API key: %v", err)
    }

    ShodanKeysMutex.Lock()
    defer ShodanKeysMutex.Unlock()

    // 检查 key 是否已存在
    for _, existingKey := range ShodanKeys {
        if existingKey == key {
            return fmt.Errorf("key already exists")
        }
    }

    // 添加新 key
    ShodanKeys = append(ShodanKeys, key)

    // 保存到文件
    return saveShodanKeysUnsafe()
}

func validateShodanKey(key string) error {
    client := &http.Client{
        Timeout: 10 * time.Second,
    }

    url := fmt.Sprintf("https://api.shodan.io/api-info?key=%s", key)
    resp, err := client.Get(url)
    if err != nil {
        return fmt.Errorf("error contacting Shodan API: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("invalid API key (status code: %d)", resp.StatusCode)
    }

    return nil
}

func DeleteShodanKey(key string) error {
    ShodanKeysMutex.Lock()
    defer ShodanKeysMutex.Unlock()

    log.Printf("Attempting to delete Shodan key: %s", key)
    
    newKeys := make([]string, 0, len(ShodanKeys))
    keyFound := false
    
    for _, k := range ShodanKeys {
        if k != key {
            newKeys = append(newKeys, k)
        } else {
            keyFound = true
        }
    }
    
    if !keyFound {
        log.Printf("Key not found: %s", key)
        return fmt.Errorf("key not found")
    }
    
    ShodanKeys = newKeys
    
    log.Printf("Key deleted. Saving updated Shodan keys...")
    err := saveShodanKeysUnsafe()
    if err != nil {
        log.Printf("Error saving Shodan keys: %v", err)
        return err
    }
    log.Printf("Shodan keys updated successfully")
    
    return nil
}

func GetMaskedShodanKeys() []string {
	ShodanKeysMutex.RLock()
	defer ShodanKeysMutex.RUnlock()

	maskedKeys := make([]string, len(ShodanKeys))
	for i, key := range ShodanKeys {
		if len(key) > 6 {
			maskedKeys[i] = key[:6] + strings.Repeat("*", len(key)-6)
		} else {
			maskedKeys[i] = key
		}
	}
	return maskedKeys
}

func ServeConfig(w http.ResponseWriter, r *http.Request) {
    ConfigMutex.RLock()
    defer ConfigMutex.RUnlock()

    config := struct {
        BlockedPaths   []string  `json:"blocked_paths"`
        AllowedIPs     []string  `json:"allowed_ips"`
        TrustedProxies []string  `json:"trusted_proxies"`
        AdminUser      AdminUser `json:"admin_user"`
    }{
        BlockedPaths:   GlobalConfig.BlockedPaths,
        AllowedIPs:     GlobalConfig.AllowedIPs,
        TrustedProxies: GlobalConfig.TrustedProxies,
        AdminUser: AdminUser{
            Username: GlobalConfig.AdminUser.Username,
            Password: "", // 不发送密码
        },
    }

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(config); err != nil {
        http.Error(w, "Error encoding config", http.StatusInternalServerError)
    }
}

func ServeShodanKeys(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        ShodanKeysMutex.RLock()
        defer ShodanKeysMutex.RUnlock()

        w.Header().Set("Content-Type", "application/json")
        if err := json.NewEncoder(w).Encode(ShodanKeys); err != nil {
            http.Error(w, "Error encoding Shodan keys", http.StatusInternalServerError)
        }
    case http.MethodPost:
        var newKey struct {
            Key string `json:"key"`
        }
        if err := json.NewDecoder(r.Body).Decode(&newKey); err != nil {
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }
        if newKey.Key == "" {
            http.Error(w, "Key cannot be empty", http.StatusBadRequest)
            return
        }
        if err := AddShodanKey(newKey.Key); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        w.WriteHeader(http.StatusCreated)
    case http.MethodDelete:
        var keyToDelete struct {
            Key string `json:"key"`
        }
        body, err := ioutil.ReadAll(r.Body)
        if err != nil {
            log.Printf("Error reading request body: %v", err)
            http.Error(w, "Error reading request body", http.StatusBadRequest)
            return
        }
        log.Printf("Received body: %s", string(body))

        if err := json.Unmarshal(body, &keyToDelete); err != nil {
            log.Printf("Error unmarshaling JSON: %v", err)
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }
        if keyToDelete.Key == "" {
            http.Error(w, "Key cannot be empty", http.StatusBadRequest)
            return
        }
        if err := DeleteShodanKey(keyToDelete.Key); err != nil {
            log.Printf("Error deleting Shodan key: %v", err)
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Key deleted successfully"))
        log.Printf("Key deleted successfully")
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

func SaveConfig() error {
	ConfigMutex.Lock()
	defer ConfigMutex.Unlock()

	// 创建一个不包含 ShodanAPIKeys 的配置副本
	configToSave := Config{
		BlockedPaths:   GlobalConfig.BlockedPaths,
		AllowedIPs:     GlobalConfig.AllowedIPs,
		TrustedProxies: GlobalConfig.TrustedProxies,
		AdminUser: AdminUser{
			Username: GlobalConfig.AdminUser.Username,
			Password: GlobalConfig.AdminUser.Password, // 保持原密码不变
		},
	}

	// 如果前端传来的密码为空，我们就使用现有的密码
	if configToSave.AdminUser.Password == "" {
		configToSave.AdminUser.Password = GlobalConfig.AdminUser.Password
	}

	data, err := yaml.Marshal(configToSave)
	if err != nil {
		return fmt.Errorf("error marshaling config: %v", err)
	}

	err = ioutil.WriteFile("config/config.yaml", data, 0644)
	if err != nil {
		return fmt.Errorf("error writing config file: %v", err)
	}

	// 更新 GlobalConfig
	GlobalConfig = configToSave

	return nil
}
