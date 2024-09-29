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

    // 复制响应体
    _, err = io.Copy(w, resp.Body)
    if err != nil {
        log.Printf("Error copying response: %v", err)
    }
}
