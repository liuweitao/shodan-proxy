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
)

func ShodanProxy(w http.ResponseWriter, r *http.Request) {
    targetURL, _ := url.Parse("https://api.shodan.io")

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

    // 设置必要的 headers
    newReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

    // 处理 API key
    q := newReq.URL.Query()
    userKey := q.Get("key")

    if userKey == "" {
        apiKey := utils.GetNextKey()
        if apiKey == "" {
            log.Printf("没有可用的 API keys")
            http.Error(w, "No API keys available", http.StatusServiceUnavailable)
            return
        }
        q.Set("key", apiKey)
        newReq.URL.RawQuery = q.Encode()
        log.Printf("使用代理的 API key")
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
