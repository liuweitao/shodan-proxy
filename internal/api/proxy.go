package api

import (
    "log"
    "net/http"
    "net/http/httputil"
    "net/url"
    "regexp"
    "io/ioutil"
    "encoding/json"
    "bytes"  // 添加这一行
    "compress/gzip"

    "shodan-proxy/pkg/api_paths"
    "shodan-proxy/internal/utils"
)

func ShodanProxy(w http.ResponseWriter, r *http.Request) {
    targetURL, _ := url.Parse("https://api.shodan.io")
    proxy := httputil.NewSingleHostReverseProxy(targetURL)

    // 打印原始请求信息
    log.Printf("Original request method: %s", r.Method)
    log.Printf("Original request URL: %s", r.URL.String())
    log.Printf("Original request headers: %v", r.Header)

    path := r.URL.Path

    // 检查路径是否被阻止
    if utils.IsPathBlocked(path) {
        log.Printf("Access denied for blocked path: %s", path)
        http.Error(w, "Access denied", http.StatusForbidden)
        return
    }

    // 检查路径是否在允许列表中
    allowed := false
    for _, allowedPath := range api_paths.AllowedPaths {
        if allowedPath.MatchString(path) {
            allowed = true
            break
        }
    }

    if !allowed {
        log.Printf("Access denied for path not in allowed list: %s", path)
        http.Error(w, "Access denied", http.StatusForbidden)
        return
    }

    log.Printf("Access allowed for path: %s", path)

    q := r.URL.Query()
    userKey := q.Get("key")

    if userKey != "" {
        log.Printf("User provided their own API key")
        // 选项1：保留用户的key
        // 不做任何操作，保留用户的key

        // 选项2：覆盖用户的key，但记录警告
        // apiKey := utils.GetNextKey()
        // if apiKey == "" {
        //     log.Printf("No API keys available")
        //     http.Error(w, "No API keys available", http.StatusServiceUnavailable)
        //     return
        // }
        // log.Printf("Warning: Overriding user-provided API key")
        // q.Set("key", apiKey)
    } else {
        apiKey := utils.GetNextKey()
        if apiKey == "" {
            log.Printf("No API keys available")
            http.Error(w, "No API keys available", http.StatusServiceUnavailable)
            return
        }
        q.Set("key", apiKey)
    }

    r.URL.Scheme = "https"
    r.URL.Host = targetURL.Host
    r.Host = targetURL.Host

    // 打印修改后的请求信息
    log.Printf("Modified request URL: %s", r.URL.String())
    log.Printf("Modified request headers: %v", r.Header)

    // 使用自定义的 Director 函数
    originalDirector := proxy.Director
    proxy.Director = func(req *http.Request) {
        originalDirector(req)
        // 打印最终的请求信息
        requestDump, err := httputil.DumpRequestOut(req, true)
        if err != nil {
            log.Printf("Error dumping request: %v", err)
        } else {
            // 在日志中隐藏 API key
            loggedRequest := string(requestDump)
            loggedRequest = regexp.MustCompile(`key=[^&]+`).ReplaceAllString(loggedRequest, "key=REDACTED")
            log.Printf("Final request to Shodan:\n%s", loggedRequest)
        }
    }

    // 使用自定义的 ModifyResponse 函数
    proxy.ModifyResponse = func(resp *http.Response) error {
        log.Printf("Response Status: %s", resp.Status)
        log.Printf("Response Headers: %v", resp.Header)

        // 检查响应状态码
        if resp.StatusCode != http.StatusOK {
            var body []byte
            var err error

            // Check if the response is gzip encoded
            if resp.Header.Get("Content-Encoding") == "gzip" {
                reader, err := gzip.NewReader(resp.Body)
                if err != nil {
                    log.Printf("Error creating gzip reader: %v", err)
                    return err
                }
                defer reader.Close()
                body, err = ioutil.ReadAll(reader)
            } else {
                body, err = ioutil.ReadAll(resp.Body)
            }

            if err != nil {
                log.Printf("Error reading response body: %v", err)
                return err
            }
            resp.Body.Close()

            // 解析错误信息
            var errorResponse struct {
                Error string `json:"error"`
            }
            if err := json.Unmarshal(body, &errorResponse); err != nil {
                log.Printf("Error parsing error response: %v", err)
                log.Printf("Raw response body: %s", string(body))
                errorResponse.Error = "Unknown error occurred"
            }

            // 创建新的响应
            newBody, _ := json.Marshal(map[string]interface{}{
                "error": errorResponse.Error,
                "status_code": resp.StatusCode,
            })

            // 设置新的响应
            resp.Body = ioutil.NopCloser(bytes.NewBuffer(newBody))
            resp.ContentLength = int64(len(newBody))
            resp.Header.Set("Content-Type", "application/json")
            resp.Header.Del("Content-Encoding") // Remove Content-Encoding header
            resp.StatusCode = http.StatusOK // 将状态码改为 200

            log.Printf("Modified error response: %s", string(newBody))
        }

        return nil
    }

    proxy.ServeHTTP(w, r)
}
