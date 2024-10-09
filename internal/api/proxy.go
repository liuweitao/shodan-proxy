package api

import (
    "bytes"
    "io"
    "io/ioutil"
    "log"
    "net/http"
    "net/http/httputil"
    "net/url"
    "regexp"
    "strings"
    "encoding/json"

    "shodan-proxy/internal/utils"
    "shodan-proxy/pkg/api_paths"
)

func ShodanProxy(w http.ResponseWriter, r *http.Request) {
    targetURL, _ := url.Parse("https://api.shodan.io")

    // 检查路径是否被允许
    pathAllowed := false
    for _, allowedPath := range api_paths.AllowedPaths {
        if allowedPath.MatchString(r.URL.Path) {
            pathAllowed = true
            break
        }
    }

    if !pathAllowed {
        http.Error(w, "This path is not allowed", http.StatusForbidden)
        return
    }

    // 检查路径是否被阻止（保留原有的检查）
    if utils.IsPathBlocked(r.URL.Path) {
        http.Error(w, "This path is blocked", http.StatusForbidden)
        return
    }

    // 创建一个新的请求，而不是使用反向代理
    newReq, err := http.NewRequest(r.Method, targetURL.String()+r.URL.Path, nil)
    if err != nil {
        http.Error(w, "Error creating request", http.StatusInternalServerError)
        return
    }

    // 复制查询参数
    newReq.URL.RawQuery = r.URL.RawQuery

    // 如果是 POST 请求，读取并设置请求体
    if r.Method == "POST" {
        body, err := ioutil.ReadAll(r.Body)
        if err != nil {
            http.Error(w, "Error reading request body", http.StatusInternalServerError)
            return
        }
        newReq.Body = ioutil.NopCloser(bytes.NewBuffer(body))
        newReq.ContentLength = int64(len(body))
        newReq.Header.Set("Content-Type", r.Header.Get("Content-Type"))
    }

    // 处理 API key
    q := newReq.URL.Query()
    userKey := q.Get("key")

    if userKey == "" || userKey == "shodanproxy" {
        apiKey := utils.GetNextKey()
        if apiKey == "" {
            log.Printf("没有可用的 API keys")
            http.Error(w, "No API keys available", http.StatusServiceUnavailable)
            return
        }
        q.Set("key", apiKey)
        newReq.URL.RawQuery = q.Encode()
        if userKey == "" {
            log.Printf("用户未提供 API key，使用代理的 API key")
        } else {
            log.Printf("用户请求使用代理的 API key")
        }
    } else {
        log.Printf("使用用户提供的 API key")
    }

    // 打印最终的请求信息
    requestDump, err := httputil.DumpRequestOut(newReq, false)
    if err != nil {
        log.Printf("Error dumping request: %v", err)
    } else {
        // 在日志中隐藏 API key
        loggedRequest := string(requestDump)
        loggedRequest = regexp.MustCompile(`key=[^&]+`).ReplaceAllString(loggedRequest, "key=REDACTED")
        log.Printf("Final request to Shodan:\n%s", loggedRequest)
    }

    // 发送请求
    client := &http.Client{}
    resp, err := client.Do(newReq)
    if err != nil {
        http.Error(w, "Error sending request to Shodan", http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    // 复制响应头
    for k, vv := range resp.Header {
        for _, v := range vv {
            w.Header().Add(k, v)
        }
    }
    w.WriteHeader(resp.StatusCode)

    // 读取响应体
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Printf("Error reading response body: %v", err)
        http.Error(w, "Error reading response from Shodan", http.StatusInternalServerError)
        return
    }

    // 如果Content-Type是JSON，尝试替换API key
    contentType := resp.Header.Get("Content-Type")
    if strings.Contains(contentType, "application/json") {
        // 将响应体解析为JSON
        var jsonBody interface{}
        err = json.Unmarshal(body, &jsonBody)
        if err == nil {
            // 递归替换JSON中的API key
            replaceAPIKey(jsonBody, userKey)
            // 重新编码JSON
            body, err = json.Marshal(jsonBody)
            if err != nil {
                log.Printf("Error re-encoding JSON: %v", err)
            }
        } else {
            log.Printf("Error parsing JSON response: %v", err)
        }
    }

    // 写入修改后的响应体
    _, err = w.Write(body)
    if err != nil {
        log.Printf("Error writing response: %v", err)
    }
}

// 递归替换JSON中的API key
func replaceAPIKey(v interface{}, newKey string) {
    switch vv := v.(type) {
    case map[string]interface{}:
        for k, v := range vv {
            if k == "api_key" {
                vv[k] = newKey
            } else {
                replaceAPIKey(v, newKey)
            }
        }
    case []interface{}:
        for _, v := range vv {
            replaceAPIKey(v, newKey)
        }
    }
}
