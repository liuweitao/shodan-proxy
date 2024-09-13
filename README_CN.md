# Shodan 代理

[English README](README.md)

Shodan 代理是一个基于 Go 语言的代理服务器，用于处理 Shodan API 请求，并提供额外功能如 IP 过滤、路径阻止和管理面板。

## 功能特性

- 代理 Shodan API 请求
- IP 白名单
- 路径阻止
- 配置管理面板
- 多个 Shodan API 密钥管理，采用轮询调用机制


## 安装

1. 克隆仓库
2. 构建 Docker 镜像：
   ```
   docker compose build
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

## 使用方法

访问 `http://localhost:8080/admin` 进入管理面板，管理设置和 API 密钥。

默认管理员账号和密码如下：
- 用户名：admin
- 密码：shodanproxy

**注意：** 出于安全考虑，强烈建议您在首次登录后立即更改默认密码。

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