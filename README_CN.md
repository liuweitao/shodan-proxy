# Shodan 代理

[English README](README.md)

Shodan 代理是一个基于 Go 语言的代理服务器，用于处理 Shodan API 请求，并提供额外功能如 IP 过滤、路径阻止和管理面板。

## 功能特性

- 代理 Shodan API 请求
- IP 白名单
- 路径阻止
- 安全的配置管理面板，带有身份验证功能
- 多个 Shodan API 密钥管理，采用轮询调用机制

## 安装

1. 创建并进入项目目录：
   ```
   mkdir shodan-proxy && cd shodan-proxy
   ```
2. 下载 compose.yaml 文件：
   ```
   curl -O https://raw.githubusercontent.com/liuweitao/shodan-proxy/main/compose.yaml
   ```
3. 启动容器：
   ```
   docker compose up -d
   ```

## 配置

编辑 `config/config.yaml` 文件以设置：

- 阻止的路径
- 允许的 IP
- 受信任的代理
- 管理员凭据

注意：确保您的 `config.yaml` 文件得到适当的保护，特别是如果它包含敏感信息如 API 密钥。

## 使用方法

### 管理面板

访问 `http://localhost:8080/admin` 进入管理面板，管理设置和 API 密钥。

默认管理员账号和密码如下：
- 用户名：admin
- 密码：shodanproxy

**注意：** 出于安全考虑，强烈建议您在首次登录后立即更改默认密码。

### API 调用示例

以下是一些 Shodan API 调用的示例，对比了官方 API 和本代理的调用方式：

1. 搜索主机信息

   官方 API:
   ```
   https://api.shodan.io/shodan/host/search?key=YOUR_API_KEY&query=apache
   ```

   本代理:
   ```
   http://localhost:8080/shodan/host/search?query=apache
   ```

2. 获取特定 IP 的信息

   官方 API:
   ```
   https://api.shodan.io/shodan/host/1.1.1.1?key=YOUR_API_KEY
   ```

   本代理:
   ```
   http://localhost:8080/shodan/host/1.1.1.1
   ```

3. 获取当前 API 计划的信息

   官方 API:
   ```
   https://api.shodan.io/api-info?key=YOUR_API_KEY
   ```

   本代理:
   ```
   http://localhost:8080/api-info
   ```

注意：
1. 使用本代理时，通常无需在每个请求中包含 API 密钥。代理会自动管理和轮询使用配置的 API 密钥。
2. 如果在调用时传入了 key 参数（例如：`http://localhost:8080/api-info?key=YOUR_API_KEY`），则代理将使用调用者提供的 key。如果没有传入 key 参数，代理将使用自身配置的 API 密钥。这种灵活性允许用户在需要时使用自己的 API 密钥，同时也能利用代理的密钥管理功能。
3. 出于安全考虑，建议在受控环境中使用此代理服务器，不要将其直接暴露在公共互联网上。

## 安全注意事项

- 始终为管理面板使用强大且唯一的密码。
- 定期更新 Docker 镜像以确保您拥有最新的安全补丁。
- 谨慎将代理服务器暴露在互联网上。建议在受控网络环境中使用它。
- 定期审查和更新您的 IP 白名单和路径阻止规则。

## 贡献

欢迎贡献！请随时提交 Pull Request。

## 许可证

本项目采用 MIT 许可证 - 详情请见 [LICENSE](LICENSE) 文件。

## 依赖管理

本项目使用 Go modules 进行依赖管理。`go.mod` 文件在 Docker 构建过程中动态生成，以确保使用最新的兼容依赖版本。如果您在 Docker 环境外进行开发，可以通过运行以下命令生成 `go.mod` 文件：

```
go mod init shodan-proxy
go mod tidy
```

## 联系方式

如有任何问题或建议，请开启一个 issue 或直接联系项目维护者。

感谢您对 Shodan 代理项目的关注！