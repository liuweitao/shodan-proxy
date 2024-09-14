# Shodan Proxy

[中文版 README](README_CN.md)

Shodan Proxy is a Go-based proxy server for Shodan API requests with additional features like IP filtering, path blocking, and an admin panel.

## Features

- Proxy Shodan API requests
- IP whitelisting
- Path blocking
- Secure configuration admin panel with authentication
- Multiple Shodan API key management with round-robin rotation

## Installation

1. Create and enter the project directory:
   ```bash
   mkdir shodan-proxy && cd shodan-proxy
   ```
2. Download the compose.yaml file:
   ```bash
   curl -O https://raw.githubusercontent.com/liuweitao/shodan-proxy/main/compose.yaml
   ```
3. Start the container:
   ```bash
   docker compose up -d
   ```

## Configuration

Edit the `config/config.yaml` file to set up:

- Blocked paths
- Allowed IPs
- Trusted proxies
- Admin credentials

Note: Ensure that your `config.yaml` file is properly secured, especially if it contains sensitive information like API keys.

## Usage

### Admin Panel

Access the admin panel at `http://localhost:8080/admin` to manage settings and API keys.

Default admin credentials:
- Username: admin
- Password: shodanproxy

**Note:** For security reasons, it is strongly recommended to change the default password after your first login.

### API Call Examples

Here are some examples of Shodan API calls, comparing the official API and this proxy's usage:

1. Search for host information

   Official API:
   ```
   https://api.shodan.io/shodan/host/search?key=YOUR_API_KEY&query=apache
   ```

   This proxy:
   ```
   http://localhost:8080/shodan/host/search?query=apache
   ```

2. Get information for a specific IP

   Official API:
   ```
   https://api.shodan.io/shodan/host/1.1.1.1?key=YOUR_API_KEY
   ```

   This proxy:
   ```
   http://localhost:8080/shodan/host/1.1.1.1
   ```

3. Get information about the current API plan

   Official API:
   ```
   https://api.shodan.io/api-info?key=YOUR_API_KEY
   ```

   This proxy:
   ```
   http://localhost:8080/api-info
   ```

Note:
1. When using this proxy, you typically don't need to include the API key in each request. The proxy will automatically manage and rotate the configured API keys.
2. If a key parameter is passed in the call (e.g., `http://localhost:8080/api-info?key=YOUR_API_KEY`), the proxy will use the key provided by the caller. If no key parameter is passed, the proxy will use its own configured API keys. This flexibility allows users to use their own API keys when needed, while also leveraging the proxy's key management functionality.
3. For security reasons, it's recommended to use this proxy server within a controlled environment and not expose it directly to the public internet.

## Security Considerations

- Always use strong, unique passwords for the admin panel.
- Regularly update the Docker image to ensure you have the latest security patches.
- Be cautious when exposing the proxy server to the internet. It's recommended to use it within a controlled network environment.
- Regularly review and update your IP whitelist and path blocking rules.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Dependencies

This project uses Go modules for dependency management. The `go.mod` file is generated during the Docker build process to ensure the latest compatible versions of dependencies are used. If you're developing outside of Docker, you can generate the `go.mod` file by running:

```
go mod init shodan-proxy
go mod tidy
```

## Contact

If you have any questions or suggestions, please open an issue or contact the project maintainer directly.

Thank you for your interest in the Shodan Proxy project!